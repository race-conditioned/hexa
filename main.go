package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"hexa/m/v2/endurance"
	"hexa/m/v2/fusion/dt"
	"hexa/m/v2/fusion/intake"
	"hexa/m/v2/horizon"
	"hexa/m/v2/horizon/ports/inbound"
	"hexa/m/v2/symphony"
)

type AppCtxI interface {
	inbound.Ctx
	Metrics() Metrics
	Limiter() Limiter
	Timeout() time.Duration
	Idempotency() IdempotencyI[inbound.Result]
}

type AppCtx struct {
	context.Context
	metrics          Metrics
	limiter          Limiter
	idempotencyStore IdempotencyI[inbound.Result]
}

func (c *AppCtx) Metrics() Metrics                          { return c.metrics }
func (c *AppCtx) Limiter() Limiter                          { return c.limiter }
func (c *AppCtx) Timeout() time.Duration                    { return 2 * time.Second }
func (c *AppCtx) Idempotency() IdempotencyI[inbound.Result] { return c.idempotencyStore }

type AppHandler = inbound.UnaryHandler[AppCtxI, inbound.Command, inbound.Result]

func main() {
	// Use case (app layer) // this is in the users code, the app layer providing use cases
	uc := NewTransferUsecase(&immediateDispatcher{}, dummyMetrics{})

	// transferEndpoints := provider.Provide(uc)
	ctx := &AppCtx{metrics: dummyMetrics{}, limiter: dummyLimiter{}, idempotencyStore: newMemIdempotency[inbound.Result]()}
	gw := horizon.NewGateway[AppCtxI](ctx)

	// pre and post
	metricsIn := symphony.PolicyStage("m_in")
	metricsOut := symphony.PolicyStage("m_out")

	pre := symphony.Order(metricsIn)
	post := symphony.Order(metricsOut)

	// policies
	idempotency := symphony.PolicyStage("idempotency")
	rateLimit := symphony.PolicyStage("rate_limit")
	timeout := symphony.PolicyStage("timeout")
	latency := symphony.PolicyStage("latency")

	s := symphony.New(pre, post, symphony.WithPre(metricsIn, CountRequests), symphony.WithPost(metricsOut, CountSuccess))

	mid := symphony.Order(DefaultPolicyOrder...)

	// Compose endpoint
	comp := symphony.Compose(
		s,
		mid,
		symphony.WithPolicy(rateLimit, symphony.Lift[AppCtxI, TransferCommand, TransferResult](RateLimit)),
		symphony.WithPolicy(timeout, symphony.Lift[AppCtxI, TransferCommand, TransferResult](Timeout)),
		symphony.WithPolicy(latency, symphony.Lift[AppCtxI, TransferCommand, TransferResult](ObserveLatency)),
		symphony.WithPolicy(idempotency, symphony.LiftCap[AppCtxI, TransferCommand, TransferResult](Idempotency)),
	)

	h := comp.Wrap(endurance.Transport(uc.SubmitTransfer, printbefore, dosomethingafter))

	gw.RegisterHandler("transfer", horizon.Adapt(h))

	spec := intake.Spec{}

	routes := []dt.Route[AppCtxI]{
		jsonRoute[TransferCommandHTTP]("transfer"),
	}

	fusion := dt.NewFusion[AppCtxI](ctx, spec, gw, routes)
	router := fusion.Build()

	httpSrv, err := dt.NewHTTPServer(
		":8080", router,
		2*time.Second,  // ReadHeaderTimeout
		5*time.Second,  // ReadTimeout
		10*time.Second, // WriteTimeout
		60*time.Second, // IdleTimeout
	)
	if err != nil {
		log.Fatal(fmt.Errorf("http server init: %w", err))
	}

	// Context + signals for graceful shutdown
	httpCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start
	go func() {
		if err := httpSrv.Start(httpCtx); err != nil {
			log.Fatal(fmt.Errorf("running server: %w", err)) // platform.Field{Key: "component", Value: "http"})
			stop()
		}
	}()

	<-httpCtx.Done()

	// Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)

	// if handler, ok := gw.Handler("transfer"); ok {
	// 	_, _ = handler(
	// 		ctx,
	// 		inbound.RequestMeta{ClientID: "cli-1", RequestID: "req-123", Protocol: "cli", Target: "transfer"},
	// 		TransferCommand{fromAccount: "A", toAccount: "B", amountCents: 100, idempotencyKey: "asedfasdf"},
	// 	)
	// }
}

func printbefore(ctx AppCtxI, meta inbound.RequestMeta, cmd TransferCommand) {
	fmt.Println("before use case!, do something")
}

func dosomethingafter(ctx AppCtxI, meta inbound.RequestMeta, cmd TransferCommand, res TransferResult, err error) {
	fmt.Println(" is error nil ? ", err)
}

func jsonRoute[T any](key string) dt.Route[AppCtxI] {
	return dt.JSONRoute[AppCtxI, T](horizon.HandlerKey(key))
}

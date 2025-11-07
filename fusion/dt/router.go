package dt

import (
	"net/http"

	"hexa/m/v2/fusion/intake"
	"hexa/m/v2/horizon/ports/inbound"
)

type Path string

func (p Path) String() string {
	return string(p)
}

// DTFusion Deuterium–Tritium fusion
type DTFusion struct {
	handlers map[Path]func(http.ResponseWriter, *http.Request)
	mux      *http.ServeMux
}

// NewFusion creates and returns a new HTTP router for the API Gateway.
func NewFusion[c inbound.Ctx](
	handlers map[Path]func(http.ResponseWriter, *http.Request),
) DTFusion {
	mux := http.NewServeMux()
	return DTFusion{handlers, mux}
}

func (dt DTFusion) RegisterFuncs() {
	for k, v := range dt.handlers {
		dt.mux.HandleFunc(k.String(), v)
	}
}

func (dt DTFusion) ApplyPolicy(spec intake.Spec) http.Handler {
	mws := []func(next http.Handler) http.Handler{}

	if spec.MaxBodyBytes > 0 {
		mws = append(mws, LimitBytes(spec.MaxBodyBytes))
	}
	if spec.MaxInFlight > 0 {
		mws = append(mws, MaxInFlight(spec.MaxInFlight))
	}
	if spec.IPRateLimiter != nil {
		//		mws = append(mws, RateLimitHTTP(spec.IPRateLimiter))
	}
	if spec.EnableReqID {
		mws = append(mws, RequestID)
	}
	if spec.EnableReqLog {
		//		mws = append(mws, RequestLogger(fusion.DefaultLogger))
	}
	if spec.EnableRecover {
		// mws = append(mws, RecovererWithLogger(fusion.DefaultLogger))
	}
	return chain(dt.mux, mws...)
}

// func thiswasoldNewRouterCode() {
// 	mux.HandleFunc("POST /transfer",
// 		Unary[inbound.TransferCommand, inbound.TransferResult](
// 			gw.TransferHandler,
// 			TransferJSONDecoder(),
// 			func(w http.ResponseWriter, res inbound.TransferResult) {
// 				writer.JSON(w, http.StatusOK, contracts.TransferResponse{
// 					TransactionID: res.TransactionID().String(),
// 					Status:        res.Status().String(),
// 					Message:       res.Message(),
// 				})
// 			},
// 			DefaultMeta,
// 		),
// 	)
//
// 	mux.HandleFunc("GET /metrics",
// 		Unary[struct{}, contracts.MetricsSnapshot](
// 			gw.MetricsHandler, // ports.UnaryHandler[struct{}, types.MetricsSnapshot]
// 			EmptyDecoder,
// 			JSONEncoder[contracts.MetricsSnapshot](http.StatusOK),
// 			DefaultMeta,
// 		),
// 	)
//
// 	mux.HandleFunc("GET /healthz",
// 		Unary[struct{}, map[string]string](
// 			gw.HealthHandler,
// 			EmptyDecoder,
// 			JSONEncoder[map[string]string](http.StatusOK),
// 			DefaultMeta,
// 		),
// 	)
//
// 	// pprof
// 	mux.HandleFunc("/debug/pprof/", pprof.Index)
// 	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
// 	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
// 	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
// 	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
//
// 	lightLimiter := limiter.NewLightLimiter(10, 20)
//
// 	return middleware.Chain(
// 		mux,
// 		middleware.RecovererWithLogger(logger),
// 		middleware.RequestID,
// 		middleware.RequestLogger(logger),
// 		middleware.RateLimitHTTP(lightLimiter.Allow),
// 		middleware.MaxInFlight(1024),
// 		middleware.LimitBytes(1<<20),
// 	)
// }

// // TransferJSONDecoder decodes a TransferCommand from a JSON HTTP request.
// func TransferJSONDecoder() Decoder[inbound.TransferCommand] {
// 	return func(r *http.Request) (inbound.TransferCommand, error) {
// 		var dto struct {
// 			From           string `json:"from"`
// 			To             string `json:"to"`
// 			Amount         int64  `json:"amount"`
// 			IdempotencyKey string `json:"idempotency_key"`
// 		}
// 		dec := json.NewDecoder(r.Body)
// 		dec.DisallowUnknownFields()
// 		if err := dec.Decode(&dto); err != nil {
// 			var maxErr *http.MaxBytesError
// 			if errors.As(err, &maxErr) {
// 				return inbound.TransferCommand{}, apperr.PayloadTooLarge("request body too large")
// 			}
// 			return inbound.TransferCommand{}, apperr.Invalid("invalid JSON payload")
// 		}
// 		return inbound.NewTransferCommand(
// 			dto.From,
// 			dto.To,
// 			dto.Amount,
// 			dto.IdempotencyKey,
// 		), nil
// 	}
// }

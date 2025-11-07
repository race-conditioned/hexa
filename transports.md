| Reaction Type                  | Elements                                             | Traits                                   | Thematic Mapping                                                   |
| ------------------------------ | ---------------------------------------------------- | ---------------------------------------- | ------------------------------------------------------------------ |
| **Deuterium–Tritium (D–T)**    | Heavy hydrogen isotopes                              | Most common, powerful, produces neutrons | 🔥 HTTP — widespread, mature, slightly “messy” (headers, overhead) |
| **Deuterium–Helium-3 (D–He³)** | Cleaner, less radioactive, rarer fuel                | Efficient, elegant                       | ⚡ GraphQL — more structured, less noise                           |
| **Proton–Boron-11 (p–B¹¹)**    | Aneutronic (no radiation), requires high temperature | Super clean but complex                  | 🧠 gRPC — high-performance, efficient, precise                     |
| **Deuterium–Deuterium (D–D)**  | Simpler, lower yield                                 | Older and rougher                        | 🌐 maybe RESTful legacy APIs                                       |
| **Proton–Proton chain**        | What powers stars                                    | Gentle but constant                      | ☀️ Internal service comms, lightweight protocols                   |

| Protocol                          | Reaction Name                    | Why It Fits                                                                      |
| --------------------------------- | -------------------------------- | -------------------------------------------------------------------------------- |
| **HTTP / REST**                   | **D–T Fusion** → `fusion/dt`     | It’s the workhorse, standard but “hot” and lossy — lots of power, some overhead. |
| **GraphQL**                       | **D–He3 Fusion** → `fusion/dhe3` | Cleaner, more efficient data flow.                                               |
| **gRPC**                          | **p–B11 Fusion** → `fusion/pb11` | Elegant, high-performance, but more complex and specialized.                     |
| **WebSocket / Streaming**         | **D–D Fusion** → `fusion/dd`     | Simpler, older, steady data stream.                                              |
| **Internal RPC or control plane** | **p–p Chain** → `fusion/pp`      | Low-energy, steady, internal power.                                              |

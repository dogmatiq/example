// Package idletimeout manages engine lifetime for the public demo hosted at
// https://example.dogmatiq.io. It is hosting infrastructure, separate from the
// illustrative application and HTTP server setup in cmd/bank.
//
// [Inject] adds idle shutdown to an otherwise ordinary application setup. It
// replaces the engine runner and command executor, wraps the HTTP handler to
// track request activity, and configures database pools to expire idle
// connections. The engine stops after an idle period, and the next request
// starts a fresh engine after shutdown completes.
//
// Database pools and the UI handler remain owned by the caller. Expiring idle
// connections prevents their TCP traffic from keeping the platform awake.
// Platform sleeping must be enabled separately. Scheduled work waits until
// the engine wakes. The HTTP server handles liveness probes separately from
// the UI requests that keep the engine awake.
package idletimeout

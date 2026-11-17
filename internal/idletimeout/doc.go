// Package idletimeout manages engine lifetime for the public demo hosted at
// https://example.dogmatiq.io. It is hosting infrastructure, separate from the
// illustrative application and HTTP server setup in cmd/bank.
//
// Manager forwards commands to the current engine. Its HTTP middleware also
// tracks requests that only read projections. The engine stops after an idle
// period, and the next request starts a fresh engine after shutdown completes.
//
// Database pools and the UI handler remain owned by the caller. Pools should
// expire idle connections so their TCP traffic does not prevent platform sleep.
// Platform sleeping must be enabled separately. Scheduled work waits until the
// engine wakes; HTTP health checks and other probes also keep it awake.
package idletimeout

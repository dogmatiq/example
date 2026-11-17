// Package publicdemo bootstraps the public demo hosted at
// https://example.dogmatiq.io. The illustrative application setup remains in
// cmd/bank; this package contains only the hosted startup path.
//
// main selects this path when DOGMA_EXAMPLE_IDLE_TIMEOUT is set to a positive
// Go duration, such as "1m". After that much time without active HTTP requests,
// the engine stops and both database pools close. The HTTP listener stays open,
// and the next request starts the application again using its persisted state.
// HTTP health checks and other probes also count as activity.
//
// This lets the hosting platform put the container to sleep when background
// traffic stops. Platform sleeping must be enabled separately. The next visitor
// may experience a cold start, and scheduled work waits until the app wakes.
package publicdemo

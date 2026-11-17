// Package idletimeout manages idle shutdown for the public demo.
// It is not part of the Dogma application or an example of engine usage.
//
// Set DOGMA_EXAMPLE_IDLE_TIMEOUT to a positive Go duration, such as "1m", to
// stop the application after that much time without HTTP requests. The HTTP
// listener stays open; the next request starts a new application instance.
// Active requests keep the application awake, and the idle period begins when
// the last request finishes. HTTP health checks and other probes also count as
// activity.
// Leave the variable unset to keep the application running continuously.
//
// Stopping the application closes its database pools and stops the engine's
// background traffic, allowing the hosting platform to put the container to
// sleep. The HTTP server should also expire idle keep-alive connections.
// Platform sleeping must be enabled separately. The next visitor may experience
// a cold start, and scheduled work becomes overdue until the application wakes.
package idletimeout

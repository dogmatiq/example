// Package idletimeout coordinates a cancelable application with HTTP activity.
// It is hosting infrastructure used by publicdemo, not part of the illustrative
// Dogma application setup.
//
// Requests wait for application initialization and keep it running until they
// finish. An idle application is canceled, and a subsequent request waits for
// its cleanup before starting it again. The HTTP server remains running and
// should expire idle keep-alive connections so they do not prevent sleeping.
package idletimeout

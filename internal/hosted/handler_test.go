package hosted

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestIdleLifecycle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var starts, stops atomic.Int32

		handler := newTestHandler(t, time.Second, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			starts.Add(1)
			defer stops.Add(1)

			return serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if _, ok := request.Context().Deadline(); !ok {
					t.Error("request has no deadline")
				}

				w.WriteHeader(http.StatusNoContent)
			}))
		})

		request(t, handler)

		time.Sleep(time.Second)
		synctest.Wait()

		if starts.Load() != 1 || stops.Load() != 1 {
			t.Fatalf("idle application: starts=%d, stops=%d", starts.Load(), stops.Load())
		}

		// Concurrent visitors must share a single replacement instance.
		var visitors sync.WaitGroup

		for range 10 {
			visitors.Go(func() {
				request(t, handler)
			})
		}

		visitors.Wait()

		if starts.Load() != 2 || stops.Load() != 1 {
			t.Fatalf("after waking: starts=%d, stops=%d", starts.Load(), stops.Load())
		}
	})
}

func TestRequestsExtendIdlePeriod(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var stops atomic.Int32

		handler := newTestHandler(t, time.Second, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			defer stops.Add(1)

			return serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
		})

		request(t, handler)
		time.Sleep(750 * time.Millisecond)

		request(t, handler)
		time.Sleep(750 * time.Millisecond)
		synctest.Wait()

		if stops.Load() != 0 {
			t.Fatal("application stopped despite recent traffic")
		}

		time.Sleep(250 * time.Millisecond)
		synctest.Wait()

		if stops.Load() != 1 {
			t.Fatal("application did not stop after the extended idle period")
		}
	})
}

func TestActiveRequestsPreventIdleShutdown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var stops atomic.Int32

		finish := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(finish)
		})

		handler := newTestHandler(t, time.Second, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			defer stops.Add(1)

			return serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				<-finish

				w.WriteHeader(http.StatusNoContent)
			}))
		})
		defer unblock()

		go request(t, handler)

		synctest.Wait()
		time.Sleep(2 * time.Second)
		synctest.Wait()

		if stops.Load() != 0 {
			t.Fatal("application stopped during an active request")
		}

		unblock()
		synctest.Wait()
		time.Sleep(time.Second)
		synctest.Wait()

		if stops.Load() != 1 {
			t.Fatal("application did not stop after the request finished")
		}
	})
}

func TestWakeWaitsForCleanup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var starts atomic.Int32

		cleanup := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(cleanup)
		})

		handler := newTestHandler(t, time.Second, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			first := starts.Add(1) == 1

			err := serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))

			if first {
				<-cleanup
			}

			return err
		})
		defer unblock()

		request(t, handler)
		time.Sleep(time.Second)
		synctest.Wait()

		go request(t, handler)

		synctest.Wait()

		if starts.Load() != 1 {
			t.Fatal("replacement started before the previous instance released its resources")
		}

		unblock()
		synctest.Wait()

		if starts.Load() != 2 {
			t.Fatal("request did not wake the application after cleanup")
		}
	})
}

func TestIdleShutdownDisabled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var stops atomic.Int32

		handler := newTestHandler(t, 0, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			defer stops.Add(1)

			return serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
		})

		time.Sleep(time.Hour)
		synctest.Wait()

		if stops.Load() != 0 {
			t.Fatal("application stopped with idle shutdown disabled")
		}

		request(t, handler)
	})
}

func TestApplicationFailureIsReported(t *testing.T) {
	for _, ready := range []bool{false, true} {
		t.Run(map[bool]string{false: "startup", true: "running"}[ready], func(t *testing.T) {
			ctx := t.Context()

			failure := errors.New("engine failed")
			handler := newHandler(ctx, testLogger(), func(
				ctx context.Context,
				_ *slog.Logger,
				serveHTTP func(context.Context, http.Handler) error,
			) error {
				if ready {
					serveCtx, stop := context.WithCancel(ctx)
					stop()

					_ = serveHTTP(serveCtx, http.NotFoundHandler())
				}

				return failure
			}, time.Second)
			defer handler.close()

			if err := handler.monitor(); !errors.Is(err, failure) {
				t.Fatalf("got %v, want engine failure", err)
			}
		})
	}
}

func TestApplicationFailureCancelsActiveRequests(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var stops atomic.Int32

		fail := make(chan struct{})
		handler := newTestHandler(t, time.Second, func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			defer stops.Add(1)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				select {
				case <-fail:
				case <-ctx.Done():
				}

				cancel()
			}()

			_ = serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				<-request.Context().Done()

				w.WriteHeader(http.StatusServiceUnavailable)
			}))

			return errors.New("engine failed")
		})

		response := httptest.NewRecorder()

		go handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

		synctest.Wait()
		close(fail)
		synctest.Wait()

		if response.Code != http.StatusServiceUnavailable || stops.Load() != 1 {
			t.Fatalf("failed application did not drain: status=%d, stops=%d", response.Code, stops.Load())
		}
	})
}

func TestShutdownDrainsActiveRequests(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var stops atomic.Int32

		finish := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(finish)
		})

		handler := newHandler(ctx, testLogger(), func(
			ctx context.Context,
			_ *slog.Logger,
			serveHTTP func(context.Context, http.Handler) error,
		) error {
			defer stops.Add(1)

			return serveHTTP(ctx, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				<-request.Context().Done()
				<-finish

				w.WriteHeader(http.StatusNoContent)
			}))
		}, time.Second)
		defer handler.close()
		defer unblock()

		go request(t, handler)

		synctest.Wait()
		cancel()
		synctest.Wait()

		if stops.Load() != 0 {
			t.Fatal("application released its resources during an active request")
		}

		unblock()
		synctest.Wait()

		if stops.Load() != 1 {
			t.Fatal("application did not release its resources after draining requests")
		}
	})
}

func TestSlowStartupDoesNotReturnSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		handler := newTestHandler(t, time.Minute, func(
			ctx context.Context,
			_ *slog.Logger,
			_ func(context.Context, http.Handler) error,
		) error {
			<-ctx.Done()

			return ctx.Err()
		})

		response := httptest.NewRecorder()

		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

		if response.Code != http.StatusServiceUnavailable {
			t.Fatalf("got status %d during slow startup", response.Code)
		}
	})
}

func TestInvalidIdleTimeout(t *testing.T) {
	for _, value := range []string{"invalid", "0", "-1s"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("DOGMA_EXAMPLE_IDLE_TIMEOUT", value)

			err := Run(context.Background(), testLogger(), nil)
			if err == nil {
				t.Fatal("invalid idle timeout was accepted")
			}
		})
	}
}

func newTestHandler(t *testing.T, idleTimeout time.Duration, run applicationFunc) *handler {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	handler := newHandler(ctx, testLogger(), run, idleTimeout)

	t.Cleanup(func() {
		cancel()
		handler.close()
	})

	go func() {
		_ = handler.monitor()
	}()

	return handler
}

func request(t *testing.T, handler http.Handler) {
	t.Helper()

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusNoContent {
		t.Errorf("got HTTP status %d, want %d", response.Code, http.StatusNoContent)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

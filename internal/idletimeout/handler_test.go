package idletimeout

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

func TestIdleAndConcurrentWake(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var starts, stops atomic.Int32
		var current int32

		controller, _ := testController(t, time.Second, func(ctx context.Context, ready func()) error {
			current = starts.Add(1)
			defer stops.Add(1)

			ready()
			<-ctx.Done()

			return ctx.Err()
		})
		handler := controller.wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if current == 0 {
				t.Error("request dispatched before initialization")
			}

			w.WriteHeader(http.StatusNoContent)
		}))

		hit(t, handler)
		time.Sleep(time.Second)
		synctest.Wait()

		if starts.Load() != 1 || stops.Load() != 1 {
			t.Fatalf("idle: starts=%d, stops=%d", starts.Load(), stops.Load())
		}

		var visitors sync.WaitGroup

		for range 10 {
			visitors.Go(func() {
				hit(t, handler)
			})
		}

		visitors.Wait()

		if starts.Load() != 2 || stops.Load() != 1 {
			t.Fatalf("wake: starts=%d, stops=%d", starts.Load(), stops.Load())
		}
	})
}

func TestRequestsKeepApplicationAwake(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		app := &testApp{}
		controller, _ := testController(t, time.Second, app.run)
		finish := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(finish)
		})
		defer unblock()

		handler := controller.wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			<-finish

			w.WriteHeader(http.StatusNoContent)
		}))

		go hit(t, handler)

		synctest.Wait()
		time.Sleep(2 * time.Second)
		synctest.Wait()

		if app.stops.Load() != 0 {
			t.Fatal("application stopped during an active request")
		}

		unblock()
		synctest.Wait()
		time.Sleep(750 * time.Millisecond)
		hit(t, handler)
		time.Sleep(750 * time.Millisecond)
		synctest.Wait()

		if app.stops.Load() != 0 {
			t.Fatal("application stopped despite recent traffic")
		}

		time.Sleep(250 * time.Millisecond)
		synctest.Wait()

		if app.stops.Load() != 1 {
			t.Fatal("application did not stop after traffic finished")
		}
	})
}

func TestWakeWaitsForCleanup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cleanup := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(cleanup)
		})
		defer unblock()

		app := &testApp{}
		app.cleanup = func() {
			if app.starts.Load() == 1 {
				<-cleanup
			}
		}
		controller, _ := testController(t, time.Second, app.run)
		handler := controller.wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		hit(t, handler)
		time.Sleep(time.Second)
		synctest.Wait()

		go hit(t, handler)

		synctest.Wait()

		if app.starts.Load() != 1 {
			t.Fatal("replacement started before cleanup finished")
		}

		unblock()
		synctest.Wait()

		if app.starts.Load() != 2 {
			t.Fatal("request did not wake the application after cleanup")
		}
	})
}

func TestIdleDisabled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		app := &testApp{}
		controller, _ := testController(t, 0, app.run)

		time.Sleep(time.Hour)
		synctest.Wait()

		if app.starts.Load() != 1 || app.stops.Load() != 0 {
			t.Fatal("application did not run continuously")
		}

		hit(t, controller.wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})))
	})
}

func TestStartupFailure(t *testing.T) {
	failure := errors.New("startup failed")
	_, result := testController(t, time.Second, func(context.Context, func()) error {
		return failure
	})

	if err := <-result; !errors.Is(err, failure) {
		t.Fatalf("got %v, want startup failure", err)
	}
}

func TestFailureCancelsRequests(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		failure := errors.New("engine failed")
		fail := make(chan struct{})
		controller, result := testController(t, time.Second, func(ctx context.Context, ready func()) error {
			ready()

			select {
			case <-fail:
				return failure
			case <-ctx.Done():
				return ctx.Err()
			}
		})
		handler := controller.wrap(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			<-request.Context().Done()

			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		response := httptest.NewRecorder()

		go handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

		synctest.Wait()
		close(fail)
		synctest.Wait()

		if err := <-result; !errors.Is(err, failure) {
			t.Fatalf("got %v, want engine failure", err)
		}

		if response.Code != http.StatusServiceUnavailable {
			t.Fatal("active request was not canceled")
		}
	})
}

func TestRequestDeadlineDuringStartup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		controller, _ := testController(t, time.Minute, func(ctx context.Context, _ func()) error {
			<-ctx.Done()

			return ctx.Err()
		})
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		response := httptest.NewRecorder()

		controller.wrap(http.NotFoundHandler()).ServeHTTP(
			response,
			httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx),
		)

		if response.Code != http.StatusServiceUnavailable {
			t.Fatalf("got status %d during startup", response.Code)
		}
	})
}

func TestInvalidIdleTimeout(t *testing.T) {
	for _, value := range []string{"invalid", "0", "-1s"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("DOGMA_EXAMPLE_IDLE_TIMEOUT", value)

			if err := Run(t.Context(), testLogger(), nil, nil); err == nil {
				t.Fatal("invalid idle timeout was accepted")
			}
		})
	}
}

type testApp struct {
	starts  atomic.Int32
	stops   atomic.Int32
	cleanup func()
}

func (app *testApp) run(ctx context.Context, ready func()) error {
	app.starts.Add(1)
	defer app.stops.Add(1)

	if app.cleanup != nil {
		defer app.cleanup()
	}

	ready()
	<-ctx.Done()

	return ctx.Err()
}

func testController(t *testing.T, timeout time.Duration, run func(context.Context, func()) error) (*controller, <-chan error) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	controller := &controller{ctx: ctx, enter: make(chan chan struct{}), leave: make(chan struct{})}
	result := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		result <- controller.run(testLogger(), timeout, run)

		cancel()
		close(done)
	}()

	t.Cleanup(func() {
		cancel()
		<-done
	})

	return controller, result
}

func hit(t *testing.T, handler http.Handler) {
	t.Helper()

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusNoContent {
		t.Errorf("got HTTP %d, want %d", response.Code, http.StatusNoContent)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

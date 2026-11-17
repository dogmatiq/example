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

	"github.com/dogmatiq/dogma"
)

func TestIdleAndConcurrentWake(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var starts, stops atomic.Int32

		manager, _ := testManager(t, time.Second, func(ctx context.Context) error {
			starts.Add(1)
			defer stops.Add(1)

			<-ctx.Done()

			return ctx.Err()
		})
		handler := manager.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
		synctest.Wait()

		if starts.Load() != 2 || stops.Load() != 1 {
			t.Fatalf("wake: starts=%d, stops=%d", starts.Load(), stops.Load())
		}
	})
}

func TestRequestsKeepApplicationAwake(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		app := &testApp{}
		manager, _ := testManager(t, time.Second, app.run)
		finish := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(finish)
		})
		defer unblock()

		handler := manager.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
		manager, _ := testManager(t, time.Second, app.run)
		handler := manager.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
		manager, _ := testManager(t, 0, app.run)

		time.Sleep(time.Hour)
		synctest.Wait()

		if app.starts.Load() != 1 || app.stops.Load() != 0 {
			t.Fatal("application did not run continuously")
		}

		hit(t, manager.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})))
	})
}

func TestStartupFailure(t *testing.T) {
	failure := errors.New("startup failed")
	_, result := testManager(t, time.Second, func(context.Context) error {
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
		manager, result := testManager(t, time.Second, func(ctx context.Context) error {

			select {
			case <-fail:
				return failure
			case <-ctx.Done():
				return ctx.Err()
			}
		})
		handler := manager.Wrap(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
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

func TestRequestDeadlineDuringShutdown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cleanup := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(cleanup)
		})
		defer unblock()

		app := &testApp{cleanup: func() {
			<-cleanup
		}}
		manager, _ := testManager(t, time.Second, app.run)
		handler := manager.Wrap(http.NotFoundHandler())

		synctest.Wait()
		time.Sleep(time.Second)
		synctest.Wait()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		response := httptest.NewRecorder()

		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx))

		if response.Code != http.StatusServiceUnavailable {
			t.Fatalf("got status %d while shutdown was pending", response.Code)
		}
	})
}

type testApp struct {
	starts  atomic.Int32
	stops   atomic.Int32
	cleanup func()
}

func (app *testApp) run(ctx context.Context) error {
	app.starts.Add(1)
	defer app.stops.Add(1)

	if app.cleanup != nil {
		defer app.cleanup()
	}

	<-ctx.Done()

	return ctx.Err()
}

type testEngine struct {
	run     func(context.Context) error
	execute func(context.Context, dogma.Command, ...dogma.ExecuteCommandOption) error
}

func (engine *testEngine) Run(ctx context.Context) error {
	return engine.run(ctx)
}

func (engine *testEngine) ExecuteCommand(ctx context.Context, command dogma.Command, options ...dogma.ExecuteCommandOption) error {
	if engine.execute != nil {
		return engine.execute(ctx, command, options...)
	}

	return nil
}

func testManager(t *testing.T, timeout time.Duration, run func(context.Context) error) (*Manager, <-chan error) {
	t.Helper()

	return testManagerFactory(t, timeout, func() Engine {
		return &testEngine{run: run}
	})
}

func testManagerFactory(t *testing.T, timeout time.Duration, factory func() Engine) (*Manager, <-chan error) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	manager := New(ctx, testLogger(), timeout, factory)
	result := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		result <- manager.Run()

		cancel()
		close(done)
	}()

	t.Cleanup(func() {
		cancel()
		<-done
	})

	return manager, result
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

package idletimeout

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
)

func TestInjectedIdleShutdown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		app := &testApp{}
		var observed []int
		var replacements int

		makeEngine := func(generation int) *testEngine {
			return &testEngine{
				run: app.run,
				execute: func(context.Context, dogma.Command, ...dogma.ExecuteCommandOption) error {
					observed = append(observed, generation)

					return nil
				},
			}
		}

		engine := makeEngine(1)
		runEngine := func() error {
			return engine.Run(ctx)
		}
		var executor dogma.CommandExecutor = engine
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		if err := Inject(
			ctx,
			"1s",
			&runEngine,
			&executor,
			&handler,
			engine,
			func() Engine {
				replacements++

				return makeEngine(replacements + 1)
			},
			testLogger(),
		); err != nil {
			t.Fatal(err)
		}

		result := make(chan error, 1)
		done := make(chan struct{})

		go func() {
			result <- runEngine()

			close(done)
		}()

		t.Cleanup(func() {
			cancel()
			<-done
		})

		hit(t, handler)

		if err := executor.ExecuteCommand(ctx, &commands.OpenAccount{}); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		synctest.Wait()

		if app.starts.Load() != 1 || app.stops.Load() != 1 || replacements != 0 {
			t.Fatal("injected runner did not start and stop the original engine")
		}

		// A request that only reads projections must also wake the engine.
		hit(t, handler)
		synctest.Wait()

		if app.starts.Load() != 2 || replacements != 1 {
			t.Fatal("injected handler did not create a replacement engine")
		}

		if err := executor.ExecuteCommand(ctx, &commands.OpenAccount{}); err != nil {
			t.Fatal(err)
		}

		if len(observed) != 2 || observed[0] != 1 || observed[1] != 2 {
			t.Fatalf("injected executor used the wrong engines: %v", observed)
		}

		cancel()

		if err := <-result; !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v, want context cancellation", err)
		}

		if app.stops.Load() != 2 {
			t.Fatal("injected runner returned before the replacement stopped")
		}
	})
}

func TestInjectRejectsInvalidTimeout(t *testing.T) {
	for _, idle := range []string{"", "invalid", "0s", "-1s", "999999999999h"} {
		t.Run(idle, func(t *testing.T) {
			engine := &testEngine{}
			var runEngine func() error
			var executor dogma.CommandExecutor = engine
			originalHandler := &http.ServeMux{}
			var handler http.Handler = originalHandler

			err := Inject(
				t.Context(),
				idle,
				&runEngine,
				&executor,
				&handler,
				engine,
				func() Engine {
					t.Fatal("invalid configuration created an engine")

					return nil
				},
				testLogger(),
			)
			if err == nil {
				t.Fatal("invalid timeout was accepted")
			}

			if runEngine != nil || executor != engine || handler != originalHandler {
				t.Fatal("invalid timeout changed the application wiring")
			}
		})
	}
}

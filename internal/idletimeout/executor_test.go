package idletimeout

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

func TestExecutorUsesCurrentEngineAndPreservesOptions(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var generation atomic.Int32
		var observed []int32
		var observers int

		command := &commands.OpenAccount{}
		manager, _ := testManagerFactory(t, time.Second, func() Engine {
			id := generation.Add(1)

			return &testEngine{
				run: func(ctx context.Context) error {
					<-ctx.Done()

					return ctx.Err()
				},
				execute: func(ctx context.Context, got dogma.Command, options ...dogma.ExecuteCommandOption) error {
					if got != command || len(options) != 2 {
						t.Fatal("command or options were not forwarded")
					}

					if options[0].(dogma.IdempotencyKeyOption).Key() != "test-key" {
						t.Fatal("idempotency key was not preserved")
					}

					observed = append(observed, id)
					_, err := options[1].(dogma.EventObserverOption).Observer()(ctx, &events.AccountOpened{})

					return err
				},
			}
		})

		execute := func() {
			t.Helper()

			err := manager.ExecuteCommand(
				t.Context(),
				command,
				dogma.WithIdempotencyKey("test-key"),
				dogma.WithEventObserver(func(context.Context, *events.AccountOpened) (bool, error) {
					observers++

					return true, nil
				}),
			)
			if err != nil {
				t.Fatal(err)
			}
		}

		execute()
		time.Sleep(time.Second)
		synctest.Wait()
		execute()

		if len(observed) != 2 || observed[0] != 1 || observed[1] != 2 || observers != 2 {
			t.Fatalf("wrong engines or observers: engines=%v, observers=%d", observed, observers)
		}
	})
}

func TestCommandKeepsEngineAwake(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		app := &testApp{}
		finish := make(chan struct{})
		unblock := sync.OnceFunc(func() {
			close(finish)
		})
		defer unblock()

		manager, _ := testManagerFactory(t, time.Second, func() Engine {
			return &testEngine{
				run: app.run,
				execute: func(context.Context, dogma.Command, ...dogma.ExecuteCommandOption) error {
					<-finish

					return nil
				},
			}
		})

		go func() {
			if err := manager.ExecuteCommand(t.Context(), &commands.OpenAccount{}); err != nil {
				t.Error(err)
			}
		}()

		synctest.Wait()
		time.Sleep(2 * time.Second)
		synctest.Wait()

		if app.stops.Load() != 0 {
			t.Fatal("engine stopped while a command call was waiting")
		}

		unblock()
		synctest.Wait()
		time.Sleep(time.Second)
		synctest.Wait()

		if app.stops.Load() != 1 {
			t.Fatal("engine did not stop after the command call finished")
		}
	})
}

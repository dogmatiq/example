package publicdemo

import (
	"io"
	"log/slog"
	"testing"
)

func TestInvalidIdleTimeout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, value := range []string{"invalid", "0", "-1s"} {
		t.Run(value, func(t *testing.T) {
			if err := Run(t.Context(), logger, value); err == nil {
				t.Fatal("invalid idle timeout was accepted")
			}
		})
	}
}

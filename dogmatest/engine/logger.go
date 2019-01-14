package engine

// Logger is an interface for logging messages.
type Logger interface {
	Logf(f string, v ...interface{})
}

// SilentLogger is a Logger implementation that does not produce any output.
var SilentLogger silentLogger

type silentLogger struct{}

func (silentLogger) Logf(f string, v ...interface{}) {}

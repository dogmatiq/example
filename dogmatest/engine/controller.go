package engine

import "context"

// Controller orchestrates the handling of a message by a handler.
type Controller interface {
	Name() string
	Handler() interface{}
	Handle(ctx context.Context, logger Logger, env *Envelope) error
	Reset()
}

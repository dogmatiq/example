package types

type Controller interface {
	Name() string
	Handler() interface{}
	Handle(env Envelope) []Envelope
	Reset()
}

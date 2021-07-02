package asyncapi

// Decoder decodes an AsyncAPI JSON Schema document.
type Decoder interface {
	Decode([]byte, interface{}) error
}

// DecodeFunc is a helper func that implements the Decoder interface.
type DecodeFunc func([]byte, interface{}) error

func (d DecodeFunc) Decode(b []byte, dst interface{}) error {
	return d(b, dst)
}

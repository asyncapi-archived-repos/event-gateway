package asyncapi

// Decoder decodes an AsyncAPI document (several formats can be supported).
// See https://github.com/asyncapi/parser-go#overview for the minimum supported schemas.
type Decoder interface {
	Decode([]byte, interface{}) error
}

// DecodeFunc is a helper func that implements the Decoder interface.
type DecodeFunc func([]byte, interface{}) error

func (d DecodeFunc) Decode(b []byte, dst interface{}) error {
	return d(b, dst)
}

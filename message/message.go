package message

// Message represents a message flowing through the wire. For example, a Kafka message.
type Message struct {
	Context Context  `json:"context"`
	Key     []byte   `json:"key,omitempty"`
	Value   []byte   `json:"value,omitempty"`
	Headers []Header `json:"headers,omitempty"`
}

// Context contains information about the context that surrounds a message.
type Context struct {
	Channel string `json:"channel"`
}

// Header represents a header of a message, if there are any.
type Header struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

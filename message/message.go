package message

import (
	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
)

// The following constants are the keys on a watermill.Message Metadata where we store valuable and needed domain metadata.
// All contain the prefix `_asyncapi_eg_` so they can be unique-ish and human-readable.
// As a note: The term `eg` is a short version of Event-Gateway.
const (
	// MetadataChannel is the key used for storing the Channel in the message Metadata.
	MetadataChannel = "_asyncapi_eg_channel"
)

// New creates a new watermillmessage.Message. It injects the Channel as Metadata.
func New(payload []byte, channel string) *watermillmessage.Message {
	msg := watermillmessage.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set(MetadataChannel, channel)

	return msg
}

package message

import (
	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
)

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

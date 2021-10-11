package message

import (
	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
)

// New creates a new watermillmessage.Message. It injects the Channel as Metadata.
func New(payload []byte, channel string) *watermillmessage.Message {
	msg := watermillmessage.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set(MetadataChannel, channel)

	return msg
}

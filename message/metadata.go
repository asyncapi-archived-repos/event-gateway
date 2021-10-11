package message

import (
	"encoding/json"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
)

// The following constants are the keys on a watermill.Message Metadata where we store valuable and needed domain metadata.
// All contain the prefix `_asyncapi_eg_` so they can be unique-ish and human-readable.
// As a note: The term `eg` is a short version of Event-Gateway.
const (
	// MetadataChannel is the key used for storing the Channel in the message Metadata.
	MetadataChannel = "_asyncapi_eg_channel"

	// MetadataValidationError is the key used for storing the Validation Error if applies.
	MetadataValidationError = "_asyncapi_eg_validation_error"
)

// UnmarshalMetadata extracts a value from the Message Metadata and unmarshals it to the given object.
func UnmarshalMetadata(msg *watermillmessage.Message, key string, unmarshalTo interface{}) error {
	raw := msg.Metadata.Get(key)
	if raw == "" {
		return nil
	}

	return json.Unmarshal([]byte(raw), unmarshalTo)
}

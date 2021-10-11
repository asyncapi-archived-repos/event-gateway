package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Validator validates a message.
// In case the message is invalid, returns a ValidationError. The second returned value is an error during validating process.
type Validator func(*watermillmessage.Message) (*ValidationError, error)

// ValidationError represents a message validation error.
type ValidationError struct {
	Timestamp time.Time `json:"ts"`
	Errors    []string  `json:"errors"`
}

func (v ValidationError) Error() string {
	return strings.Join(v.Errors, " | ")
}

// NewValidationError creates a new ValidationError.
func NewValidationError(ts time.Time, errors ...string) *ValidationError {
	return &ValidationError{Timestamp: ts, Errors: errors}
}

// ValidationErrorFromMessage extracts a ValidationError from the message Metadata if exists.
func ValidationErrorFromMessage(msg *watermillmessage.Message) (*ValidationError, error) {
	validationErr := new(ValidationError)
	return validationErr, UnmarshalMetadata(msg, MetadataValidationError, validationErr)
}

// ValidationErrorToMessage sets a ValidationError to the given message Metadata.
func ValidationErrorToMessage(msg *watermillmessage.Message, validationErr *ValidationError) error {
	raw, err := json.Marshal(validationErr)
	if err != nil {
		return err
	}

	msg.Metadata.Set(MetadataValidationError, string(raw))

	return nil
}

// JSONSchemaMessageValidator validates a message payload based on a map of Json Schema, where the key can be any identifier  (depends on who implements it).
// For example, the identifier can be its channel name, message ID, etc.
func JSONSchemaMessageValidator(messageSchemas map[string]gojsonschema.JSONLoader, idProvider func(msg *watermillmessage.Message) string) (Validator, error) {
	return func(msg *watermillmessage.Message) (*ValidationError, error) {
		msgID := idProvider(msg)
		msgSchema, ok := messageSchemas[msgID]
		if !ok {
			return nil, nil
		}

		result, err := gojsonschema.Validate(msgSchema, gojsonschema.NewBytesLoader(msg.Payload))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error validating JSON Schema for message %s", msgID))
		}

		if !result.Valid() {
			errs := make([]string, len(result.Errors()))
			for i := 0; i < len(result.Errors()); i++ {
				errs[i] = result.Errors()[i].String()
			}

			return NewValidationError(time.Now(), errs...), nil
		}

		return nil, nil
	}, nil
}

package message

import (
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
	Timestamp time.Time                 `json:"ts"`
	Msg       *watermillmessage.Message `json:"msg"`
	Errors    []string                  `json:"errors"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("message %q is invalid. Validation errors: %s", v.Msg.UUID, strings.Join(v.Errors, " | "))
}

// NewValidationError creates a new ValidationError.
func NewValidationError(msg *watermillmessage.Message, ts time.Time, errors ...string) *ValidationError {
	return &ValidationError{Msg: msg, Timestamp: ts, Errors: errors}
}

// ValidationErrorNotifier notifies whenever a ValidationError happens.
type ValidationErrorNotifier func(validationError *ValidationError) error

// ValidationErrorToChanNotifier notifies to a given chan when a ValidationError happens.
func ValidationErrorToChanNotifier(errChan chan *ValidationError) ValidationErrorNotifier {
	return func(validationError *ValidationError) error {
		// TODO Blocking or non blocking? Shall we just fire and forget via goroutine instead?
		errChan <- validationError

		return nil
	}
}

// JSONSchemaMessageValidator validates a message payload based on a map of Json Schema, where the key can be any identifier  (depends on who implements it).
// For example, the identifier can be it's channel name, message ID, etc.
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

			return NewValidationError(msg, time.Now(), errs...), nil
		}

		return nil, nil
	}, nil
}

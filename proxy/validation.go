package proxy

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Message represents a message flowing through the wire. For example, a Kafka message.
type Message struct {
	Context MessageContext
	Key     []byte
	Value   []byte
	Headers []MessageHeader
}

// MessageContext contains information about the context that surrounds a message.
type MessageContext struct {
	Channel string
}

// MessageHeader represents a header of a message, if there are any.
type MessageHeader struct {
	Key   []byte
	Value []byte
}

// ValidationError represents a message validation error.
type ValidationError struct {
	Msg    *Message
	Result *gojsonschema.Result
}

func (v ValidationError) String() string {
	errs := make([]string, len(v.Result.Errors()))
	for i, err := range v.Result.Errors() {
		errs[i] = err.String()
	}

	return strings.Join(errs, " | ")
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

// MessageValidator validates a message.
// Returns a boolean indicating if the message is valid, and an error if something went wrong.
type MessageValidator func(*Message) (*ValidationError, error)

// NotifyOnValidationError is a MessageValidator that notifies ValidationError from a given MessageValidator output to the given channel.
func NotifyOnValidationError(validator MessageValidator, notifier ValidationErrorNotifier) MessageValidator {
	return func(msg *Message) (*ValidationError, error) {
		validationErr, err := validator(msg)
		if err != nil {
			return nil, err
		}

		if validationErr != nil {
			if err := notifier(validationErr); err != nil {
				return nil, errors.Wrap(err, "error notifying validation error")
			}

			return validationErr, nil
		}

		return nil, nil
	}
}

// JSONSchemaMessageValidator validates a message payload based on a map of Json Schema, where the key can be any identifier  (depends on who implements it).
// For example, the identifier can be it's channel name, message ID, etc.
func JSONSchemaMessageValidator(messageSchemas map[string]gojsonschema.JSONLoader, idProvider func(msg *Message) string) (MessageValidator, error) {
	return func(msg *Message) (*ValidationError, error) {
		msgID := idProvider(msg)
		msgSchema, ok := messageSchemas[msgID]
		if !ok {
			return nil, nil
		}

		result, err := gojsonschema.Validate(msgSchema, gojsonschema.NewBytesLoader(msg.Value))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error validating JSON Schema for message %s", msgID))
		}

		if !result.Valid() {
			return &ValidationError{
				Msg:    msg,
				Result: result,
			}, nil
		}

		return nil, nil
	}, nil
}

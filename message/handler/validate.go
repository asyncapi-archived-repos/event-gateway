package handler

import (
	"encoding/json"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/asyncapi/event-gateway/message"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ErrMessageIsInvalid is the error used when a message did not pass validation and failWhenInvalid option was set to true.
var ErrMessageIsInvalid = errors.New("Message is invalid and failWhenInvalid was set to true")

// ValidateMessage validates a message. If invalid, It injects the validation error into the message Metadata.
// By default, next handler will always be called, including whenever the message is invalid.
// In case you want to make it fail, set failWhenInvalid to true.
func ValidateMessage(validator message.Validator, failWhenInvalid bool) watermillmessage.HandlerFunc {
	return func(msg *watermillmessage.Message) ([]*watermillmessage.Message, error) {
		validationErr, err := validator(msg)
		if err != nil {
			return nil, err
		}

		if validationErr != nil {
			if failWhenInvalid {
				err = ErrMessageIsInvalid
			}

			rawErr, err := json.Marshal(validationErr)
			if err != nil {
				logrus.WithError(err).Error("Error marshaling message.ValidationError")
			}

			// Injecting the validation error into the message Metadata.
			msg.Metadata.Set(message.MetadataValidationError, string(rawErr))
		}

		return []*watermillmessage.Message{msg}, err
	}
}

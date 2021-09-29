package handler

import (
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/asyncapi/event-gateway/message"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ErrMessageIsInvalid is the error used when a message did not pass validation and failWhenInvalid option was set to true.
var ErrMessageIsInvalid = errors.New("Message is invalid and failWhenInvalid was set to true")

// ValidateMessage validates a message. Optionally notifies if a notifier is set.
// By default, next handler will always be called, including whenever the message is invalid. If you want to make it fail then, set failWhenInvalid to true.
func ValidateMessage(validator message.Validator, notifier message.ValidationErrorNotifier, failWhenInvalid bool) watermillmessage.HandlerFunc {
	return func(msg *watermillmessage.Message) ([]*watermillmessage.Message, error) {
		validationErr, err := validator(msg)
		if err != nil {
			return nil, err
		}

		if validationErr != nil {
			if notifier != nil {
				if err := notifier(validationErr); err != nil {
					logrus.WithError(err).Error("error notifying message validation error")
				}
			} else {
				logrus.WithError(validationErr).Error("Message is invalid")
			}

			if failWhenInvalid {
				return nil, ErrMessageIsInvalid
			}
		}

		return []*watermillmessage.Message{msg}, nil
	}
}

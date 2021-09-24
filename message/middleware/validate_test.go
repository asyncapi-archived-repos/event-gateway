package middleware

import (
	"testing"
	"time"

	"github.com/asyncapi/event-gateway/message"
	"github.com/stretchr/testify/assert"
)

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name            string
		validator       message.Validator
		notifier        message.ValidationErrorNotifier
		failWhenInvalid bool
		expectedErr     error
	}{
		{
			name: "Message is valid",
			validator: func(*message.Message) (*message.ValidationError, error) {
				return nil, nil
			},
		},
		{
			name:      "Message is invalid. No notifier is set. failWhenInvalid = false",
			validator: invalidMessageValidator,
		},
		{
			name:            "Message is invalid. No notifier is set. failWhenInvalid = true",
			failWhenInvalid: true,
			expectedErr:     ErrMessageIsInvalid,
			validator:       invalidMessageValidator,
		},
		{
			name: "Message is invalid. Notifier is set.",
			notifier: func(validationError *message.ValidationError) error {
				assert.NotNil(t, validationError)
				return nil
			},
			validator: invalidMessageValidator,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := &message.Message{}
			returnedMsg, err := ValidateMessage(test.validator, test.notifier, test.failWhenInvalid)(msg)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
				assert.Nil(t, returnedMsg)
			} else {
				assert.NoError(t, err)
				assert.Same(t, msg, returnedMsg)
			}
		})
	}
}

func invalidMessageValidator(m *message.Message) (*message.ValidationError, error) {
	return &message.ValidationError{
		Timestamp: time.Now(),
		Msg:       m,
		Errors:    []string{"testing error"},
	}, nil
}

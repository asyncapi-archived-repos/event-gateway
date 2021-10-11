package handler

import (
	"testing"
	"time"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/asyncapi/event-gateway/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name                  string
		validator             message.Validator
		failWhenInvalid       bool
		expectedValidationErr string
		expectedErr           error
	}{
		{
			name: "Message is valid",
			validator: func(*watermillmessage.Message) (*message.ValidationError, error) {
				return nil, nil
			},
		},
		{
			name:                  "Message is invalid. failWhenInvalid = false",
			validator:             invalidMessageValidator,
			expectedValidationErr: "testing error",
		},
		{
			name:                  "Message is invalid.failWhenInvalid = true",
			validator:             invalidMessageValidator,
			failWhenInvalid:       true,
			expectedValidationErr: "testing error",
			expectedErr:           ErrMessageIsInvalid,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := message.New([]byte{}, test.name)
			returnedMsgs, err := ValidateMessage(test.validator, test.failWhenInvalid)(msg)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, returnedMsgs, 1)
			assert.Same(t, msg, returnedMsgs[0])

			validationErr, err := message.ValidationErrorFromMessage(msg)
			require.NoError(t, err)

			if test.expectedValidationErr != "" {
				assert.Equal(t, test.expectedValidationErr, validationErr.Error())
			}
		})
	}
}

func invalidMessageValidator(_ *watermillmessage.Message) (*message.ValidationError, error) {
	return &message.ValidationError{
		Timestamp: time.Now(),
		Errors:    []string{"testing error"},
	}, nil
}

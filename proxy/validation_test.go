package proxy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

func TestNotifyOnValidationError(t *testing.T) {
	expectedMessage := generateTestMessage()
	validator := func(msg *Message) (*ValidationError, error) {
		assert.Equal(t, expectedMessage, msg)
		return NewValidationError(msg, time.Now(), "This is a validation error"), nil
	}

	var notified bool
	notifier := func(validationError *ValidationError) error {
		notified = true
		return nil
	}

	validationErr, err := NotifyOnValidationError(validator, notifier)(expectedMessage)
	assert.NoError(t, err)
	assert.NotEmpty(t, validationErr.Errors)
	assert.True(t, notified)
}

func generateTestMessage() *Message {
	return &Message{
		Context: MessageContext{Channel: "test"},
		Value:   []byte(`Hello World!`),
	}
}

func TestJsonSchemaMessageValidator(t *testing.T) {
	schema := `{
   "properties":{
      "command":{
         "description":"Whether to turn on or off the light.",
         "enum":[
            "on",
            "off"
         ],
         "type":"string"
      }
   },
   "type":"object"
}`

	tests := []struct {
		name    string
		valid   bool
		payload []byte
	}{
		{
			name:    "Valid payload",
			payload: []byte(`{"command": "on"}`),
			valid:   true,
		},
		{
			name:    "Invalid payload",
			payload: []byte(`{"command": 123123}`),
			valid:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedMsg := generateTestMessage()
			expectedMsg.Value = test.payload

			idToSchemaMap := map[string]gojsonschema.JSONLoader{
				expectedMsg.Context.Channel: gojsonschema.NewStringLoader(schema),
			}

			validator, err := JSONSchemaMessageValidator(idToSchemaMap, func(msg *Message) string {
				assert.Equal(t, expectedMsg, msg)
				return msg.Context.Channel
			})
			assert.NoError(t, err)

			validationErr, err := validator(expectedMsg)
			assert.NoError(t, err)

			if test.valid {
				assert.Nil(t, validationErr)
			} else {
				assert.NotNil(t, validationErr)
				assert.NotEmpty(t, validationErr.Errors)
			}
		})
	}
}

package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

func TestValidationError_String(t *testing.T) {
	validationErr := generateTestValidationError(nil)
	assert.Equal(t, "AnIntegerField: Invalid type. Expected: integer, given: string | AStringField: Invalid type. Expected: string, given: integer", validationErr.String())
}

func TestNotifyOnValidationError(t *testing.T) {
	expectedMessage := generateTestMessage()
	validator := func(msg *Message) (*ValidationError, error) {
		assert.Equal(t, expectedMessage, msg)
		return generateTestValidationError(msg), nil
	}

	var notified bool
	notifier := func(validationError *ValidationError) error {
		notified = true
		return nil
	}

	validationErr, err := NotifyOnValidationError(validator, notifier)(expectedMessage)
	assert.NoError(t, err)
	assert.False(t, validationErr.Result.Valid())
	assert.True(t, notified)
}

func generateTestMessage() *Message {
	return &Message{
		Context: MessageContext{Channel: "test"},
		Value:   []byte(`Hello World!`),
	}
}

func generateTestValidationError(msg *Message) *ValidationError {
	validationErr := &ValidationError{
		Msg:    msg,
		Result: &gojsonschema.Result{},
	}

	addTestErrors(validationErr)
	return validationErr
}

func addTestErrors(validationErr *ValidationError) {
	badTypeErr := &gojsonschema.InvalidTypeError{}
	badTypeErr.SetContext(gojsonschema.NewJsonContext("AnIntegerField", nil))
	badTypeErr.SetDetails(gojsonschema.ErrorDetails{
		"expected": gojsonschema.TYPE_INTEGER,
		"given":    gojsonschema.TYPE_STRING,
	})
	badTypeErr.SetDescriptionFormat(gojsonschema.Locale.InvalidType())
	validationErr.Result.AddError(badTypeErr, badTypeErr.Details())

	badTypeErr2 := &gojsonschema.InvalidTypeError{}
	badTypeErr2.SetContext(gojsonschema.NewJsonContext("AStringField", nil))
	badTypeErr2.SetDetails(gojsonschema.ErrorDetails{
		"expected": gojsonschema.TYPE_STRING,
		"given":    gojsonschema.TYPE_INTEGER,
	})
	badTypeErr2.SetDescriptionFormat(gojsonschema.Locale.InvalidType())
	validationErr.Result.AddError(badTypeErr2, badTypeErr2.Details())
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
				assert.False(t, validationErr.Result.Valid())
			}
		})
	}
}

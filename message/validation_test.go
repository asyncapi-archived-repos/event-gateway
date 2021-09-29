package message

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

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
			expectedMsg.Payload = test.payload

			idToSchemaMap := map[string]gojsonschema.JSONLoader{
				expectedMsg.Metadata.Get(MetadataChannel): gojsonschema.NewStringLoader(schema),
			}

			validator, err := JSONSchemaMessageValidator(idToSchemaMap, func(msg *watermillmessage.Message) string {
				assert.Equal(t, expectedMsg, msg)
				return msg.Metadata.Get(MetadataChannel)
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

func generateTestMessage() *watermillmessage.Message {
	msg := watermillmessage.NewMessage(watermill.NewUUID(), []byte(`Hello World!`))
	msg.Metadata.Set(MetadataChannel, "the-channel")
	return msg
}

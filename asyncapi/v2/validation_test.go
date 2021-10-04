package v2

import (
	"testing"

	"github.com/asyncapi/event-gateway/asyncapi"
	"github.com/asyncapi/event-gateway/message"
	"github.com/stretchr/testify/assert"
)

func TestFromDocJsonSchemaMessageValidator(t *testing.T) {
	tests := []struct {
		name    string
		valid   bool
		schema  *Schema
		payload []byte
	}{
		{
			name: "Valid payload",
			schema: &Schema{
				TypeField: "object",
				PropertiesField: Schemas{
					"AnIntergerField": &Schema{
						MaximumField:  refFloat64(10),
						MinimumField:  refFloat64(3),
						RequiredField: []string{"AnIntergerField"},
						TypeField:     "number",
					},
				},
			},
			payload: []byte(`{"AnIntergerField": 5}`),
			valid:   true,
		},
		{
			name: "Valid multiple payloads",
			schema: &Schema{
				TypeField: "object",
				OneOfField: []asyncapi.Schema{
					&Schema{
						PropertiesField: Schemas{
							"AnIntergerField": &Schema{
								MaximumField:  refFloat64(10),
								MinimumField:  refFloat64(3),
								RequiredField: []string{"AnIntergerField"},
								TypeField:     "number",
							},
						},
						AdditionalPropertiesField: false,
					},
					&Schema{
						PropertiesField: Schemas{
							"AStringField": &Schema{
								RequiredField: []string{"AStringField"},
								TypeField:     "string",
							},
						},
						AdditionalPropertiesField: false,
					},
				},
			},
			payload: []byte(`{"AStringField": "hello!"}`),
			valid:   true,
		},
		{
			name: "Invalid payload",
			schema: &Schema{
				TypeField: "object",
				PropertiesField: Schemas{
					"AnIntergerField": &Schema{
						MaximumField:  refFloat64(10),
						MinimumField:  refFloat64(3),
						RequiredField: []string{"AnIntergerField"},
						TypeField:     "number",
					},
				},
			},
			payload: []byte(`{"AnIntergerField": 1}`),
			valid:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Doc generation
			channel := NewChannel(t.Name())
			channel.Publish = NewPublishOperation(&Message{PayloadField: test.schema})
			doc := Document{ChannelsField: map[string]Channel{t.Name(): *channel}}

			// Test
			validator, err := FromDocJSONSchemaMessageValidator(doc)
			assert.NoError(t, err)

			msg := message.New(test.payload, t.Name())
			validationErr, err := validator(msg)
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

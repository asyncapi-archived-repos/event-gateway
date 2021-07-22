package v2

import (
	"testing"

	"github.com/asyncapi/event-gateway/proxy"
	"github.com/stretchr/testify/assert"
)

func TestFromDocJsonSchemaMessageValidator(t *testing.T) {
	msg := &Message{
		PayloadField: &Schema{
			TypeField: "object",
			PropertiesField: Schemas{
				"AnIntergerField": &Schema{
					Extendable:    Extendable{},
					MaximumField:  refFloat64(10),
					MinimumField:  refFloat64(3),
					RequiredField: []string{"AnIntergerField"},
					TypeField:     "number",
				},
			},
		},
	}
	channel := NewChannel("test")
	channel.Subscribe = NewSubscribeOperation(msg)

	doc := Document{
		Extendable: Extendable{},
		ChannelsField: map[string]Channel{
			"test": *channel,
		},
	}

	tests := []struct {
		name    string
		valid   bool
		payload []byte
	}{
		{
			name:    "Valid payload",
			payload: []byte(`{"AnIntergerField": 5}`),
			valid:   true,
		},
		{
			name:    "Invalid payload",
			payload: []byte(`{"AnIntergerField": 1}`),
			valid:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator, err := FromDocJSONSchemaMessageValidator(doc)
			assert.NoError(t, err)

			msg := &proxy.Message{
				Context: proxy.MessageContext{
					Channel: "test",
				},
				Value: test.payload,
			}
			validationErr, err := validator(msg)
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

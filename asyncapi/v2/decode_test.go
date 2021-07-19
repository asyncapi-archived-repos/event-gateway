package v2

import (
	"strings"
	"testing"

	"github.com/asyncapi/event-gateway/asyncapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeFromFile(t *testing.T) {
	doc := new(Document)
	require.NoError(t, Decode([]byte("testdata/example-kafka.yaml"), doc))

	require.Len(t, doc.Servers(), 1)
	s, ok := doc.Server("test")
	assert.True(t, ok)
	assert.Equal(t, s, doc.Servers()[0])

	assert.True(t, s.HasName())
	assert.Equal(t, "test", s.Name())
	assert.True(t, s.HasDescription())
	assert.Equal(t, "Test broker", s.Description())
	assert.True(t, s.HasProtocol())
	assert.Equal(t, "kafka-secure", s.Protocol())
	assert.True(t, s.HasURL())
	assert.Equal(t, "localhost:9092", s.URL())
	assert.True(t, s.HasExtension(asyncapi.ExtensionEventGatewayListener))
	assert.Equal(t, "proxy:28002", s.Extension(asyncapi.ExtensionEventGatewayListener))
	assert.True(t, s.HasExtension(asyncapi.ExtensionEventGatewayDialMapping))
	assert.Equal(t, "broker:9092", s.Extension(asyncapi.ExtensionEventGatewayDialMapping))
	assert.Empty(t, s.Variables())

	channelPaths := []string{
		"smartylighting.streetlights.1.0.event.{streetlightId}.lighting.measured",
		"smartylighting.streetlights.1.0.action.{streetlightId}.turn.on",
		"smartylighting.streetlights.1.0.action.{streetlightId}.turn.off",
		"smartylighting.streetlights.1.0.action.{streetlightId}.dim",
	}

	require.Len(t, doc.Channels(), 4)
	for _, c := range doc.Channels() {
		assert.Containsf(t, channelPaths, c.Path(), "Channel path %q is not one of: %s", c.Path(), strings.Join(channelPaths, ", "))
		assert.Len(t, c.Parameters(), 1)
		assert.Len(t, c.Operations(), 1)
		assert.Len(t, c.Messages(), 1)
		for _, o := range c.Operations() {
			assert.Containsf(t, []asyncapi.OperationType{OperationTypePublish, OperationTypeSubscribe}, o.Type(), "Operation type %q is not one of %s, %s", o.Type(), OperationTypePublish, OperationTypeSubscribe)
		}

		for _, m := range c.Messages() {
			assert.NotNil(t, m.Payload())
		}
	}
}

//nolint:misspell,funlen
func TestDecodeFromPlainText(t *testing.T) {
	raw := []byte(`
asyncapi: '2.0.0'
info:
  title: Streetlights API
  version: '1.0.0'
  description: |
    The Smartylighting Streetlights API allows you
    to remotely manage the city lights.
  license:
    name: Apache 2.0
    url: 'https://www.apache.org/licenses/LICENSE-2.0'
servers:
  mosquitto:
    url: mqtt://test.mosquitto.org
    protocol: mqtt
channels:
  light/measured:
    publish:
      summary: Inform about environmental lighting conditions for a particular streetlight.
      operationId: onLightMeasured
      message:
        name: LightMeasured
        payload:
          type: object
          properties:
            id:
              type: integer
              minimum: 0
              description: Id of the streetlight.
            lumens:
              type: integer
              minimum: 0
              description: Light intensity measured in lumens.
            sentAt:
              type: string
              format: date-time
              description: Date and time when the message was sent.`)

	doc := new(Document)
	require.NoError(t, Decode(raw, doc))

	require.Len(t, doc.Servers(), 1)
	s, ok := doc.Server("mosquitto")
	assert.True(t, ok)
	assert.Equal(t, s, doc.Servers()[0])

	assert.True(t, s.HasName())
	assert.Equal(t, "mosquitto", s.Name())
	assert.False(t, s.HasDescription())
	assert.True(t, s.HasProtocol())
	assert.Equal(t, "mqtt", s.Protocol())
	assert.True(t, s.HasURL())
	assert.Equal(t, "mqtt://test.mosquitto.org", s.URL())
	assert.False(t, s.HasExtension(asyncapi.ExtensionEventGatewayListener))
	assert.False(t, s.HasExtension(asyncapi.ExtensionEventGatewayDialMapping))
	assert.Empty(t, s.Variables())

	assert.Len(t, doc.ApplicationPublishableChannels(), 1)
	assert.Len(t, doc.ApplicationPublishOperations(), 1)
	assert.Len(t, doc.ApplicationPublishableMessages(), 1)

	assert.Empty(t, doc.ApplicationSubscribableChannels())
	assert.Empty(t, doc.ApplicationSubscribeOperations())
	assert.Empty(t, doc.ApplicationSubscribableMessages())

	assert.Len(t, doc.ClientSubscribableChannels(), 1)
	assert.Len(t, doc.ClientSubscribeOperations(), 1)
	assert.Len(t, doc.ClientSubscribableMessages(), 1)

	assert.Empty(t, doc.ClientPublishableChannels())
	assert.Empty(t, doc.ClientPublishOperations())
	assert.Empty(t, doc.ClientPublishableMessages())

	channels := doc.Channels()
	require.Len(t, channels, 1)
	assert.Equal(t, "light/measured", channels[0].Path())
	assert.Empty(t, channels[0].Parameters())
	assert.False(t, channels[0].HasDescription())

	operations := channels[0].Operations()
	require.Len(t, operations, 1)
	assert.Equal(t, OperationTypePublish, operations[0].Type())
	assert.True(t, operations[0].IsApplicationPublishing())
	assert.False(t, operations[0].IsApplicationSubscribing())
	assert.True(t, operations[0].IsClientSubscribing())
	assert.False(t, operations[0].IsClientPublishing())
	assert.False(t, operations[0].HasDescription())
	assert.True(t, operations[0].HasSummary())
	assert.Equal(t, "Inform about environmental lighting conditions for a particular streetlight.", operations[0].Summary())
	assert.Equal(t, "onLightMeasured", operations[0].ID())

	messages := operations[0].Messages()
	require.Len(t, messages, 1)

	assert.Equal(t, "LightMeasured", messages[0].Name())
	assert.False(t, messages[0].HasSummary())
	assert.False(t, messages[0].HasDescription())
	assert.False(t, messages[0].HasTitle())
	assert.Empty(t, messages[0].ContentType())

	payload := messages[0].Payload()
	require.NotNil(t, payload)

	assert.Equal(t, []string{"object"}, payload.Type())
	properties := payload.Properties()
	require.Len(t, properties, 3)

	expectedProperties := map[string]asyncapi.Schema{
		"id": &Schema{
			DescriptionField: "Id of the streetlight.",
			MinimumField:     refFloat64(0),
			TypeField:        "integer",
		},
		"lumens": &Schema{
			DescriptionField: "Light intensity measured in lumens.",
			MinimumField:     refFloat64(0),
			TypeField:        "integer",
		},
		"sentAt": &Schema{
			DescriptionField: "Date and time when the message was sent.",
			FormatField:      "date-time",
			TypeField:        "string",
		},
	}

	assert.Equal(t, expectedProperties, properties)
}

func refFloat64(v float64) *float64 {
	return &v
}

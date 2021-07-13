package v2

import (
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
}

//nolint:misspell
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
}

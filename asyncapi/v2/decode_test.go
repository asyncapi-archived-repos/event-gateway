package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen,misspell
func TestDecode(t *testing.T) {
	tests := []struct {
		name          string
		doc           []byte
		dst           interface{}
		expectedError error
		check         func(t *testing.T, dst interface{})
	}{
		{
			name: "from yaml file example-kafka.yaml",
			doc:  []byte("testdata/example-kafka.yaml"),
			dst:  &Document{},
			check: func(t *testing.T, dst interface{}) {
				require.IsType(t, (*Document)(nil), dst)
				doc := dst.(*Document)
				require.Len(t, doc.Servers(), 1)
				s := doc.Servers()[0]

				assert.True(t, s.HasName())
				assert.Equal(t, "test", s.Name())
				assert.True(t, s.HasDescription())
				assert.Equal(t, "Test broker", s.Description())
				assert.True(t, s.HasProtocol())
				assert.Equal(t, "kafka-secure", s.Protocol())
				assert.True(t, s.HasURL())
				assert.Equal(t, "localhost:9092", s.URL())
				assert.True(t, s.HasExtension("x-eventgateway-listener"))
				assert.Equal(t, "proxy:28002", s.Extension("x-eventgateway-listener"))
				assert.True(t, s.HasExtension("x-eventgateway-dial-mapping"))
				assert.Equal(t, "broker:9092", s.Extension("x-eventgateway-dial-mapping"))
				assert.Empty(t, s.Variables())
			},
		},
		{
			name: "from plain text yaml",
			doc: []byte(`
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
              description: Date and time when the message was sent.`),
			dst: &Document{},
			check: func(t *testing.T, dst interface{}) {
				require.IsType(t, (*Document)(nil), dst)
				doc := dst.(*Document)
				require.Len(t, doc.Servers(), 1)
				s := doc.Servers()[0]

				assert.True(t, s.HasName())
				assert.Equal(t, "mosquitto", s.Name())
				assert.False(t, s.HasDescription())
				assert.True(t, s.HasProtocol())
				assert.Equal(t, "mqtt", s.Protocol())
				assert.True(t, s.HasURL())
				assert.Equal(t, "mqtt://test.mosquitto.org", s.URL())
				assert.False(t, s.HasExtension("x-eventgateway-listener"))
				assert.False(t, s.HasExtension("x-eventgateway-dial-mapping"))
				assert.Empty(t, s.Variables())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Decode(test.doc, test.dst)
			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.check != nil {
				test.check(t, test.dst)
			}
		})
	}
}

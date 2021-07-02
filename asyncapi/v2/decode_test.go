package v2

import (
	"testing"

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
				assert.IsType(t, (*Document)(nil), dst)
				assert.Len(t, dst.(*Document).Servers(), 1)
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
				assert.IsType(t, (*Document)(nil), dst)
				doc := dst.(*Document)

				assert.Len(t, doc.Servers(), 1)
				assert.Equal(t, "mosquitto", doc.Servers()[0].Name())
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

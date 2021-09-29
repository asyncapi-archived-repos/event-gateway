package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	payload := []byte(`Hello!`)
	channel := "random-channel"
	msg := New(payload, channel)

	assert.Equal(t, payload, []byte(msg.Payload))
	assert.Equal(t, channel, msg.Metadata.Get(MetadataChannel))
}

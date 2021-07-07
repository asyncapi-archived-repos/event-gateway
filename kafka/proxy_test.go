package kafka

import (
	"bytes"
	"testing"

	"github.com/asyncapi/event-gateway/proxy"
	kafkaproxy "github.com/grepplabs/kafka-proxy/proxy"
	kafkaprotocol "github.com/grepplabs/kafka-proxy/proxy/protocol"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewKafka(t *testing.T) {
	tests := []struct {
		name        string
		c           ProxyConfig
		expectedErr error
	}{
		{
			name: "Proxy is created when config is valid",
			c: ProxyConfig{
				BrokersMapping: []string{"localhost:9092, localhost:28002"},
			},
		},
		{
			name:        "Proxy creation errors when config is invalid",
			c:           ProxyConfig{},
			expectedErr: errors.New("Brokers mapping is required"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewProxy(test.c)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, (proxy.Proxy)(nil), p)
			}
		})
	}
}

func TestRequestKeyHandler_Handle(t *testing.T) {
	tests := []struct {
		name              string
		request           []byte
		shouldReply       bool
		apiKey            int16
		shouldSkipRequest bool
	}{
		{
			name:        "Valid message",
			request:     []byte{0, 0, 0, 7, 0, 16, 99, 111, 110, 115, 111, 108, 101, 45, 112, 114, 111, 100, 117, 99, 101, 114, 255, 255, 0, 1}, // payload: 'valid message'
			shouldReply: true,
		},
		{
			name:        "Invalid message",
			request:     []byte{0, 0, 0, 8, 0, 16, 99, 111, 110, 115, 111, 108, 101, 45, 112, 114, 111, 100, 117, 99, 101, 114, 255, 255, 0, 1}, // payload: 'invalid message'
			shouldReply: true,
		},
		{
			name:              "Other Requests (different than Produce type) are skipped",
			request:           []byte{0, 0, 0, 1}, // fake payload that should not be read.
			apiKey:            int16(42),
			shouldReply:       true,
			shouldSkipRequest: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kv := &kafkaprotocol.RequestKeyVersion{
				ApiKey: test.apiKey,                  // default is 0, which is a Produce Request
				Length: int32(len(test.request)) + 4, // 4 bytes are ApiKey + Version located in all request headers (already read by the time of validating the msg).
			}

			readBytes := bytes.NewBuffer(nil)
			var h requestKeyHandler
			shouldReply, err := h.Handle(kv, bytes.NewReader(test.request), &kafkaproxy.RequestsLoopContext{}, readBytes)
			assert.NoError(t, err)
			assert.Equal(t, test.shouldReply, shouldReply)

			if test.shouldSkipRequest {
				assert.Empty(t, readBytes.Bytes()) // Payload is never read.
			} else {
				assert.Equal(t, readBytes.Len(), len(test.request))
			}
		})
	}
}

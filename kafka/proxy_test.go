package kafka

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/asyncapi/event-gateway/proxy"
	kafkaproxy "github.com/grepplabs/kafka-proxy/proxy"
	kafkaprotocol "github.com/grepplabs/kafka-proxy/proxy/protocol"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

// Note: Taking V8 as random version.
// correlation_id: 0, 0, 0, 6
// client_id_size: 0, 16
// client_id: 99, 111, 110, 115, 111, 108, 101, 45, 112, 114, 111, 100, 117, 99, 101, 114
// transactional_id_size: 255, 255
// acks: 0, 1
var produceRequestV8Headers = []byte{0, 0, 0, 6, 0, 16, 99, 111, 110, 115, 111, 108, 101, 45, 112, 114, 111, 100, 117, 99, 101, 114, 255, 255, 0, 1}

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
			request:     generateProduceRequestV8("valid message"),
			shouldReply: true,
		},
		{
			name:        "Invalid message",
			request:     generateProduceRequestV8("invalid message"),
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
			log := logrustest.NewGlobal()
			test.request = append(produceRequestV8Headers, test.request...) // This appends the Produce Request headers which we don't care for this test.
			kv := &kafkaprotocol.RequestKeyVersion{
				ApiVersion: 8,                            // All test data was grabbed from a Produce Request version 8.
				ApiKey:     test.apiKey,                  // default is 0, which is a Produce Request
				Length:     int32(len(test.request) + 4), // 4 bytes are ApiKey + Version located in all request headers (already read by the time of validating the msg).
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

			for _, l := range log.AllEntries() {
				assert.NotEqualf(t, l.Level, logrus.ErrorLevel, "%q logged error unexpected", l.Message) // We don't have a notification mechanism for errors yet
			}
		})
	}
}

func generateProduceRequestV8(payload string) []byte {
	raw := []byte{0, 0, 5, 220, 0, 0, 0, 1, 0, 4, 100, 101, 109, 111, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 81, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 69, 255, 255, 255, 255, 2, 42, 190, 231, 201, 0, 0, 0, 0, 0, 0, 0, 0, 1, 122, 129, 58, 51, 194, 0, 0, 1, 122, 129, 58, 51, 194, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 1, 38, 0, 0, 0, 1}

	// Based on https://kafka.apache.org/documentation/#record
	messageLen := make([]byte, 1)
	binary.PutVarint(messageLen, int64(len(payload)))

	bytesPayload := append(messageLen, []byte(payload)...)
	raw = append(raw, bytesPayload...)

	return append(raw, 0) // No headers
}

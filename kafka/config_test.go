package kafka

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestProxyConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      ProxyConfig
		expectedErr error
	}{
		{
			name: "Valid config. Only broker mapping",
			config: ProxyConfig{
				BrokersMapping: []string{"broker.mybrokers.org:9092,:9092"},
			},
		},
		{
			name: "Valid config. Broker mapping + Dial mapping",
			config: ProxyConfig{
				BrokersMapping:     []string{"broker.mybrokers.org:9092,:9092"},
				DialAddressMapping: []string{"broker.mybrokers.org:9092,192.168.1.10:9092"},
			},
		},
		{
			name:        "Invalid config. No broker mapping",
			expectedErr: errors.New("BrokersMapping is mandatory"),
		},
		{
			name: "Invalid config. Both broker and proxy can't listen to the same port within same host",
			config: ProxyConfig{
				BrokersMapping: []string{"localhost:9092,:9092"},
			},
			expectedErr: errors.New("broker and proxy can't listen to the same port on the same host. Broker is already listening at localhost:9092. Please configure a different listener port"),
		},
		{
			name: "Invalid config. Both broker and proxy are the same",
			config: ProxyConfig{
				BrokersMapping: []string{"broker.mybrokers.org:9092,broker.mybrokers.org:9092"},
			},
			expectedErr: errors.New("broker and proxy can't listen to the same port on the same host. Broker is already listening at broker.mybrokers.org:9092. Please configure a different listener port"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

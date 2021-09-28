package config

import (
	"testing"

	"github.com/asyncapi/event-gateway/kafka"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestKafkaProxy_ProxyConfig(t *testing.T) {
	tests := []struct {
		name                string
		config              *KafkaProxy
		doc                 []byte
		expectedProxyConfig func(*testing.T, *kafka.ProxyConfig) *kafka.ProxyConfig
		expectedErr         error
	}{
		{
			name: "Valid config. Only one broker",
			config: &KafkaProxy{
				BrokerFromServer: "test",
				// Also testing extra flags.
				ExtraFlags: pipeSeparatedValues{Values: []string{"arg1=arg1value", "arg2=arg2value"}},
			},
			expectedProxyConfig: func(_ *testing.T, _ *kafka.ProxyConfig) *kafka.ProxyConfig {
				return &kafka.ProxyConfig{
					BrokersMapping: []string{"broker.mybrokers.org:9092,:9092"},
					ExtraConfig:    []string{"arg1=arg1value", "arg2=arg2value"},
				}
			},
			doc: []byte(`testdata/simple-kafka.yaml`),
		},
		{
			name: "Valid config. Only one broker + enable message validation",
			config: &KafkaProxy{
				BrokerFromServer: "test",
				MessageValidation: MessageValidation{
					Enabled: true,
				},
			},
			expectedProxyConfig: func(t *testing.T, c *kafka.ProxyConfig) *kafka.ProxyConfig {
				assert.Equal(t, []string{"broker.mybrokers.org:9092,:9092"}, c.BrokersMapping)
				assert.Len(t, c.MessageHandlers, 1)
				return nil
			},
			doc: []byte(`testdata/simple-kafka.yaml`),
		},
		{
			name: "Valid config. Only one broker + Override listener port",
			config: &KafkaProxy{
				BrokerFromServer: "test",
			},
			expectedProxyConfig: func(t *testing.T, c *kafka.ProxyConfig) *kafka.ProxyConfig {
				return &kafka.ProxyConfig{
					BrokersMapping: []string{"broker.mybrokers.org:9092,:28002"},
				}
			},
			doc: []byte(`testdata/override-port-kafka.yaml`),
		},
		{
			name:   "Valid config. All brokers + Override listener port",
			config: &KafkaProxy{},
			expectedProxyConfig: func(t *testing.T, c *kafka.ProxyConfig) *kafka.ProxyConfig {
				return &kafka.ProxyConfig{
					BrokersMapping: []string{"broker.mybrokers.org:9092,:28002"},
				}
			},
			doc: []byte(`testdata/override-port-kafka.yaml`),
		},
		{
			name:   "Valid config. All brokers + multiple dial mapping",
			config: &KafkaProxy{},
			expectedProxyConfig: func(t *testing.T, c *kafka.ProxyConfig) *kafka.ProxyConfig {
				return &kafka.ProxyConfig{
					BrokersMapping:     []string{"broker.mybrokers.org:9092,:9092"},
					DialAddressMapping: []string{"0.0.0.0:28002,kafkaproxy.myapp.org:28002", "0.0.0.0:28003,kafkaproxy.myapp.org:28003"},
				}
			},
			doc: []byte(`testdata/dial-mapping-kafka.yaml`),
		},
		{
			name:        "Invalid config. Both broker and proxy are the same",
			config:      &KafkaProxy{},
			expectedErr: errors.New("broker and proxy can't listen to the same port on the same host. Broker is already listening at localhost:9092. Please configure a different listener port"),
			doc:         []byte(`testdata/invalid-same-address.yaml`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxyConfig, err := test.config.ProxyConfig(test.doc, false)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedProxyConfig != nil {
				if expectedConf := test.expectedProxyConfig(t, proxyConfig); expectedConf != nil {
					assert.EqualValues(t, expectedConf, proxyConfig)
				}
			}
		})
	}
}

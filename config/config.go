package config

import (
	"errors"
	"strings"

	"github.com/asyncapi/event-gateway/kafka"
)

// App holds the config for the whole application.
type App struct {
	Debug       bool
	AsyncAPIDoc []byte     `split_words:"true"`
	KafkaProxy  KafkaProxy `split_words:"true"`
}

// Validate validates the config.
func (c App) Validate() error {
	if len(c.AsyncAPIDoc) == 0 && len(c.KafkaProxy.BrokersMapping.Values) == 0 {
		return errors.New("either AsyncAPIDoc or KafkaProxyBrokersMapping config should be provided")
	}

	return nil
}

// ProxyConfig creates a config struct for the Kafka Proxy.
func (c App) ProxyConfig() (kafka.ProxyConfig, error) {
	return c.KafkaProxy.ProxyConfig(c.AsyncAPIDoc, c.Debug)
}

type pipeSeparatedValues struct {
	Values []string
}

func (b *pipeSeparatedValues) Set(value string) error {
	b.Values = strings.Split(value, "|")
	return nil
}

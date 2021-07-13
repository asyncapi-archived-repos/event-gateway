package config

import (
	"strings"

	"github.com/asyncapi/event-gateway/kafka"
)

// App holds the config for the whole application.
type App struct {
	Debug       bool
	AsyncAPIDoc []byte     `split_words:"true"`
	KafkaProxy  KafkaProxy `split_words:"true"`
}

// ProxyConfig creates a config struct for the Kafka Proxy.
func (c App) ProxyConfig() (*kafka.ProxyConfig, error) {
	return c.KafkaProxy.ProxyConfig(c.AsyncAPIDoc, c.Debug)
}

type pipeSeparatedValues struct {
	Values []string
}

func (b *pipeSeparatedValues) Set(value string) error {
	b.Values = strings.Split(value, "|")
	return nil
}

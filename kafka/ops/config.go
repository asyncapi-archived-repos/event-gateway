package ops

import (
	"crypto/tls"

	"github.com/Shopify/sarama"
)

type Config struct {
	Connection struct {
		Addrs []string
		TLS   struct {
			Enabled bool
			Config  *tls.Config
		}
	}
}

func (c *Config) Sarama() *sarama.Config {
	return sarama.NewConfig() // TODO
}

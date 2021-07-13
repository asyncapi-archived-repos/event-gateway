package config

import (
	"fmt"
	"strings"

	v2 "github.com/asyncapi/event-gateway/asyncapi/v2"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/pkg/errors"
)

// KafkaProxy holds the config for later configuring a Kafka proxy.
type KafkaProxy struct {
	ServerName         string              `split_words:"true"`
	BrokersMapping     pipeSeparatedValues `split_words:"true"`
	BrokersDialMapping pipeSeparatedValues `split_words:"true"`
	ExtraFlags         pipeSeparatedValues `split_words:"true"`
}

// ProxyConfig creates a config struct for the Kafka Proxy based on a given AsyncAPI doc (if provided).
func (c *KafkaProxy) ProxyConfig(doc []byte, debug bool) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig
	if len(doc) > 0 {
		conf, err := configFromDoc(doc)
		if err != nil {
			return kafkaProxyConfig, err
		}

		kafkaProxyConfig = conf
	} else {
		kafkaProxyConfig = kafka.ProxyConfig{
			BrokersMapping:     c.BrokersMapping.Values,
			DialAddressMapping: c.BrokersDialMapping.Values,
		}
	}

	kafkaProxyConfig.Debug = debug
	kafkaProxyConfig.ExtraConfig = c.ExtraFlags.Values

	return kafkaProxyConfig, nil
}

func configFromDoc(d []byte) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig

	doc := new(v2.Document)
	if err := v2.Decode(d, doc); err != nil {
		return kafkaProxyConfig, errors.Wrap(err, "error decoding AsyncAPI json doc to Document struct")
	}

	for _, s := range doc.Servers() {
		if strings.HasPrefix(s.Protocol(), "kafka") {
			listenAt := s.Extension("x-eventgateway-listener")
			if listenAt == nil {
				return kafkaProxyConfig, errors.New("x-eventgateway-listener extension is mandatory for opening ports to listen connections")
			}

			// TODO should be only port config also accepted?
			kafkaProxyConfig.BrokersMapping = append(kafkaProxyConfig.BrokersMapping, fmt.Sprintf("%s,%s", s.URL(), listenAt))

			if dialMapping := s.Extension("x-eventgateway-dial-mapping"); dialMapping != nil {
				kafkaProxyConfig.DialAddressMapping = append(kafkaProxyConfig.DialAddressMapping, fmt.Sprintf("%s,%s", s.URL(), dialMapping))
			}
		}
	}

	return kafkaProxyConfig, nil
}

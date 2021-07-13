package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/asyncapi/event-gateway/asyncapi"
	v2 "github.com/asyncapi/event-gateway/asyncapi/v2"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/pkg/errors"
)

// KafkaProxy holds the config for later configuring a Kafka proxy.
type KafkaProxy struct {
	BrokerFromServer   string              `split_words:"true"`
	BrokersMapping     pipeSeparatedValues `split_words:"true"`
	BrokersDialMapping pipeSeparatedValues `split_words:"true"`
	ExtraFlags         pipeSeparatedValues `split_words:"true"`
}

// ProxyConfig creates a config struct for the Kafka Proxy based on a given AsyncAPI doc (if provided).
func (c *KafkaProxy) ProxyConfig(doc []byte, debug bool) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig
	if len(doc) == 0 && len(c.BrokersMapping.Values) == 0 {
		return kafkaProxyConfig, errors.New("either AsyncAPIDoc or KafkaProxyBrokersMapping config should be provided")
	}

	if c.BrokerFromServer != "" && len(doc) == 0 {
		return kafkaProxyConfig, errors.New("AsyncAPIDoc should be provided when setting BrokerFromServer")
	}

	if len(doc) > 0 {
		conf, err := c.configFromDoc(doc)
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

	if err := kafkaProxyConfig.Validate(); err != nil {
		return kafkaProxyConfig, err
	}

	return kafkaProxyConfig, nil
}

func (c *KafkaProxy) configFromDoc(d []byte) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig

	doc := new(v2.Document)
	if err := v2.Decode(d, doc); err != nil {
		return kafkaProxyConfig, errors.Wrap(err, "error decoding AsyncAPI json doc to Document struct")
	}

	var err error
	if c.BrokerFromServer != "" {
		kafkaProxyConfig, err = kafkaProxyConfigFromServer(c.BrokerFromServer, doc)
	} else {
		kafkaProxyConfig, err = kafkaProxyConfigFromAllServers(doc)
	}

	return kafkaProxyConfig, err
}

func isValidKafkaProtocol(s asyncapi.Server) bool {
	return strings.HasPrefix(s.Protocol(), "kafka")
}

func kafkaProxyConfigFromAllServers(doc asyncapi.Document) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig
	for _, s := range doc.Servers() {
		if isValidKafkaProtocol(s) {
			l := s.Extension(asyncapi.ExtensionEventGatewayListener)
			listenAt, ok := l.(string)
			if listenAt == "" || !ok {
				return kafkaProxyConfig, fmt.Errorf("please specify either %s extension, env vars or an AsyncAPI doc in orderr to set the Kafka proxy listener(s)", asyncapi.ExtensionEventGatewayListener)
			}

			kafkaProxyConfig.BrokersMapping = append(kafkaProxyConfig.BrokersMapping, fmt.Sprintf("%s,%s", s.URL(), listenAt))

			if dialMapping := s.Extension(asyncapi.ExtensionEventGatewayDialMapping); dialMapping != nil {
				kafkaProxyConfig.DialAddressMapping = append(kafkaProxyConfig.DialAddressMapping, fmt.Sprintf("%s,%s", s.URL(), dialMapping))
			}
		}
	}

	return kafkaProxyConfig, nil
}

func kafkaProxyConfigFromServer(name string, doc asyncapi.Document) (kafka.ProxyConfig, error) {
	var kafkaProxyConfig kafka.ProxyConfig
	s, ok := doc.Server(name)
	if !ok {
		return kafkaProxyConfig, fmt.Errorf("server %s not found in the provided AsyncAPI doc", name)
	}

	if !isValidKafkaProtocol(s) {
		return kafkaProxyConfig, fmt.Errorf("server %s has no kafka protocol configured but '%s'", name, s.Protocol())
	}

	// Only one broker will be configured
	_, port, err := net.SplitHostPort(s.URL())
	if err != nil {
		return kafkaProxyConfig, errors.Wrapf(err, "error getting port from broker %s. URL:%s", s.Name(), s.URL())
	}

	kafkaProxyConfig.BrokersMapping = append(kafkaProxyConfig.BrokersMapping, fmt.Sprintf("%s,:%s", s.URL(), port))

	return kafkaProxyConfig, nil
}

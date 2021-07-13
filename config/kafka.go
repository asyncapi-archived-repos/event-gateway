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
func (c *KafkaProxy) ProxyConfig(doc []byte, debug bool) (*kafka.ProxyConfig, error) {
	if len(doc) == 0 && len(c.BrokersMapping.Values) == 0 {
		return nil, errors.New("either AsyncAPIDoc or KafkaProxyBrokersMapping config should be provided")
	}

	if c.BrokerFromServer != "" && len(doc) == 0 {
		return nil, errors.New("AsyncAPIDoc should be provided when setting BrokerFromServer")
	}

	var kafkaProxyConfig *kafka.ProxyConfig
	var err error
	if len(doc) > 0 {
		kafkaProxyConfig, err = c.configFromDoc(doc)
	} else {
		kafkaProxyConfig, err = kafka.NewProxyConfig(c.BrokersMapping.Values, kafka.WithDialAddressMapping(c.BrokersDialMapping.Values), kafka.WithExtra(c.ExtraFlags.Values))
	}

	if err != nil {
		return nil, err
	}

	kafkaProxyConfig.Debug = debug

	return kafkaProxyConfig, nil
}

func (c *KafkaProxy) configFromDoc(d []byte) (*kafka.ProxyConfig, error) {
	doc := new(v2.Document)
	if err := v2.Decode(d, doc); err != nil {
		return nil, errors.Wrap(err, "error decoding AsyncAPI json doc to Document struct")
	}

	if c.BrokerFromServer != "" {
		return kafkaProxyConfigFromServer(c.BrokerFromServer, doc)
	}

	return kafkaProxyConfigFromAllServers(doc.Servers())
}

func isValidKafkaProtocol(s asyncapi.Server) bool {
	return strings.HasPrefix(s.Protocol(), "kafka")
}

func kafkaProxyConfigFromAllServers(servers []asyncapi.Server) (*kafka.ProxyConfig, error) {
	var brokersMapping []string
	var dialAddressMapping []string
	for _, s := range servers {
		if isValidKafkaProtocol(s) {
			l := s.Extension(asyncapi.ExtensionEventGatewayListener)
			listenAt, ok := l.(string)
			if listenAt == "" || !ok {
				return nil, fmt.Errorf("please specify either %s extension, env vars or an AsyncAPI doc in orderr to set the Kafka proxy listener(s)", asyncapi.ExtensionEventGatewayListener)
			}

			brokersMapping = append(brokersMapping, fmt.Sprintf("%s,%s", s.URL(), listenAt))

			if dialMapping := s.Extension(asyncapi.ExtensionEventGatewayDialMapping); dialMapping != nil {
				dialAddressMapping = append(dialAddressMapping, fmt.Sprintf("%s,%s", s.URL(), dialMapping))
			}
		}
	}

	return kafka.NewProxyConfig(brokersMapping, kafka.WithDialAddressMapping(dialAddressMapping))
}

func kafkaProxyConfigFromServer(name string, doc asyncapi.Document) (*kafka.ProxyConfig, error) {
	s, ok := doc.Server(name)
	if !ok {
		return nil, fmt.Errorf("server %s not found in the provided AsyncAPI doc", name)
	}

	if !isValidKafkaProtocol(s) {
		return nil, fmt.Errorf("server %s has no kafka protocol configured but '%s'", name, s.Protocol())
	}

	// Only one broker will be configured
	_, port, err := net.SplitHostPort(s.URL())
	if err != nil {
		return nil, errors.Wrapf(err, "error getting port from broker %s. URL:%s", s.Name(), s.URL())
	}

	var opts []kafka.Option
	if dialMapping := s.Extension(asyncapi.ExtensionEventGatewayDialMapping); dialMapping != nil {
		opts = append(opts, kafka.WithDialAddressMapping([]string{fmt.Sprintf("%s,%s", s.URL(), dialMapping)}))
	}

	return kafka.NewProxyConfig([]string{fmt.Sprintf("%s,:%s", s.URL(), port)}, opts...)
}

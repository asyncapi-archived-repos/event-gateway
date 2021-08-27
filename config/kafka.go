package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/asyncapi/event-gateway/proxy"

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
	MessageValidation  MessageValidation   `split_words:"true"`
}

// MessageValidation holds the config about message validation.
type MessageValidation struct {
	Enabled  bool
	Notifier proxy.ValidationErrorNotifier
}

// NotifyValidationErrorOnChan sets a channel as ValidationError notifier.
func NotifyValidationErrorOnChan(errChan chan *proxy.ValidationError) Opt {
	return func(app *App) {
		app.KafkaProxy.MessageValidation.Notifier = proxy.ValidationErrorToChanNotifier(errChan)
	}
}

// NewKafkaProxy creates a KafkaProxy with defaults.
func NewKafkaProxy() *KafkaProxy {
	return &KafkaProxy{MessageValidation: MessageValidation{
		Enabled: true,
	}}
}

// ProxyConfig creates a config struct for the Kafka Proxy based on a given AsyncAPI doc (if provided).
func (c *KafkaProxy) ProxyConfig(doc []byte, debug bool, messageHandlers ...kafka.MessageHandler) (*kafka.ProxyConfig, error) {
	if len(doc) == 0 && len(c.BrokersMapping.Values) == 0 {
		return nil, errors.New("either AsyncAPIDoc or KafkaProxyBrokersMapping config should be provided")
	}

	if c.BrokerFromServer != "" && len(doc) == 0 {
		return nil, errors.New("AsyncAPIDoc should be provided when setting BrokerFromServer")
	}

	var kafkaProxyConfig *kafka.ProxyConfig
	var err error
	if len(doc) > 0 {
		kafkaProxyConfig, err = c.configFromDoc(doc, kafka.WithExtra(c.ExtraFlags.Values))
	} else {
		kafkaProxyConfig, err = kafka.NewProxyConfig(c.BrokersMapping.Values, kafka.WithDialAddressMapping(c.BrokersDialMapping.Values), kafka.WithExtra(c.ExtraFlags.Values))
	}

	if err != nil {
		return nil, err
	}

	kafkaProxyConfig.Debug = debug
	kafkaProxyConfig.MessageHandlers = append(kafkaProxyConfig.MessageHandlers, messageHandlers...)

	return kafkaProxyConfig, nil
}

func (c *KafkaProxy) configFromDoc(d []byte, opts ...kafka.Option) (*kafka.ProxyConfig, error) {
	doc := new(v2.Document)
	if err := v2.Decode(d, doc); err != nil {
		return nil, errors.Wrap(err, "error decoding AsyncAPI json doc to Document struct")
	}

	if c.MessageValidation.Enabled {
		validator, err := v2.FromDocJSONSchemaMessageValidator(doc)
		if err != nil {
			return nil, errors.Wrap(err, "error creating message validator")
		}

		if notifier := c.MessageValidation.Notifier; notifier != nil {
			validator = proxy.NotifyOnValidationError(validator, notifier)
		}

		opts = append(opts, kafka.WithMessageHandlers(validateMessageHandler(validator)))
	}

	if c.BrokerFromServer != "" {
		return kafkaProxyConfigFromServer(c.BrokerFromServer, doc, opts...)
	}

	return kafkaProxyConfigFromAllServers(doc.Servers(), opts...)
}

func isValidKafkaProtocol(s asyncapi.Server) bool {
	return strings.HasPrefix(s.Protocol(), "kafka")
}

func kafkaProxyConfigFromAllServers(servers []asyncapi.Server, opts ...kafka.Option) (*kafka.ProxyConfig, error) {
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

	opts = append(opts, kafka.WithDialAddressMapping(dialAddressMapping))

	return kafka.NewProxyConfig(brokersMapping, opts...)
}

func kafkaProxyConfigFromServer(name string, doc asyncapi.Document, opts ...kafka.Option) (*kafka.ProxyConfig, error) {
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

	if dialMapping := s.Extension(asyncapi.ExtensionEventGatewayDialMapping); dialMapping != nil {
		opts = append(opts, kafka.WithDialAddressMapping([]string{fmt.Sprintf("%s,%s", s.URL(), dialMapping)}))
	}

	return kafka.NewProxyConfig([]string{fmt.Sprintf("%s,:%s", s.URL(), port)}, opts...)
}

func validateMessageHandler(validator proxy.MessageValidator) kafka.MessageHandler {
	return func(msg kafka.Message) error {
		pMsg := &proxy.Message{
			Context: proxy.MessageContext{
				Channel: msg.Context.Topic,
			},
			Key:   msg.Key,
			Value: msg.Value,
		}

		if len(msg.Headers) > 0 {
			pMsg.Headers = make([]proxy.MessageHeader, len(msg.Headers))
			for i := 0; i < len(msg.Headers); i++ {
				pMsg.Headers[i] = proxy.MessageHeader{
					Key:   msg.Headers[i].Key,
					Value: msg.Headers[i].Value,
				}
			}
		}

		_, err := validator(pMsg)

		return err
	}
}

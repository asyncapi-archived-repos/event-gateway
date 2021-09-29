package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/asyncapi/event-gateway/asyncapi"
	v2 "github.com/asyncapi/event-gateway/asyncapi/v2"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/asyncapi/event-gateway/message"
	"github.com/asyncapi/event-gateway/message/handler"
	"github.com/pkg/errors"
)

// KafkaProxy holds the config for later configuring a Kafka proxy.
type KafkaProxy struct {
	BrokerFromServer  string              `split_words:"true" desc:"When configuring from an AsyncAPI doc, this allows the user to only configure one server instead of all"`
	MessageValidation MessageValidation   `split_words:"true" desc:""`
	ExtraFlags        pipeSeparatedValues `split_words:"true" desc:"Advanced configuration. Configure any flag from https://github.com/grepplabs/kafka-proxy/blob/4f3b89fbaecb3eb82426f5dcff5f76188ea9a9dc/cmd/kafka-proxy/server.go#L85-L195. Multiple values can be configured by using pipe separation (|)"`
}

// MessageValidation holds the config about message validation.
type MessageValidation struct {
	Enabled  bool                            `default:"true" desc:"Enable or disable validation of Kafka messages"`
	Notifier message.ValidationErrorNotifier `envconfig:"-"`
}

// NotifyValidationErrorOnChan sets a channel as ValidationError notifier.
func NotifyValidationErrorOnChan(errChan chan *message.ValidationError) Opt {
	return func(app *App) {
		app.KafkaProxy.MessageValidation.Notifier = message.ValidationErrorToChanNotifier(errChan)
	}
}

// NewKafkaProxy creates a KafkaProxy with defaults.
func NewKafkaProxy() *KafkaProxy {
	return &KafkaProxy{MessageValidation: MessageValidation{
		Enabled: true,
	}}
}

// ProxyConfig creates a config struct for the Kafka Proxy based on a given AsyncAPI doc (if provided).
func (c *KafkaProxy) ProxyConfig(d []byte, debug bool) (*kafka.ProxyConfig, error) {
	if len(d) == 0 {
		return nil, errors.New("AsyncAPIDoc config should be provided")
	}

	doc := new(v2.Document)
	if err := v2.Decode(d, doc); err != nil {
		return nil, errors.Wrap(err, "error decoding AsyncAPI json doc to Document struct")
	}

	opts := []kafka.Option{kafka.WithExtra(c.ExtraFlags.Values), kafka.WithDebug(debug)}
	if c.MessageValidation.Enabled {
		validator, err := v2.FromDocJSONSchemaMessageValidator(doc)
		if err != nil {
			return nil, errors.Wrap(err, "error creating message validator")
		}

		opts = append(opts, kafka.WithMessageHandler(handler.ValidateMessage(validator, c.MessageValidation.Notifier, false)))
	}

	servers := doc.Servers()
	if c.BrokerFromServer != "" {
		// Pick up only the specified server
		s, ok := doc.Server(c.BrokerFromServer)
		if !ok {
			return nil, fmt.Errorf("server %s not found in the provided AsyncAPI doc", s.Name())
		}

		if !isValidKafkaProtocol(s) {
			return nil, fmt.Errorf("server %s has no kafka protocol configured but '%s'", s.Name(), s.Protocol())
		}

		servers = []asyncapi.Server{s}
	}

	return kafkaProxyConfigFromServers(servers, opts...)
}

func isValidKafkaProtocol(s asyncapi.Server) bool {
	return strings.HasPrefix(s.Protocol(), "kafka")
}

func kafkaProxyConfigFromServers(servers []asyncapi.Server, opts ...kafka.Option) (*kafka.ProxyConfig, error) {
	brokersMapping, dialAddressMapping, err := extractAddressMappingFromServers(servers...)
	if err != nil {
		return nil, err
	}

	if len(brokersMapping) == 0 {
		return nil, errors.New("No Kafka brokers were found when configuring")
	}

	if len(dialAddressMapping) > 0 {
		opts = append(opts, kafka.WithDialAddressMapping(dialAddressMapping))
	}

	return kafka.NewProxyConfig(brokersMapping, opts...)
}

func extractAddressMappingFromServers(servers ...asyncapi.Server) (brokersMapping []string, dialAddressMapping []string, err error) {
	for _, s := range servers {
		if !isValidKafkaProtocol(s) {
			continue
		}

		var listenAt string
		// If extension is configured, it overrides the value of the port.
		if overridePort := s.Extension(asyncapi.ExtensionEventGatewayListener); overridePort != nil {
			if val := fmt.Sprintf("%v", overridePort); val != "" { // Convert value to string rep as can be either string or number
				if host, _, _ := net.SplitHostPort(val); host == "" {
					val = ":" + val // If no host, prefix with : as localhost is inferred
				}
				listenAt = val
			}
		} else {
			// Use the same port as remote but locally.
			_, val, err := net.SplitHostPort(s.URL())
			if err != nil {
				return nil, nil, errors.Wrapf(err, "error getting port from broker %s. URL:%s", s.Name(), s.URL())
			}

			listenAt = ":" + val // Prefix with : as localhost is inferred
		}

		brokersMapping = append(brokersMapping, fmt.Sprintf("%s,%s", s.URL(), listenAt))
		if dialMapping := s.Extension(asyncapi.ExtensionEventGatewayDialMapping); dialMapping != nil {
			dialAddressMapping = append(dialAddressMapping, strings.Split(dialMapping.(string), "|")...)
		}
	}

	return brokersMapping, dialAddressMapping, nil
}

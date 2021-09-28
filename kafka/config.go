package kafka

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/asyncapi/event-gateway/message"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var localHostIpv4 = regexp.MustCompile(`127\.0\.0\.\d+`)

// ProxyConfig holds the configuration for the Kafka Proxy.
type ProxyConfig struct {
	BrokersMapping     []string
	DialAddressMapping []string
	ExtraConfig        []string
	MessageHandlers    []message.Handler
	Debug              bool
}

// Option represents a functional configuration for the Proxy.
type Option func(*ProxyConfig) error

// WithMessageHandlers configures sets a given list of message.Middlware as proxy message handlers.
func WithMessageHandlers(handlers ...message.Handler) Option {
	return func(c *ProxyConfig) error {
		c.MessageHandlers = append(c.MessageHandlers, handlers...)
		return nil
	}
}

// WithDebug enables debug.
func WithDebug() Option {
	return func(c *ProxyConfig) error {
		c.Debug = true
		return nil
	}
}

// WithDialAddressMapping configures Dial Address Mapping.
func WithDialAddressMapping(mapping []string) Option {
	return func(c *ProxyConfig) error {
		c.DialAddressMapping = mapping
		return nil
	}
}

// WithExtra configures extra parameters.
func WithExtra(extra []string) Option {
	return func(c *ProxyConfig) error {
		c.ExtraConfig = extra
		return nil
	}
}

// NewProxyConfig creates a new ProxyConfig.
func NewProxyConfig(brokersMapping []string, opts ...Option) (*ProxyConfig, error) {
	c := &ProxyConfig{BrokersMapping: brokersMapping}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, c.Validate()
}

// Validate validates ProxyConfig.
func (c *ProxyConfig) Validate() error {
	if len(c.BrokersMapping) == 0 {
		return errors.New("BrokersMapping is mandatory")
	}

	invalidFormatMsg := "BrokersMapping should be in form 'remotehost:remoteport,localhost:localport"
	for _, m := range c.BrokersMapping {
		v := strings.Split(m, ",")
		if len(v) != 2 {
			return errors.New(invalidFormatMsg)
		}

		remoteHost, remotePort, err := net.SplitHostPort(v[0])
		if err != nil {
			return errors.Wrap(err, invalidFormatMsg)
		}

		localHost, localPort, err := net.SplitHostPort(v[1])
		if err != nil {
			return errors.Wrap(err, invalidFormatMsg)
		}

		if remoteHost == localHost && remotePort == localPort || (isLocalHost(remoteHost) && isLocalHost(localHost) && remotePort == localPort) {
			return fmt.Errorf("broker and proxy can't listen to the same port on the same host. Broker is already listening at %s. Please configure a different listener port", v[0])
		}
	}

	if len(c.MessageHandlers) == 0 {
		logrus.Warn("There are no message handlers configured")
	}

	return nil
}

func isLocalHost(host string) bool {
	return host == "" ||
		host == "::1" ||
		host == "0:0:0:0:0:0:0:1" ||
		localHostIpv4.MatchString(host) ||
		host == "localhost"
}

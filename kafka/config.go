package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"strings"

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var localHostIpv4 = regexp.MustCompile(`127\.0\.0\.\d+`)

// ProxyConfig holds the configuration for the Kafka Proxy.
type ProxyConfig struct {
	BrokersMapping     []string
	DialAddressMapping []string
	ExtraConfig        []string
	MessageHandler     watermillmessage.HandlerFunc
	MessagePublisher   watermillmessage.Publisher
	PublishToTopic     string
	MessageSubscriber  watermillmessage.Subscriber
	TLS                *TLSConfig
	Debug              bool
}

// TLSConfig holds configuration for TLS.
type TLSConfig struct {
	Enable             bool
	InsecureSkipVerify bool   `split_words:"true"`
	ClientCertFile     string `split_words:"true"`
	ClientKeyFile      string `split_words:"true"`
	CAChainCertFile    string `split_words:"true"`
}

// Config returns a *tls.Config based on current config.
func (c *TLSConfig) Config() (*tls.Config, error) {
	cfg := &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify} //nolint:gosec

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCertFile, c.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	if c.CAChainCertFile != "" {
		caCertPEMBlock, err := ioutil.ReadFile(c.CAChainCertFile)
		if err != nil {
			return nil, err
		}
		rootCAs := x509.NewCertPool()
		if ok := rootCAs.AppendCertsFromPEM(caCertPEMBlock); !ok {
			return nil, errors.New("Failed to parse client root certificate")
		}

		cfg.RootCAs = rootCAs
	}

	return cfg, nil
}

// ProxyOption represents a functional configuration for the Proxy.
type ProxyOption func(*ProxyConfig) error

// WithMessageHandler ...
func WithMessageHandler(handler watermillmessage.HandlerFunc) ProxyOption {
	return func(c *ProxyConfig) error {
		c.MessageHandler = handler
		return nil
	}
}

// WithMessagePublisher configures a publisher where the messages will be published after being handled.
func WithMessagePublisher(publisher watermillmessage.Publisher, topic string) ProxyOption {
	return func(c *ProxyConfig) error {
		c.MessagePublisher = publisher
		c.PublishToTopic = topic
		return nil
	}
}

// WithMessageSubscriber configures a subscriber subscribed to the messages published by the configured c.MessagePublisher.
func WithMessageSubscriber(subscriber watermillmessage.Subscriber) ProxyOption {
	return func(c *ProxyConfig) error {
		c.MessageSubscriber = subscriber
		return nil
	}
}

// WithDebug enables/disables debug.
func WithDebug(enabled bool) ProxyOption {
	return func(c *ProxyConfig) error {
		c.Debug = enabled
		return nil
	}
}

// WithDialAddressMapping configures Dial Address Mapping.
func WithDialAddressMapping(mapping []string) ProxyOption {
	return func(c *ProxyConfig) error {
		c.DialAddressMapping = mapping
		return nil
	}
}

// WithExtra configures extra parameters.
func WithExtra(extra []string) ProxyOption {
	return func(c *ProxyConfig) error {
		c.ExtraConfig = extra
		return nil
	}
}

// NewProxyConfig creates a new ProxyConfig.
func NewProxyConfig(brokersMapping []string, opts ...ProxyOption) (*ProxyConfig, error) {
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

	if c.MessageHandler == nil {
		logrus.Warn("There is no message handler configured")
		return nil
	} else if c.MessagePublisher != nil && c.PublishToTopic == "" || c.MessagePublisher == nil && c.PublishToTopic != "" {
		return fmt.Errorf("MessagePublisher and PublishToTopic should be set together")
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

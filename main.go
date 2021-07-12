package main

import (
	"context"
	"fmt"
	"strings"

	v2 "github.com/asyncapi/event-gateway/asyncapi/v2"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type config struct {
	Debug                        bool
	AsyncAPIDoc                  []byte              `split_words:"true"`
	KafkaProxyBrokersMapping     pipeSeparatedValues ` split_words:"true"`
	KafkaProxyBrokersDialMapping pipeSeparatedValues `split_words:"true"`
	KafkaProxyExtraFlags         pipeSeparatedValues `split_words:"true"`
}

type pipeSeparatedValues struct {
	values []string
}

func (b *pipeSeparatedValues) Set(value string) error { //nolint:unparam
	b.values = strings.Split(value, "|")
	return nil
}

func main() {
	var c config
	if err := envconfig.Process("eventgateway", &c); err != nil {
		logrus.Fatal(err)
	}

	if len(c.AsyncAPIDoc) == 0 && len(c.KafkaProxyBrokersMapping.values) == 0 {
		logrus.Fatalln("Either AsyncAPIDoc or KafkaProxyBrokersMapping config should be provided")
	}

	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var kafkaProxyConfig kafka.ProxyConfig
	if len(c.AsyncAPIDoc) > 0 {
		conf, err := configFromDoc(c.AsyncAPIDoc)
		if err != nil {
			logrus.WithError(err).Fatal()
		}

		kafkaProxyConfig = conf
	} else {
		kafkaProxyConfig = kafka.ProxyConfig{
			BrokersMapping:     c.KafkaProxyBrokersMapping.values,
			DialAddressMapping: c.KafkaProxyBrokersDialMapping.values,
		}
	}

	kafkaProxyConfig.Debug = c.Debug
	kafkaProxyConfig.ExtraConfig = c.KafkaProxyExtraFlags.values

	kafkaProxy, err := kafka.NewProxy(kafkaProxyConfig)
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := kafkaProxy(context.Background()); err != nil {
		logrus.Fatalln(err)
	}
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

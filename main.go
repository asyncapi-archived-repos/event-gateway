package main

import (
	"context"
	"strings"

	"github.com/asyncapi/event-gateway/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	Debug                        bool
	KafkaProxyBrokersMapping     pipeSeparatedValues `required:"true" split_words:"true"`
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

	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	kafkaProxyConfig := kafka.ProxyConfig{
		BrokersMapping:     c.KafkaProxyBrokersMapping.values,
		DialAddressMapping: c.KafkaProxyBrokersDialMapping.values,
		ExtraConfig:        c.KafkaProxyExtraFlags.values,
		Debug:              c.Debug,
	}

	kafkaProxy, err := kafka.NewProxy(proxyConfig)
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := kafkaProxy(context.Background()); err != nil {
		logrus.Fatalln(err)
	}
}

package main

import (
	"context"

	"github.com/asyncapi/event-gateway/config"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	var c config.App
	if err := envconfig.Process("eventgateway", &c); err != nil {
		logrus.WithError(err).Fatal()
	}

	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	kafkaProxyConfig, err := c.ProxyConfig()
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	kafkaProxy, err := kafka.NewProxy(kafkaProxyConfig)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	if err := kafkaProxy(context.Background()); err != nil {
		logrus.WithError(err).Fatal()
	}
}

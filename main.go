package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asyncapi/event-gateway/proxy"

	"github.com/asyncapi/event-gateway/config"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	validationErrChan := make(chan *proxy.ValidationError)
	c := config.NewApp(config.NotifyValidationErrorOnChan(validationErrChan))

	if err := envconfig.Process("eventgateway", c); err != nil {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleInterruptions(cancel)

	// At this moment, we do nothing else.
	go logValidationErrors(ctx, validationErrChan)

	if err := kafkaProxy(context.Background()); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func logValidationErrors(ctx context.Context, validationErrChan chan *proxy.ValidationError) {
	for {
		select {
		case validationErr, ok := <-validationErrChan:
			if !ok {
				return
			}

			logrus.WithField("validation_errors", validationErr.String()).Errorf("error validating message")
		case <-ctx.Done():
			return
		}
	}
}

func handleInterruptions(cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-c
		logrus.WithField("signal", s).Info("Stopping AsyncAPI Event-Gateway...")
		cancel()
		time.Sleep(time.Second)
		os.Exit(0)
	}()
}

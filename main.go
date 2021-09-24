package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asyncapi/event-gateway/config"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/asyncapi/event-gateway/message"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
)

const configPrefix = "eventgateway"

func main() {
	validationErrChan := make(chan *message.ValidationError)
	c := config.NewApp(config.NotifyValidationErrorOnChan(validationErrChan))

	if err := envconfig.Process(configPrefix, c); err != nil {
		_ = envconfig.Usage(configPrefix, c)
		logrus.WithError(err).Fatal()
	}

	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	kafkaProxyConfig, err := c.ProxyConfig()
	if err != nil {
		_ = envconfig.Usage(configPrefix, c)
		logrus.WithError(err).Fatal()
	}

	kafkaProxy, err := kafka.NewProxy(kafkaProxyConfig)
	if err != nil {
		_ = envconfig.Usage(configPrefix, c)
		logrus.WithError(err).Fatal()
	}

	m := melody.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleInterruptions(cancel, func() error {
		return m.CloseWithMsg(melody.FormatCloseMessage(1000, "The server says goodbye :)"))
	})

	r := chi.NewRouter()
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ = m.HandleRequest(w, r)
	})

	go func() {
		address := fmt.Sprintf(":%v", c.WSServerPort)
		logrus.Infof("Websocket server listening on %s", address)
		if err := http.ListenAndServe(address, r); err != nil {
			logrus.WithError(err).Fatal("error running websocket server")
		}
	}()

	go handleValidationErrors(ctx, validationErrChan, m)
	if err := kafkaProxy(context.Background()); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func handleValidationErrors(ctx context.Context, validationErrChan chan *message.ValidationError, m *melody.Melody) {
	for {
		select {
		case validationErr, ok := <-validationErrChan:
			if !ok {
				return
			}

			content, err := json.Marshal(validationErr)
			if err != nil {
				logrus.WithError(err).Error("error encoding validationError")
				content = []byte(`For some reason, this validation error can't be seen. Please drop us an issue on github.'`)
			}

			logrus.WithField("validation_error", string(content)).Debug("a message validation error has been found")

			if err := m.Broadcast(content); err != nil {
				logrus.WithError(err).Error("error broadcasting message to all ws sessions")
			}
		case <-ctx.Done():
			return
		}
	}
}

func handleInterruptions(cancel context.CancelFunc, funcs ...func() error) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-c
		logrus.WithField("signal", s).Info("Stopping AsyncAPI Event-Gateway...")
		cancel()
		for _, f := range funcs {
			if err := f(); err != nil {
				logrus.WithError(err).Error("error shutting down")
			}
		}
		time.Sleep(time.Second)
		os.Exit(0)
	}()
}

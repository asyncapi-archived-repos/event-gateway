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

	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/asyncapi/event-gateway/config"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/asyncapi/event-gateway/message"
	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const configPrefix = "eventgateway"

func main() {
	c := config.NewApp()
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

	watermillLogger := message.NewWatermillLogrusLogger(logrus.StandardLogger())
	messageRouter, err := watermillmessage.NewRouter(watermillmessage.RouterConfig{}, watermillLogger)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	defer messageRouter.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := melody.New()
	handleInterruptions(cancel, func() error {
		return m.CloseWithMsg(melody.FormatCloseMessage(1000, "The server says goodbye :)"))
	})

	if kafkaProxyConfig.MessageSubscriber != nil {
		defer kafkaProxyConfig.MessageSubscriber.Close()
		defer kafkaProxyConfig.MessagePublisher.Close()
		messageRouter.AddNoPublisherHandler("consume-validation-errors-from-kafka", c.KafkaProxy.MessageValidation.PublishToKafkaTopic, kafkaProxyConfig.MessageSubscriber, validationErrorsHandler(m))
	}

	kafkaProxy, err := kafka.NewProxy(kafkaProxyConfig, messageRouter)
	if err != nil {
		_ = envconfig.Usage(configPrefix, c)
		logrus.WithError(err).Fatal()
	}

	runWebsocketServer(c.WSServerPort, "/ws", m)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return messageRouter.Run(ctx) // Note: Can not be called until fully configured.
	})
	group.Go(func() error {
		return kafkaProxy(ctx)
	})

	if err := group.Wait(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func runWebsocketServer(port int, path string, m *melody.Melody) {
	r := chi.NewRouter()
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		_ = m.HandleRequest(w, r)
	})

	go func() {
		address := fmt.Sprintf(":%v", port)
		logrus.Infof("Websocket server listening on %s", address)
		if err := http.ListenAndServe(address, r); err != nil {
			logrus.WithError(err).Fatal("error running websocket server")
		}
	}()
}

func validationErrorsHandler(m *melody.Melody) watermillmessage.NoPublishHandlerFunc {
	return func(msg *watermillmessage.Message) error {
		validationError, err := message.ValidationErrorFromMessage(msg)
		if err != nil {
			logrus.WithError(err).WithField("uuid", msg.UUID).Error("Couldn't determine if the message is rather invalid or not")
			return nil // no retry
		}

		if validationError == nil {
			return nil
		}

		content, err := json.Marshal(msg)
		if err != nil {
			logrus.WithError(err).Error("error marshaling message")
			content = []byte(fmt.Sprintf(
				"The message %s is invalid: %s. However, there was an error during encoding and couldn't be shown completely. Please drop us an issue on github.",
				msg.UUID,
				validationError.Error(),
			))
		}

		logrus.WithError(validationError).Debug("Message is invalid")

		if err := m.Broadcast(content); err != nil {
			logrus.WithError(err).Error("error broadcasting message to all ws sessions")
		}

		return nil
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

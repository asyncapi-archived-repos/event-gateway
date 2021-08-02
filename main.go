package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asyncapi/event-gateway/config"
	"github.com/asyncapi/event-gateway/kafka"
	"github.com/asyncapi/event-gateway/proxy"
	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/olahol/melody"
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

	m := melody.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleInterruptions(cancel, func() error {
		return m.CloseWithMsg(melody.FormatCloseMessage(1000, "The server says goodbye :)"))
	})

	r := chi.NewRouter()
	r.Get("/", http.FileServer(http.Dir("./web/static")).ServeHTTP)
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

func handleValidationErrors(ctx context.Context, validationErrChan chan *proxy.ValidationError, m *melody.Melody) {
	// Using logrus for handling ws output text.
	noopLogger := logrus.New()
	noopLogger.Out = io.Discard

	// Using logrus formatter for ws output text.
	f := logrus.TextFormatter{
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	}

	for {
		select {
		case validationErr, ok := <-validationErrChan:
			if !ok {
				return
			}

			logrus.WithField("validation_errors", validationErr.String()).Debug("error validating message")
			entry := noopLogger.WithFields(logrus.Fields{
				"channel":           validationErr.Msg.Context.Channel,
				"key":               string(validationErr.Msg.Key),
				"headers":           validationErr.Msg.Headers,
				"validation_errors": validationErr.Result.Errors(),
			}).WithTime(time.Now())
			entry.Level = logrus.InfoLevel

			txt, _ := f.Format(entry)
			if err := m.Broadcast(txt); err != nil {
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

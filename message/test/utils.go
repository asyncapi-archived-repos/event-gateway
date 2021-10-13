package test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/asyncapi/event-gateway/message"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------- //
// Functions and helpers useful for testing  //
// ----------------------------------------- //

// NoopPublisher is a publisher that publishes nothing.
// Implements watermillmessage.Publisher interface.
type NoopPublisher struct{}

func (n NoopPublisher) Publish(_ string, _ ...*watermillmessage.Message) error {
	return nil
}

func (n NoopPublisher) Close() error {
	return nil
}

// ReliablePublisher creates a publisher with a subscriber that will make the test fail if the number of published messages does not reach `messagesToAssert` before a given timeout.
// Should be only used on tests.
func ReliablePublisher(t *testing.T, topic string, messagesToAssert int, timeout time.Duration) (publisher *gochannel.GoChannel, done chan struct{}) {
	p := gochannel.NewGoChannel(gochannel.Config{Persistent: true}, watermill.NopLogger{})
	msgs, err := p.Subscribe(context.Background(), topic)
	if err != nil {
		t.Fatal(err)
	}

	done = make(chan struct{})
	var messagesReceived int
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		timer := time.NewTimer(timeout)
		for {
			select {
			case <-msgs:
				messagesReceived++
				if messagesReceived >= messagesToAssert {
					return
				}
			case <-timer.C:
				t.Errorf("Publisher was expected to publish at least %v, but only %v were published within %s", messagesToAssert, messagesReceived, timeout)
				return
			}
		}
	}()

	return p, done
}

// AssertCalledHandlerFunc returns a new HandlerFunc based on a given one that will make the test fail if the given handler has not been called at least `times` before a given timeout.
func AssertCalledHandlerFunc(t *testing.T, h watermillmessage.HandlerFunc, times uint32, timeout time.Duration) (handler watermillmessage.HandlerFunc, done chan struct{}) {
	var calledTimes uint32
	handler = func(msg *watermillmessage.Message) ([]*watermillmessage.Message, error) {
		atomic.AddUint32(&calledTimes, +1)
		return h(msg)
	}

	done = make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		timer := time.NewTimer(timeout)
		for {
			select {
			case <-timer.C:
				t.Errorf("Handler was expected to be called at least %v, but only %v times were called within %s", times, atomic.LoadUint32(&calledTimes), timeout)
				return
			default:
				if atomic.LoadUint32(&calledTimes) >= times {
					return
				}
			}
		}
	}()

	return handler, done
}

// NewRouterWithLogs creates a watermillmessage.Router with some needed test cleanup actions. It also configures a logger.
func NewRouterWithLogs(t *testing.T, logger *logrus.Logger) *watermillmessage.Router {
	return newRouter(t, message.NewWatermillLogrusLogger(logger))
}

// NewRouter creates a watermillmessage.Router with some needed test cleanup actions.
func NewRouter(t *testing.T) *watermillmessage.Router {
	return newRouter(t, watermill.NopLogger{})
}

func newRouter(t *testing.T, logger watermill.LoggerAdapter) *watermillmessage.Router {
	r, err := watermillmessage.NewRouter(watermillmessage.RouterConfig{}, logger)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		assert.NoError(t, r.Close())
	})

	return r
}

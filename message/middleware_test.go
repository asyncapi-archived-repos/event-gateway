package message

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_EnsureOrderOfExecution(t *testing.T) {
	var timesCalled uint32
	msg := &Message{}
	returnedMsg, err := Chain(
		orderAwareMiddleware(t, 0, &timesCalled),
		orderAwareMiddleware(t, 1, &timesCalled),
		orderAwareMiddleware(t, 2, &timesCalled),
	)(msg)
	assert.NoError(t, err)
	assert.Same(t, msg, returnedMsg)
}

func TestMiddleware_OneErrorShouldStopExecution(t *testing.T) {
	m1 := func(m *Message) (*Message, error) {
		return m, nil
	}

	m2 := func(m *Message) (*Message, error) {
		return nil, errors.New("whatever error")
	}

	m3 := func(m *Message) (*Message, error) {
		assert.FailNow(t, "This should not be reached ever")
		return m, nil
	}

	msg := &Message{}
	returnedMsg, err := Chain(m1, m2, m3)(msg)
	assert.EqualError(t, err, "whatever error")
	assert.Nil(t, returnedMsg)
}

func TestMiddleware_MissingReturnShouldError(t *testing.T) {
	missingReturnMiddleware := func(m *Message) (*Message, error) {
		return nil, nil
	}

	msg := &Message{}
	returnedMsg, err := Chain(missingReturnMiddleware)(msg)
	assert.EqualError(t, err, ErrMiddlewareMissingReturn.Error())
	assert.Nil(t, returnedMsg)
}

func orderAwareMiddleware(t *testing.T, expectedCounter uint32, counter *uint32) Middleware {
	return func(m *Message) (*Message, error) {
		assert.Equalf(t, expectedCounter, *counter, "This middleware was expected to be called in position %d, but was %d", expectedCounter, *counter)
		*counter++
		return m, nil
	}
}

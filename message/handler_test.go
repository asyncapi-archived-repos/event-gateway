package message

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandler_EnsureOrderOfExecution(t *testing.T) {
	var timesCalled uint32
	msg := &Message{}
	returnedMsg, err := HandlersChain(
		orderAwareHandler(t, 0, &timesCalled),
		orderAwareHandler(t, 1, &timesCalled),
		orderAwareHandler(t, 2, &timesCalled),
	)(msg)
	assert.NoError(t, err)
	assert.Same(t, msg, returnedMsg)
}

func TestHandler_OneErrorShouldStopExecution(t *testing.T) {
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
	returnedMsg, err := HandlersChain(m1, m2, m3)(msg)
	assert.EqualError(t, err, "whatever error")
	assert.Nil(t, returnedMsg)
}

func TestHandler_MissingReturnShouldError(t *testing.T) {
	missingReturnHandler := func(m *Message) (*Message, error) {
		return nil, nil
	}

	msg := &Message{}
	returnedMsg, err := HandlersChain(missingReturnHandler)(msg)
	assert.EqualError(t, err, ErrHandlerMissingReturn.Error())
	assert.Nil(t, returnedMsg)
}

func orderAwareHandler(t *testing.T, expectedCounter uint32, counter *uint32) Handler {
	return func(m *Message) (*Message, error) {
		assert.Equalf(t, expectedCounter, *counter, "This handler was expected to be called in position %d, but was %d", expectedCounter, *counter)
		*counter++
		return m, nil
	}
}

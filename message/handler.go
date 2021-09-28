package message

import "github.com/pkg/errors"

// ErrHandlerMissingReturn is the error used when a handler did not return either a message or an error but one is expected.
var ErrHandlerMissingReturn = errors.New("a handler did not return either a message or an error but one is expected")

// Handler is a function applied to a Message and returns the same Message after processing it or modifying it.
// A handler should return either a Message or an error.
type Handler func(m *Message) (*Message, error)

// HandlersChain calls all given handlers in a pipeline-fashioned process. Order of handlers is preserved during execution.
func HandlersChain(handlers ...Handler) Handler {
	return func(m *Message) (*Message, error) {
		var err error
		for _, handler := range handlers {
			m, err = handler(m)
			if err != nil {
				break
			}

			if m == nil {
				return nil, ErrHandlerMissingReturn
			}
		}

		return m, err
	}
}

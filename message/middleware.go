package message

import "github.com/pkg/errors"

// ErrMiddlewareMissingReturn is the error used when a middleware did not return either a message or an error but one is expected.
var ErrMiddlewareMissingReturn = errors.New("a middleware did not return either a message or an error but one is expected")

// Middleware is a function applied to a Message and returns the same Message after processing it or modifying it.
// A middleware should return either a Message or an error.
type Middleware func(m *Message) (*Message, error)

// Chain calls all given middlewares in a pipeline-fashioned process. Order of middlewares is preserved during execution.
func Chain(middlewares ...Middleware) Middleware {
	return func(m *Message) (*Message, error) {
		var err error
		for _, middleware := range middlewares {
			m, err = middleware(m)
			if err != nil {
				break
			}

			if m == nil {
				return nil, ErrMiddlewareMissingReturn
			}
		}

		return m, err
	}
}

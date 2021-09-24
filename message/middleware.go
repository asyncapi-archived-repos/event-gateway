package message

import "github.com/pkg/errors"

// ErrMiddlewareMissingReturn is the error used when a middleware did not return either a message or an error but one is expected.
var ErrMiddlewareMissingReturn = errors.New("a middleware did not return either a message or an error but one is expected")

type Middleware func(m *Message) (*Message, error) // TODO consider even not returning error here

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

package proxy

import "context"

// Proxy represents an event-gateway proxy.
type Proxy func(ctx context.Context) error

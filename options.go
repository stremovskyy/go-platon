package go_platon

import (
	"time"

	internalhttp "github.com/stremovskyy/go-platon/internal/http"
)

// ClientOption configures the underlying HTTP client.
type ClientOption func(*internalhttp.Options)

func WithTimeout(d time.Duration) ClientOption {
	return func(o *internalhttp.Options) {
		o.Timeout = d
	}
}

func WithKeepAlive(d time.Duration) ClientOption {
	return func(o *internalhttp.Options) {
		o.KeepAlive = d
	}
}

func WithMaxIdleConns(n int) ClientOption {
	return func(o *internalhttp.Options) {
		o.MaxIdleConns = n
	}
}

func WithIdleConnTimeout(d time.Duration) ClientOption {
	return func(o *internalhttp.Options) {
		o.IdleConnTimeout = d
	}
}

// NewClient creates a client with custom HTTP options.
func NewClient(opts ...ClientOption) Platon {
	o := internalhttp.DefaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}

	return &client{
		platonClient: internalhttp.NewClient(o),
	}
}


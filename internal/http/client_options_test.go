package http

import (
	"net/http"
	"testing"
)

func TestNewClient_NilOptions_UsesDefaults(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Fatalf("expected client, got nil")
	}
	if c.options == nil {
		t.Fatalf("expected default options, got nil")
	}

	defaults := DefaultOptions()
	if c.options.Timeout != defaults.Timeout {
		t.Fatalf("timeout mismatch: want %v, got %v", defaults.Timeout, c.options.Timeout)
	}
	if c.options.DialTimeout != defaults.DialTimeout {
		t.Fatalf("dial timeout mismatch: want %v, got %v", defaults.DialTimeout, c.options.DialTimeout)
	}
	if c.options.MaxIdleConnsPerHost != defaults.MaxIdleConnsPerHost {
		t.Fatalf(
			"max idle conns per host mismatch: want %d, got %d",
			defaults.MaxIdleConnsPerHost,
			c.options.MaxIdleConnsPerHost,
		)
	}
}

func TestNewClient_NormalizesInvalidOptions(t *testing.T) {
	c := NewClient(&Options{
		Timeout:               -1,
		KeepAlive:             0,
		DialTimeout:           -1,
		TLSHandshakeTimeout:   0,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 0,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   0,
		MaxConnsPerHost:       -1,
		IdleConnTimeout:       0,
	})

	defaults := DefaultOptions()
	if c.options.Timeout != defaults.Timeout {
		t.Fatalf("timeout mismatch: want %v, got %v", defaults.Timeout, c.options.Timeout)
	}
	if c.options.MaxConnsPerHost != defaults.MaxConnsPerHost {
		t.Fatalf("max conns per host mismatch: want %d, got %d", defaults.MaxConnsPerHost, c.options.MaxConnsPerHost)
	}
	if c.options.ResponseHeaderTimeout != defaults.ResponseHeaderTimeout {
		t.Fatalf(
			"response header timeout mismatch: want %v, got %v",
			defaults.ResponseHeaderTimeout,
			c.options.ResponseHeaderTimeout,
		)
	}
}

func TestNewClient_TransportIsHardenedByDefault(t *testing.T) {
	c := NewClient(nil)
	transport, ok := c.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("transport type mismatch: got %T", c.client.Transport)
	}
	if !transport.ForceAttemptHTTP2 {
		t.Fatalf("expected ForceAttemptHTTP2=true")
	}
	if transport.Proxy == nil {
		t.Fatalf("expected proxy function to be configured")
	}
	if c.client.CheckRedirect == nil {
		t.Fatalf("expected check redirect function to be configured")
	}
}

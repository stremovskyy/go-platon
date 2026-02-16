/*
 * MIT License
 *
 * Copyright (c) 2026 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
	c := NewClient(
		&Options{
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
		},
	)

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

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

package go_platon

import (
	"net/http"
	"time"

	internalhttp "github.com/stremovskyy/go-platon/internal/http"
	"github.com/stremovskyy/recorder"
)

type clientConfig struct {
	httpOptions *internalhttp.Options
	httpClient  *http.Client
	recorder    recorder.Recorder
}

func defaultClientConfig() *clientConfig {
	return &clientConfig{
		httpOptions: internalhttp.DefaultOptions(),
	}
}

// Option configures the platon client.
type Option func(*clientConfig)

// ClientOption is a backward-compatible alias for Option.
type ClientOption = Option

func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) {
		c.httpOptions.Timeout = d
	}
}

func WithKeepAlive(d time.Duration) Option {
	return func(c *clientConfig) {
		c.httpOptions.KeepAlive = d
	}
}

func WithMaxIdleConns(n int) Option {
	return func(c *clientConfig) {
		c.httpOptions.MaxIdleConns = n
	}
}

func WithIdleConnTimeout(d time.Duration) Option {
	return func(c *clientConfig) {
		c.httpOptions.IdleConnTimeout = d
	}
}

// WithClient overrides the default underlying net/http client.
func WithClient(cl *http.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = cl
		if cl != nil {
			// Keep request context timeout consistent with the explicitly provided client.
			c.httpOptions.Timeout = cl.Timeout
		}
	}
}

// WithRecorder attaches a recorder to the client.
func WithRecorder(r recorder.Recorder) Option {
	return func(c *clientConfig) {
		c.recorder = r
	}
}

// NewClient creates a platon client with custom options.
func NewClient(opts ...Option) Platon {
	cfg := defaultClientConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}

	httpClient := internalhttp.NewClient(cfg.httpOptions)
	if cfg.httpClient != nil {
		httpClient.SetClient(cfg.httpClient)
	}
	if cfg.recorder != nil {
		httpClient.SetRecorder(cfg.recorder)
	}

	return &client{
		platonClient: httpClient,
	}
}

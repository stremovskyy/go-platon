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

import "time"

// Options for http client
type Options struct {
	Timeout               time.Duration
	KeepAlive             time.Duration
	DialTimeout           time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	IsDebug               bool
}

func DefaultOptions() *Options {
	return &Options{
		Timeout:               15 * time.Second,
		KeepAlive:             30 * time.Second,
		DialTimeout:           10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		IsDebug:               false,
	}
}

func normalizeOptions(options *Options) *Options {
	defaults := DefaultOptions()
	if options == nil {
		return defaults
	}

	normalized := *options

	if normalized.Timeout <= 0 {
		normalized.Timeout = defaults.Timeout
	}
	if normalized.KeepAlive <= 0 {
		normalized.KeepAlive = defaults.KeepAlive
	}
	if normalized.DialTimeout <= 0 {
		normalized.DialTimeout = defaults.DialTimeout
	}
	if normalized.TLSHandshakeTimeout <= 0 {
		normalized.TLSHandshakeTimeout = defaults.TLSHandshakeTimeout
	}
	if normalized.ResponseHeaderTimeout <= 0 {
		normalized.ResponseHeaderTimeout = defaults.ResponseHeaderTimeout
	}
	if normalized.ExpectContinueTimeout <= 0 {
		normalized.ExpectContinueTimeout = defaults.ExpectContinueTimeout
	}
	if normalized.MaxIdleConns <= 0 {
		normalized.MaxIdleConns = defaults.MaxIdleConns
	}
	if normalized.MaxIdleConnsPerHost <= 0 {
		normalized.MaxIdleConnsPerHost = defaults.MaxIdleConnsPerHost
	}
	if normalized.MaxConnsPerHost <= 0 {
		normalized.MaxConnsPerHost = defaults.MaxConnsPerHost
	}
	if normalized.IdleConnTimeout <= 0 {
		normalized.IdleConnTimeout = defaults.IdleConnTimeout
	}

	return &normalized
}

type CtxKey string

const (
	CtxKeyRequestID CtxKey = "request_id"
)

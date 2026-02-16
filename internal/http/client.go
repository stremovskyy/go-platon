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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/log"
	"github.com/stremovskyy/go-platon/platon"
	"github.com/stremovskyy/recorder"
)

type Client struct {
	client   *http.Client
	options  *Options
	logger   *log.Logger
	recorder recorder.Recorder
}

const maxResponseBodyBytes = 4 << 20 // 4 MiB

// Api handles Platon API request.
func (c *Client) Api(apiRequest *platon.Request, apiURL string) (*platon.Response, error) {
	return c.sendURLEncodedRequest(apiURL, apiRequest, c.logger)
}

// WithRecorder attaches a recorder to the client.
func (c *Client) WithRecorder(rec recorder.Recorder) *Client {
	c.recorder = rec

	return c
}

// SetClient allows replacing the underlying net/http client.
func (c *Client) SetClient(cl *http.Client) {
	c.client = cl
}

// SetRecorder allows setting a recorder explicitly.
func (c *Client) SetRecorder(r recorder.Recorder) {
	c.recorder = r
}

func (c *Client) sendURLEncodedRequest(apiURL string, unsignedRequest *platon.Request, logger *log.Logger) (*platon.Response, error) {
	requestID := uuid.New().String()
	logger.Debug("API URL: %v", apiURL)
	logger.Debug("Request ID: %v", requestID)

	if unsignedRequest == nil {
		return nil, c.logAndReturnError("request is nil", platon.ErrRequestIsNil, logger, requestID, nil)
	}

	signedRequest, err := unsignedRequest.SignAndPrepare()
	if err != nil {
		return nil, c.logAndReturnError("cannot sign request", err, logger, requestID, nil)
	}

	encodedForm, err := encodeRequestMap(signedRequest.ToMap())
	if err != nil {
		return nil, c.logAndReturnError("cannot encode request", err, logger, requestID, nil)
	}
	logger.Debug("Request (%s):\n%s", FormURLEncodedContentType, PrettyPrintFormURLEncodedBody(encodedForm))

	ctx := context.Background()
	if c.options != nil && c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}
	ctx = context.WithValue(ctx, CtxKeyRequestID, requestID)

	tags := tagsRetriever(signedRequest)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(encodedForm))
	if err != nil {
		return nil, c.logAndReturnError("cannot create request", err, logger, requestID, tags)
	}
	c.setHeaders(req, requestID)

	if c.recorder != nil {
		if err := c.recorder.RecordRequest(ctx, nil, requestID, []byte(encodedForm), tags); err != nil {
			logger.Error("cannot record request: %v", err)
		}
	}

	if c.client == nil {
		return nil, c.logAndReturnError("http client is nil", fmt.Errorf("http client is nil"), logger, requestID, tags)
	}

	tStart := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send request", err, logger, requestID, tags)
	}
	if resp == nil {
		return nil, c.logAndReturnError(
			"invalid response",
			fmt.Errorf("http response is nil"),
			logger,
			requestID,
			tags,
		)
	}
	if resp.Body == nil {
		return nil, c.logAndReturnError(
			"invalid response",
			fmt.Errorf("http response body is nil"),
			logger,
			requestID,
			tags,
		)
	}
	logger.Debug("Request time: %v", time.Since(tStart))

	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes+1))
	if err != nil {
		return nil, c.logAndReturnError("cannot read response", err, logger, requestID, tags)
	}

	logger.Debug("Response: %v", FormatBodyForDebug(resp.Header.Get("Content-Type"), raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	if len(raw) == 0 {
		return nil, c.logAndReturnError("no response bytes", fmt.Errorf("empty response"), logger, requestID, tags)
	}
	if len(raw) > maxResponseBodyBytes {
		return nil, c.logAndReturnError(
			"response too large",
			fmt.Errorf("response exceeds %d bytes", maxResponseBodyBytes),
			logger,
			requestID,
			tags,
		)
	}

	if c.recorder != nil {
		if err := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); err != nil {
			logger.Error("cannot record response: %v", err)
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, c.logAndReturnError(
			"unexpected response status",
			fmt.Errorf("status=%d body=%s", resp.StatusCode, truncateBodyForError(raw)),
			logger,
			requestID,
			tags,
		)
	}

	response, err := platon.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, c.logAndReturnError("cannot unmarshal response", err, logger, requestID, tags)
	}

	return response, response.GetError()
}

func encodeRequestMap(requestMap map[string]interface{}) (string, error) {
	formValues := url.Values{}

	for key, value := range requestMap {
		if value == nil {
			continue
		}

		switch typed := value.(type) {
		case string:
			formValues.Set(key, typed)
		case []byte:
			formValues.Set(key, string(typed))
		default:
			rawValue, err := json.Marshal(value)
			if err != nil {
				return "", fmt.Errorf("cannot marshal field %q: %w", key, err)
			}
			formValues.Set(key, string(rawValue))
		}
	}

	return formValues.Encode(), nil
}

// logAndReturnError logs an error and optionally records it.
func (c *Client) logAndReturnError(msg string, err error, logger *log.Logger, requestID string, tags map[string]string) error {
	logger.Error("%s: %v", msg, err)

	if c.recorder != nil {
		ctx := context.WithValue(context.Background(), CtxKeyRequestID, requestID)
		if err := c.recorder.RecordError(ctx, nil, requestID, err, tags); err != nil {
			logger.Error("cannot record error: %v", err)
		}
	}

	return err
}

// setHeaders sets common headers for all requests.
func (c *Client) setHeaders(req *http.Request, requestID string) {
	req.Header.Set("Content-Type", FormURLEncodedContentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GO PLATON/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)
}

// safeClose ensures the body is closed properly and logs any error.
func (c *Client) safeClose(body io.ReadCloser, logger *log.Logger) {
	if err := body.Close(); err != nil {
		logger.Error("cannot close response body: %v", err)
	}
}

func tagsRetriever(request *platon.Request) map[string]string {
	tags := make(map[string]string)
	if request == nil {
		return tags
	}

	if request.Action != "" {
		tags["action"] = request.Action
	}
	if request.OrderID != nil {
		tags["order_id"] = *request.OrderID
	}
	if request.TransId != nil {
		tags["trans_id"] = *request.TransId
	}

	return tags
}

func truncateBodyForError(raw []byte) string {
	const max = 512
	if len(raw) <= max {
		return string(raw)
	}
	return string(raw[:max]) + "...(truncated)"
}

// NewClient initializes a new HTTP client with options.
func NewClient(options *Options) *Client {
	options = normalizeOptions(options)

	dialer := &net.Dialer{
		Timeout:   options.DialTimeout,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          options.MaxIdleConns,
		MaxIdleConnsPerHost:   options.MaxIdleConnsPerHost,
		MaxConnsPerHost:       options.MaxConnsPerHost,
		IdleConnTimeout:       options.IdleConnTimeout,
		TLSHandshakeTimeout:   options.TLSHandshakeTimeout,
		ResponseHeaderTimeout: options.ResponseHeaderTimeout,
		ExpectContinueTimeout: options.ExpectContinueTimeout,
		DisableCompression:    true,
	}

	cl := &http.Client{
		Transport: tr,
		Timeout:   options.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client{
		client:  cl,
		options: options,
		logger:  log.NewLogger("Platon HTTP: "),
	}
}

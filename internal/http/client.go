package http

import (
	"bytes"
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
	client         *http.Client
	options        *Options
	logger         *log.Logger
	applePayLogger *log.Logger
	recorder       recorder.Recorder
}

// Api handles the standard Platon API request.
func (c *Client) Api(apiRequest *platon.Request, url string) (*platon.Response, error) {
	return c.sendURLEncodedRequest(url, apiRequest, c.logger)
}

// WithRecorder attaches a recorder to the client and returns the client for method chaining.
func (c *Client) WithRecorder(rec recorder.Recorder) *Client {
	c.recorder = rec
	return c
}

// sendRequest handles sending an HTTP request and processing the response.
func (c *Client) sendRequest(apiURL string, unsignedRequest *platon.Request, logger *log.Logger) (*platon.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
	defer cancel()
	requestID := uuid.New().String()
	logger.Debug("API URL: %v", apiURL)
	logger.Debug("Request ID: %v", requestID)

	signedRequest, err := unsignedRequest.SignAndPrepare()
	if err != nil {
		return nil, c.logAndReturnError("cannot sign request", err, logger, requestID, nil)
	}

	jsonBody, err := json.Marshal(signedRequest)
	if err != nil {
		return nil, c.logAndReturnError("cannot marshal request", err, logger, requestID, nil)
	}

	logger.Debug("Request: %v", string(jsonBody))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, c.logAndReturnError("cannot create request", err, logger, requestID, nil)
	}

	c.setCommonHeaders(req, requestID)
	req.Header.Set("Content-Type", "application/json")

	// Generate and log cURL command
	curlCmd := generateCurlCommand(req, jsonBody)
	logger.Debug("cURL equivalent: %s", curlCmd)

	tags := make(map[string]string)
	if signedRequest != nil {
		tags["action"] = string(signedRequest.Action)
	}
	if c.recorder != nil {
		if err := c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags); err != nil {
			logger.Error("cannot record request: %v", err)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send request", err, logger, requestID, tags)
	}
	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read response", err, logger, requestID, tags)
	}

	logger.Debug("Response status: %v", resp.StatusCode)

	if len(raw) == 0 {
		return nil, c.logAndReturnError("No response bytes", fmt.Errorf("empty response"), logger, requestID, tags)
	}

	logger.Debug("Response: %v", string(raw))

	if c.recorder != nil {
		if err := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); err != nil {
			logger.Error("cannot record response: %v", err)
		}
	}

	response, err := platon.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, c.logAndReturnError("cannot unmarshal response", err, logger, requestID, tags)
	}

	return response, response.GetError()
}

func (c *Client) sendURLEncodedRequest(apiURL string, unsignedRequest *platon.Request, logger *log.Logger) (*platon.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
	defer cancel()
	requestID := uuid.New().String()
	logger.Debug("API URL: %v", apiURL)
	logger.Debug("Request ID: %v", requestID)

	signedRequest, err := unsignedRequest.SignAndPrepare()
	if err != nil {
		return nil, c.logAndReturnError("cannot sign request", err, logger, requestID, nil)
	}

	// Convert the request to a map for URL encoding
	requestMap := signedRequest.ToMap()
	if err != nil {
		return nil, c.logAndReturnError("cannot convert request to map", err, logger, requestID, nil)
	}

	// Create form data with URL encoding
	formValues := url.Values{}
	for key, value := range requestMap {
		// Convert values to string as needed
		strValue := ""
		switch v := value.(type) {
		case string:
			strValue = v
		case []byte:
			strValue = string(v)
		default:
			// For complex types, use JSON string representation
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, c.logAndReturnError("cannot marshal field value", err, logger, requestID, nil)
			}
			strValue = string(jsonBytes)
		}
		formValues.Set(key, strValue)
	}

	encodedForm := formValues.Encode()
	logger.Debug("URL-encoded request: %v", encodedForm)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(encodedForm))
	if err != nil {
		return nil, c.logAndReturnError("cannot create request", err, logger, requestID, nil)
	}

	c.setCommonHeaders(req, requestID)
	// Set headers for URL-encoded form (must not be overwritten by common headers).
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tags := make(map[string]string)
	if signedRequest != nil {
		tags["action"] = string(signedRequest.Action)
	}
	if c.recorder != nil {
		if err := c.recorder.RecordRequest(ctx, nil, requestID, []byte(encodedForm), tags); err != nil {
			logger.Error("cannot record request: %v", err)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.logAndReturnError("cannot send request", err, logger, requestID, tags)
	}
	defer c.safeClose(resp.Body, logger)

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, c.logAndReturnError("cannot read response", err, logger, requestID, tags)
	}

	logger.Debug("Response status: %v", resp.StatusCode)

	if len(raw) == 0 {
		return nil, c.logAndReturnError("No response bytes", fmt.Errorf("empty response"), logger, requestID, tags)
	}

	logger.Debug("Response: %v", string(raw))

	if c.recorder != nil {
		if err := c.recorder.RecordResponse(ctx, nil, requestID, raw, tags); err != nil {
			logger.Error("cannot record response: %v", err)
		}
	}

	response, err := platon.UnmarshalJSONResponse(raw)
	if err != nil {
		return nil, c.logAndReturnError("cannot unmarshal response", err, logger, requestID, tags)
	}

	return response, response.GetError()
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

// setCommonHeaders sets headers common to all requests.
// Content-Type must be set by the concrete request sender (JSON vs x-www-form-urlencoded).
func (c *Client) setCommonHeaders(req *http.Request, requestID string) {
	headers := map[string]string{
		"Accept":       "application/json",
		"User-Agent":   "GO PLATON/" + consts.Version,
		"X-Request-ID": requestID,
		"Api-Version":  consts.ApiVersion,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// safeClose ensures the body is closed properly and logs any error.
func (c *Client) safeClose(body io.ReadCloser, logger *log.Logger) {
	if err := body.Close(); err != nil {
		logger.Error("cannot close response body: %v", err)
	}
}

// NewClient initializes a new HTTP client with options.
func NewClient(options *Options) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext:        dialer.DialContext,
	}

	return &Client{
		client:         &http.Client{Transport: tr, Timeout: options.Timeout},
		options:        options,
		logger:         log.NewLogger("Platon HTTP: "),
		applePayLogger: log.NewLogger("Platon ApplePay: "),
	}
}

func getApiURL(formID platon.ActionCode) string {
	if formID == platon.ActionCodeSALE {
		return consts.ApiVerifyURL
	}
	return "ERROR_URL"
}

func generateCurlCommand(req *http.Request, body []byte) string {
	curl := fmt.Sprintf("curl -X %s '%s'", req.Method, req.URL.String())

	for key, values := range req.Header {
		for _, value := range values {
			curl += fmt.Sprintf(" -H '%s: %s'", key, value)
		}
	}

	if len(body) > 0 {
		curl += fmt.Sprintf(" -d '%s'", string(body))
	}

	return curl
}

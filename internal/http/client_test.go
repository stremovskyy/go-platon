package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestApi_UsesFormURLEncodedContentType(t *testing.T) {
	var gotContentType string
	var gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ACCEPTED"}`))
	}))
	defer srv.Close()

	auth := &platon.Auth{Key: "k", Secret: "secret123"}
	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := platon.NewRequest(platon.ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(platon.HashTypeCardTokenPayment)

	c := NewClient(DefaultOptions())
	resp, err := c.Api(req, srv.URL)
	if err != nil {
		t.Fatalf("Api() error: %v", err)
	}
	if resp == nil || resp.Result == nil || *resp.Result != platon.ResultAccepted {
		t.Fatalf("unexpected response: %+v", resp)
	}

	if gotContentType != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type mismatch: want application/x-www-form-urlencoded, got %q", gotContentType)
	}

	if !strings.Contains(gotBody, "client_key=clientKey") {
		t.Fatalf("expected body to contain client_key, got %q", gotBody)
	}
	if !strings.Contains(gotBody, "card_token=") {
		t.Fatalf("expected body to contain card_token, got %q", gotBody)
	}
}

func TestApi_ReturnsErrorOnNon2xxStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`gateway down`))
	}))
	defer srv.Close()

	auth := &platon.Auth{Key: "k", Secret: "secret123"}
	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := platon.NewRequest(platon.ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(platon.HashTypeCardTokenPayment)

	c := NewClient(DefaultOptions())
	_, err := c.Api(req, srv.URL)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "status=502") {
		t.Fatalf("expected status code in error, got %q", err.Error())
	}
}

func TestApi_ReturnsErrorWhenResponseIsTooLarge(t *testing.T) {
	tooLarge := bytes.Repeat([]byte("x"), maxResponseBodyBytes+16)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(tooLarge)
	}))
	defer srv.Close()

	auth := &platon.Auth{Key: "k", Secret: "secret123"}
	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := platon.NewRequest(platon.ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(platon.HashTypeCardTokenPayment)

	c := NewClient(DefaultOptions())
	_, err := c.Api(req, srv.URL)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "response exceeds") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestApi_ReturnsErrorOnNilResponseBody(t *testing.T) {
	auth := &platon.Auth{Key: "k", Secret: "secret123"}
	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := platon.NewRequest(platon.ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(platon.HashTypeCardTokenPayment)

	c := NewClient(DefaultOptions())
	c.SetClient(&http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       nil,
			}, nil
		}),
	})

	_, err := c.Api(req, "https://example.com")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "response body is nil") && !strings.Contains(err.Error(), "empty response") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestApi_ReturnsDeclinedErrorFromReason(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"DECLINED","decline_reason":"102: Token is not active","error_message":null}`))
	}))
	defer srv.Close()

	auth := &platon.Auth{Key: "k", Secret: "secret123"}
	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := platon.NewRequest(platon.ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(platon.HashTypeCardTokenPayment)

	c := NewClient(DefaultOptions())
	resp, err := c.Api(req, srv.URL)
	if err == nil {
		t.Fatalf("expected decline error, got nil")
	}
	if !strings.Contains(err.Error(), "102: Token is not active") {
		t.Fatalf("expected decline reason in error, got %q", err.Error())
	}
	if resp == nil {
		t.Fatalf("expected response payload with decline_reason")
	}
	if resp.DeclineReason != "102: Token is not active" {
		t.Fatalf("unexpected decline reason: %q", resp.DeclineReason)
	}
}

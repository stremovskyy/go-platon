package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

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


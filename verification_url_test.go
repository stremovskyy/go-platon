package go_platon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stremovskyy/go-platon/platon"
)

func TestResolveClientServerVerificationURL_UsesLocationHeader(t *testing.T) {
	wantURL := "https://secure.platononline.com/payment/purchase?token=ABC123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method mismatch: want %q, got %q", http.MethodPost, r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("content-type mismatch: want application/x-www-form-urlencoded, got %q", got)
		}

		w.Header().Set("Location", wantURL)
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()

	form := &platon.ClientServerVerificationForm{
		Method:   http.MethodPost,
		Endpoint: server.URL,
		Fields: map[string]string{
			"payment": "CC",
			"key":     "client",
			"url":     "https://merchant.example/success",
			"data":    "payload",
			"sign":    "signature",
		},
	}

	urlResult, err := resolveClientServerVerificationURL(form)
	if err != nil {
		t.Fatalf("resolveClientServerVerificationURL() error: %v", err)
	}
	if urlResult.String() != wantURL {
		t.Fatalf("URL mismatch: want %q, got %q", wantURL, urlResult.String())
	}
}

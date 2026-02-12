package go_platon

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

func TestBuildClientServerVerificationForm(t *testing.T) {
	paymentID := "order-1"
	req := &Request{
		Merchant: &Merchant{
			MerchantKey:     "CLIENT_KEY",
			SecretKey:       "SECRET_KEY",
			SuccessRedirect: "https://merchant.example/success",
		},
		PaymentData: &PaymentData{
			PaymentID:   &paymentID,
			Currency:    currency.UAH,
			Description: "Verify card",
		},
	}

	form, err := BuildClientServerVerificationForm(req)
	if err != nil {
		t.Fatalf("BuildClientServerVerificationForm() error: %v", err)
	}

	if form.Method != "POST" {
		t.Fatalf("method mismatch: want POST, got %q", form.Method)
	}
	if form.Endpoint != consts.ApiPaymentAuthURL {
		t.Fatalf("endpoint mismatch: want %q, got %q", consts.ApiPaymentAuthURL, form.Endpoint)
	}

	fields := form.Fields
	if fields["payment"] != "CC" {
		t.Fatalf("payment mismatch: want CC, got %q", fields["payment"])
	}
	if fields["key"] != "CLIENT_KEY" {
		t.Fatalf("key mismatch: want CLIENT_KEY, got %q", fields["key"])
	}
	if fields["url"] != "https://merchant.example/success" {
		t.Fatalf("url mismatch: got %q", fields["url"])
	}
	if fields["formid"] != "verify" {
		t.Fatalf("formid mismatch: want verify, got %q", fields["formid"])
	}
	if fields["req_token"] != "Y" {
		t.Fatalf("req_token mismatch: want Y, got %q", fields["req_token"])
	}
	if fields["sign"] != "72e8c7944a9b9422b05e21ecbdce48bb" {
		t.Fatalf("sign mismatch: got %q", fields["sign"])
	}

	rawData, err := base64.StdEncoding.DecodeString(fields["data"])
	if err != nil {
		t.Fatalf("cannot decode data: %v", err)
	}

	var payload struct {
		Amount      string `json:"amount"`
		Description string `json:"description"`
		Currency    string `json:"currency"`
		Recurring   string `json:"recurring"`
		Order       string `json:"order"`
	}
	if err := json.Unmarshal(rawData, &payload); err != nil {
		t.Fatalf("cannot decode JSON payload: %v", err)
	}

	if payload.Amount != "0.40" {
		t.Fatalf("amount mismatch: want 0.40, got %q", payload.Amount)
	}
	if payload.Description != "Verify card" {
		t.Fatalf("description mismatch: got %q", payload.Description)
	}
	if payload.Currency != "UAH" {
		t.Fatalf("currency mismatch: got %q", payload.Currency)
	}
	if payload.Recurring != "Y" {
		t.Fatalf("recurring mismatch: got %q", payload.Recurring)
	}
	if payload.Order != "order-1" {
		t.Fatalf("order mismatch: got %q", payload.Order)
	}
}

func TestBuildClientServerVerificationForm_WithoutOrderID(t *testing.T) {
	req := &Request{
		Merchant: &Merchant{
			MerchantKey:     "CLIENT_KEY",
			SecretKey:       "SECRET_KEY",
			SuccessRedirect: "https://merchant.example/success",
		},
		PaymentData: &PaymentData{
			Currency:    currency.UAH,
			Description: "Verify card",
		},
	}

	form, err := BuildClientServerVerificationForm(req)
	if err != nil {
		t.Fatalf("BuildClientServerVerificationForm() error: %v", err)
	}

	rawData, err := base64.StdEncoding.DecodeString(form.Fields["data"])
	if err != nil {
		t.Fatalf("cannot decode data: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(rawData, &payload); err != nil {
		t.Fatalf("cannot decode JSON payload: %v", err)
	}

	if _, exists := payload["order"]; exists {
		t.Fatalf("unexpected order field in payload")
	}
}

func TestBuildClientServerVerificationForm_Validation(t *testing.T) {
	validPaymentID := "order-1"
	valid := &Request{
		Merchant: &Merchant{
			MerchantKey:     "CLIENT_KEY",
			SecretKey:       "SECRET_KEY",
			SuccessRedirect: "https://merchant.example/success",
		},
		PaymentData: &PaymentData{
			PaymentID:   &validPaymentID,
			Currency:    currency.UAH,
			Description: "Verify card",
		},
	}

	tests := []struct {
		name      string
		req       *Request
		wantError string
	}{
		{name: "nil request", req: nil, wantError: "request is nil"},
		{
			name: "missing merchant",
			req: &Request{
				PaymentData: valid.PaymentData,
			},
			wantError: "merchant is required",
		},
		{
			name: "missing merchant key",
			req: &Request{
				Merchant: &Merchant{
					SecretKey:       "SECRET_KEY",
					SuccessRedirect: "https://merchant.example/success",
				},
				PaymentData: valid.PaymentData,
			},
			wantError: "client_key is required",
		},
		{
			name: "missing secret key",
			req: &Request{
				Merchant: &Merchant{
					MerchantKey:     "CLIENT_KEY",
					SuccessRedirect: "https://merchant.example/success",
				},
				PaymentData: valid.PaymentData,
			},
			wantError: "secret key is required",
		},
		{
			name: "missing redirect URL",
			req: &Request{
				Merchant: &Merchant{
					MerchantKey: "CLIENT_KEY",
					SecretKey:   "SECRET_KEY",
				},
				PaymentData: valid.PaymentData,
			},
			wantError: "redirect URL is required",
		},
		{
			name: "missing description",
			req: &Request{
				Merchant: valid.Merchant,
				PaymentData: &PaymentData{
					Currency: currency.UAH,
				},
			},
			wantError: "order_description is required",
		},
		{
			name: "missing currency",
			req: &Request{
				Merchant: valid.Merchant,
				PaymentData: &PaymentData{
					Description: "Verify card",
				},
			},
			wantError: "order_currency is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := BuildClientServerVerificationForm(tc.req)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tc.req == nil {
				if err != platon.ErrRequestIsNil {
					t.Fatalf("expected ErrRequestIsNil, got %v", err)
				}
				return
			}
			if !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("error mismatch: want contains %q, got %q", tc.wantError, err.Error())
			}
		})
	}
}

func TestVerification_ValidatesRequestBeforeNetworkCall(t *testing.T) {
	c := &client{}
	req := &Request{}

	result, err := c.Verification(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result")
	}
	if !strings.Contains(err.Error(), "merchant is required") {
		t.Fatalf("error mismatch: got %q", err.Error())
	}
}

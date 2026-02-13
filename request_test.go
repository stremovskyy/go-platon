package go_platon

import (
	"math"
	"testing"
)

func TestRequest_GetAmount_UsesMinorUnits(t *testing.T) {
	req := &Request{
		PaymentData: &PaymentData{
			Amount: 199,
		},
	}

	got := req.GetAmount()
	if math.Abs(float64(got)-1.99) > 1e-6 {
		t.Fatalf("GetAmount() mismatch: want 1.99, got %.8f", got)
	}
}

func TestRequest_NilReceiver_GettersAreSafe(t *testing.T) {
	var req *Request

	if auth := req.GetAuth(); auth == nil || auth.Key != "EMPTY_KEY" || auth.Secret != "EMPTY_SECRET" {
		t.Fatalf("GetAuth() expected fallback auth, got %#v", auth)
	}
	if req.GetSuccessRedirect() != "" {
		t.Fatalf("GetSuccessRedirect() expected empty value")
	}
	if req.GetFailRedirect() != "" {
		t.Fatalf("GetFailRedirect() expected empty value")
	}
	if req.GetPlatonPaymentID() != 0 {
		t.Fatalf("GetPlatonPaymentID() expected zero value")
	}
	if req.GetPlatonTransID() != nil {
		t.Fatalf("GetPlatonTransID() expected nil")
	}
	if req.GetCardToken() != nil {
		t.Fatalf("GetCardToken() expected nil")
	}
	if req.GetCardPan() != nil {
		t.Fatalf("GetCardPan() expected nil")
	}
	if req.GetPaymentID() != nil {
		t.Fatalf("GetPaymentID() expected nil")
	}
	if req.GetPayerEmail() != nil {
		t.Fatalf("GetPayerEmail() expected nil")
	}
	if req.GetPayerPhone() != nil {
		t.Fatalf("GetPayerPhone() expected nil")
	}
	if req.GetAmount() != 0 {
		t.Fatalf("GetAmount() expected zero value")
	}
	if req.GetDescription() != "" {
		t.Fatalf("GetDescription() expected empty value")
	}
	if req.GetCurrency() != "" {
		t.Fatalf("GetCurrency() expected empty value")
	}
	if req.IsMobile() {
		t.Fatalf("IsMobile() expected false")
	}
	if _, err := req.GetAppleContainer(); err == nil {
		t.Fatalf("GetAppleContainer() expected error")
	}
	if req.IsApplePay() {
		t.Fatalf("IsApplePay() expected false")
	}
	if _, err := req.GetGoogleToken(); err == nil {
		t.Fatalf("GetGoogleToken() expected error")
	}
	if req.GetTrackingData() != nil {
		t.Fatalf("GetTrackingData() expected nil")
	}
	if splitRules, err := req.GetSplitRules(); err != nil || splitRules != nil {
		t.Fatalf("GetSplitRules() expected nil,nil, got %v,%v", splitRules, err)
	}
	if req.GetSubmerchantID() != nil {
		t.Fatalf("GetSubmerchantID() expected nil")
	}
	if req.GetReceiverTIN() != nil {
		t.Fatalf("GetReceiverTIN() expected nil")
	}
	if req.GetRelatedIDs() != nil {
		t.Fatalf("GetRelatedIDs() expected nil")
	}
	if metadata := req.GetMetadata(); metadata == nil || len(metadata) != 0 {
		t.Fatalf("GetMetadata() expected empty map")
	}
	if req.GetMerchantKey() != "" {
		t.Fatalf("GetMerchantKey() expected empty value")
	}
	if req.GetClientIP() != nil {
		t.Fatalf("GetClientIP() expected nil")
	}
	if req.GetTermsURL() != nil {
		t.Fatalf("GetTermsURL() expected nil")
	}
	if req.GetCardNumber() != nil {
		t.Fatalf("GetCardNumber() expected nil")
	}
	if req.GetCardExpMonth() != nil {
		t.Fatalf("GetCardExpMonth() expected nil")
	}
	if req.GetCardExpYear() != nil {
		t.Fatalf("GetCardExpYear() expected nil")
	}
	if req.GetCardCvv2() != nil {
		t.Fatalf("GetCardCvv2() expected nil")
	}
}

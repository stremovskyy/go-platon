package go_platon

import (
	"testing"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

func TestCredit_CardToken_DryRun_BuildsA2CRequest(t *testing.T) {
	var capturedEndpoint string
	var capturedRequest *platon.Request

	c := &client{}
	request := &Request{
		Merchant: &Merchant{
			MerchantKey: "CLIENT_KEY",
			SecretKey:   "CLIENT_PASS",
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("ORDER-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "A2C payout",
		},
		PaymentMethod: &PaymentMethod{
			Card: &Card{Token: ref("CARD_TOKEN")},
		},
	}

	_, err := c.Credit(request, DryRun(func(endpoint string, payload any) {
		capturedEndpoint = endpoint
		capturedRequest, _ = payload.(*platon.Request)
	}))
	if err != nil {
		t.Fatalf("Credit() unexpected error: %v", err)
	}

	if capturedEndpoint != consts.ApiP2PUnqURL {
		t.Fatalf("Credit() endpoint mismatch: want %q, got %q", consts.ApiP2PUnqURL, capturedEndpoint)
	}
	if capturedRequest == nil {
		t.Fatal("Credit() captured request is nil")
	}
	if capturedRequest.Action != platon.ActionCodeCREDIT2CARD.String() {
		t.Fatalf("Credit() action mismatch: want %q, got %q", platon.ActionCodeCREDIT2CARD.String(), capturedRequest.Action)
	}
	if capturedRequest.HashType != platon.HashTypeCredit2CardToken {
		t.Fatalf("Credit() hash type mismatch: want %q, got %q", platon.HashTypeCredit2CardToken, capturedRequest.HashType)
	}
	if capturedRequest.PayerFirstName == nil || *capturedRequest.PayerFirstName == "" {
		t.Fatal("Credit() payer_first_name should be filled")
	}
	if capturedRequest.PayerCountry == nil || *capturedRequest.PayerCountry == "" {
		t.Fatal("Credit() payer_country should be filled")
	}
}

func TestStatus_DryRun_A2CFlow_UsesP2PEndpointAndHash(t *testing.T) {
	var capturedEndpoint string
	var capturedRequest *platon.Request

	c := &client{}
	request := &Request{
		Merchant: &Merchant{
			MerchantKey: "CLIENT_KEY",
			SecretKey:   "CLIENT_PASS",
		},
		PaymentData: &PaymentData{
			PaymentID: ref("ORDER-2"),
			Metadata: map[string]string{
				platonMetaFlow: platonFlowA2C,
			},
		},
	}

	_, err := c.Status(request, DryRun(func(endpoint string, payload any) {
		capturedEndpoint = endpoint
		capturedRequest, _ = payload.(*platon.Request)
	}))
	if err != nil {
		t.Fatalf("Status() unexpected error: %v", err)
	}

	if capturedEndpoint != consts.ApiP2PUnqURL {
		t.Fatalf("Status() endpoint mismatch: want %q, got %q", consts.ApiP2PUnqURL, capturedEndpoint)
	}
	if capturedRequest == nil {
		t.Fatal("Status() captured request is nil")
	}
	if capturedRequest.HashType != platon.HashTypeGetTransStatusByOrder {
		t.Fatalf("Status() hash type mismatch: want %q, got %q", platon.HashTypeGetTransStatusByOrder, capturedRequest.HashType)
	}
	if capturedRequest.Action != platon.ActionCodeGetTransStatusByOrder.String() {
		t.Fatalf("Status() action mismatch: want %q, got %q", platon.ActionCodeGetTransStatusByOrder.String(), capturedRequest.Action)
	}
}

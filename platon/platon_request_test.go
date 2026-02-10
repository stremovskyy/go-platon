package platon

import (
	"testing"

	"github.com/stremovskyy/go-platon/currency"
)

func TestSignAndPrepare_VerificationSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "verify-1"
	desc := "verification"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	pan := "4111111111111111"
	month := "01"
	year := "2026"
	cvv := "123"

	req := NewRequest(ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithChannelNoAmountVerification().
		WithOrderID(&orderID).
		WithOrderAmount(VerifyNoAmount.String()).
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithCardNumber(&pan).
		WithCardExpMonth(&month).
		WithCardExpYear(&year).
		WithCardCvv2(&cvv).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		WithReqToken(true).
		WithRecurringInitFlag(true).
		SignForAction(HashTypeVerification)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "bcc927a61aee5b183d13f1154e2ea5e2"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_CardPaymentSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"
	desc := "payment"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	pan := "4111111111111111"
	month := "01"
	year := "2026"
	cvv := "123"

	req := NewRequest(ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithCardNumber(&pan).
		WithCardExpMonth(&month).
		WithCardExpYear(&year).
		WithCardCvv2(&cvv).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(HashTypeCardPayment)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	// Same signature scheme as verification (email + secret + first6/last4).
	const want = "bcc927a61aee5b183d13f1154e2ea5e2"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}

	if signed.ReqToken == nil || *signed.ReqToken != "N" {
		t.Fatalf("req_token default mismatch: want N, got %v", signed.ReqToken)
	}
	if signed.RecurringInit == nil || *signed.RecurringInit != "N" {
		t.Fatalf("recurring_init default mismatch: want N, got %v", signed.RecurringInit)
	}
}

func TestSignAndPrepare_CardTokenPaymentSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"
	desc := "one-click"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	token := "TOKEN123"

	req := NewRequest(ActionCodeSALE).
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
		SignForAction(HashTypeCardTokenPayment)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "03838ac02c89b98621f95ec98a68aa14"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_ApplePaySignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"
	desc := "apple"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	data := "ZGF0YQ==" // "data" base64

	req := NewRequest(ActionCodeAPPLEPAY).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		WithApplePayData(&data).
		SignForAction(HashTypeApplePay)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "02d1662d7a7eb526b1c939639a914ec6"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_RecurringSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"
	desc := "recurring"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	token := "TOKEN123"
	ext3 := "recurring"

	req := NewRequest(ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithCardToken(&token).
		WithExt3(&ext3).
		WithOrderID(&orderID).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithPayerEmail(&email).
		SignForAction(HashTypeRecurring)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "03838ac02c89b98621f95ec98a68aa14"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_GetTransStatusSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"
	cardHashPart := "4111111111" // first6+last4

	req := NewRequest(ActionCodeGetTransStatus).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithPayerEmail(&email).
		SignForAction(HashTypeGetTransStatus)
	req.CardHashPart = &cardHashPart

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "77a4785689636b4d3875ec7acf47d5e2"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_CaptureSignatureAndMap(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"
	cardHashPart := "4111111111" // first6+last4

	req := NewRequest(ActionCodeCAPTURE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithHashEmail(&email).
		SignForAction(HashTypeCapture).
		WithCardHashPart(&cardHashPart)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "77a4785689636b4d3875ec7acf47d5e2"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}

	// Internal-only fields must not be serialized.
	m := signed.ToMap()
	if _, ok := m["hash_email"]; ok {
		t.Fatalf("unexpected serialized key: hash_email")
	}
	if _, ok := m["card_hash_part"]; ok {
		t.Fatalf("unexpected serialized key: card_hash_part")
	}
	if _, ok := m["amount"]; !ok {
		t.Fatalf("expected serialized key: amount")
	}
	if _, ok := m["split_rules"]; ok {
		t.Fatalf("unexpected serialized key: split_rules")
	}
}

func TestSignAndPrepare_CreditVoidSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"
	cardHashPart := "4111111111" // first6+last4

	req := NewRequest(ActionCodeCREDITVOID).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithHashEmail(&email).
		SignForAction(HashTypeCreditVoid).
		WithCardHashPart(&cardHashPart)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "77a4785689636b4d3875ec7acf47d5e2"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_CaptureWithSplitRules(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"
	cardHashPart := "4111111111" // first6+last4

	req := NewRequest(ActionCodeCAPTURE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("10.00").
		WithSplitRules(SplitRules{
			"submerchant_01": "2.50",
			"submerchant_02": "7.50",
		}).
		WithHashEmail(&email).
		SignForAction(HashTypeCapture).
		WithCardHashPart(&cardHashPart)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	m := signed.ToMap()
	if _, ok := m["split_rules"]; !ok {
		t.Fatalf("expected serialized key: split_rules")
	}
}

func TestSignAndPrepare_CreditVoidSplitRulesExceedAmount(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"
	cardHashPart := "4111111111" // first6+last4

	req := NewRequest(ActionCodeCREDITVOID).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithSplitRules(SplitRules{
			"submerchant_01": "0.70",
			"submerchant_02": "0.40",
		}).
		WithHashEmail(&email).
		SignForAction(HashTypeCreditVoid).
		WithCardHashPart(&cardHashPart)

	if _, err := req.SignAndPrepare(); err == nil {
		t.Fatalf("expected split rules validation error, got nil")
	}
}

func TestSignAndPrepare_GetSubmerchantSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}
	submerchantID := "12345678"

	req := NewRequest(ActionCodeGetSubmerchant).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithSubmerchantID(&submerchantID).
		SignForAction(HashTypeGetSubmerchant)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "15f549d19f26ce89022396a649c4ac9f"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_OrderAmountValidation(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"
	desc := "payment"
	ip := "127.0.0.1"
	term := "https://example.com/3ds"
	email := "payer@example.com"
	phone := "380631234567"
	pan := "4111111111111111"
	month := "01"
	year := "2026"
	cvv := "123"

	req := NewRequest(ActionCodeSALE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		WithOrderAmount("1000"). // invalid format (must be 1000.00)
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithPayerIP(&ip).
		WithTermsURL(&term).
		WithCardNumber(&pan).
		WithCardExpMonth(&month).
		WithCardExpYear(&year).
		WithCardCvv2(&cvv).
		WithPayerEmail(&email).
		WithPayerPhone(&phone).
		SignForAction(HashTypeCardPayment)

	if _, err := req.SignAndPrepare(); err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

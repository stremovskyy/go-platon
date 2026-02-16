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

	req := NewRequest(ActionCodeGetTransStatus).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithPayerEmail(&email).
		SignForAction(HashTypeGetTransStatus)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "ef374c28b6398c097e0b3d6230deebd6"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_CaptureSignatureAndMap(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"

	req := NewRequest(ActionCodeCAPTURE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithHashEmail(&email).
		SignForAction(HashTypeCapture)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "ef374c28b6398c097e0b3d6230deebd6"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}

	// Internal-only fields must not be serialized.
	m := signed.ToMap()
	if _, ok := m["hash_email"]; ok {
		t.Fatalf("unexpected serialized key: hash_email")
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

	req := NewRequest(ActionCodeCREDITVOID).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithHashEmail(&email).
		SignForAction(HashTypeCreditVoid)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "ef374c28b6398c097e0b3d6230deebd6"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_Credit2CardSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-a2c-pan"
	desc := "a2c payout"
	pan := "4111111111111111"
	firstName := "John"
	lastName := "Doe"
	address := "Main st 1"
	country := "UA"
	state := "UA"
	city := "Kyiv"
	zip := "01001"

	req := NewRequest(ActionCodeCREDIT2CARD).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		WithAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithCardNumber(&pan).
		WithPayerFirstName(&firstName).
		WithPayerLastName(&lastName).
		WithPayerAddress(&address).
		WithPayerCountry(&country).
		WithPayerState(&state).
		WithPayerCity(&city).
		WithPayerZip(&zip).
		SignForAction(HashTypeCredit2Card)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "cbe775dd3121bd75d6636a42a3cf65cc"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_Credit2CardTokenSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-a2c-token"
	desc := "a2c payout"
	token := "TOKEN123"
	firstName := "John"
	lastName := "Doe"
	address := "Main st 1"
	country := "UA"
	state := "UA"
	city := "Kyiv"
	zip := "01001"

	req := NewRequest(ActionCodeCREDIT2CARD).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		WithAmount("1.00").
		ForCurrency(currency.UAH).
		WithDescription(desc).
		WithCardToken(&token).
		WithPayerFirstName(&firstName).
		WithPayerLastName(&lastName).
		WithPayerAddress(&address).
		WithPayerCountry(&country).
		WithPayerState(&state).
		WithPayerCity(&city).
		WithPayerZip(&zip).
		SignForAction(HashTypeCredit2CardToken)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "9d63d6b5b3de7807899d10e08f00864a"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_GetTransStatusByOrderSignature(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	orderID := "order-123"

	req := NewRequest(ActionCodeGetTransStatusByOrder).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithOrderID(&orderID).
		SignForAction(HashTypeGetTransStatusByOrder)

	signed, err := req.SignAndPrepare()
	if err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}

	const want = "b6a84d3306211abea3704548513662d6"
	if signed.Hash != want {
		t.Fatalf("hash mismatch: want %s, got %s", want, signed.Hash)
	}
}

func TestSignAndPrepare_CaptureWithSplitRules(t *testing.T) {
	auth := &Auth{Key: "k", Secret: "secret123"}

	email := "payer@example.com"
	transID := "632508054"

	req := NewRequest(ActionCodeCAPTURE).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("10.00").
		WithSplitRules(
			SplitRules{
				"submerchant_01": "2.50",
				"submerchant_02": "7.50",
			},
		).
		WithHashEmail(&email).
		SignForAction(HashTypeCapture)

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

	req := NewRequest(ActionCodeCREDITVOID).
		WithAuth(auth).
		WithClientKey("clientKey").
		WithTransID(&transID).
		WithAmount("1.00").
		WithSplitRules(
			SplitRules{
				"submerchant_01": "0.70",
				"submerchant_02": "0.40",
			},
		).
		WithHashEmail(&email).
		SignForAction(HashTypeCreditVoid)

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

func TestRequest_NilReceiver_SignAndPrepare(t *testing.T) {
	var req *Request

	signed, err := req.SignAndPrepare()
	if err == nil {
		t.Fatalf("expected error for nil request receiver, got nil")
	}
	if signed != nil {
		t.Fatalf("expected nil signed request, got %#v", signed)
	}
}

func TestRequest_NilReceiver_SignForAction(t *testing.T) {
	var req *Request
	if got := req.SignForAction(HashTypeCardPayment); got != nil {
		t.Fatalf("expected nil receiver to stay nil, got %#v", got)
	}
}

func TestRequest_NilReceiver_ToMap(t *testing.T) {
	var req *Request

	result := req.ToMap()
	if result == nil {
		t.Fatalf("expected non-nil map for nil receiver")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestRequest_NilReceiver_BuilderChainIsSafe(t *testing.T) {
	var req *Request

	orderID := "order-1"
	transID := "trans-1"
	email := "payer@example.com"
	value := "value"

	got := req.
		WithAuth(&Auth{Key: "k", Secret: "s"}).
		WithClientKey("k").
		WithReqToken(true).
		WithRecToken().
		WithRecurringInitFlag(true).
		WithRecurringInit().
		WithAsync(true).
		UseAsync().
		WithChannelNoAmountVerification().
		WithPayerIP(nil).
		WithTermsURL(&value).
		WithCardNumber(&value).
		WithCardToken(&value).
		WithCardExpMonth(&value).
		WithCardExpYear(&value).
		WithCardCvv2(&value).
		WithPayerEmail(&email).
		WithPayerPhone(&value).
		WithPayerFirstName(&value).
		WithPayerLastName(&value).
		WithApplePayData(&value).
		WithGooglePayToken(&value).
		WithPaymentToken(&value).
		WithHoldAuth().
		WithVerifyAmount(0).
		WithOrderAmountMinorUnits(100).
		WithOrderAmount("1.00").
		ForCurrency(currency.UAH).
		WithSubmerchantID(&value).
		WithDescription("desc").
		WithOrderID(&orderID).
		WithRecurringFirstTransID(&transID).
		WithTransID(&transID).
		WithAmountMinorUnits(100).
		WithAmount("1.00").
		WithSplitRules(SplitRules{"submerchant": "1.00"}).
		WithImmediately(true).
		WithHashEmail(&email).
		WithExt3(&value).
		SignForAction(HashTypeCardPayment)

	if got != nil {
		t.Fatalf("expected nil request after nil receiver builder chain, got %#v", got)
	}
}

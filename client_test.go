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

package go_platon

import (
	"encoding/base64"
	"testing"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

func ref(s string) *string { return &s }

func TestBuildIAPaymentRequest_ApplePay(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	// Minimal Apple Pay container for GetAppleContainer(): it extracts top-level "token".
	containerJSON := `{"token":{"foo":"bar"}}`
	containerB64 := base64.StdEncoding.EncodeToString([]byte(containerJSON))

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			AppleContainer: &containerB64,
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	apiReq, apiURL, err := c.buildIAPaymentRequest(req, false)
	if err != nil {
		t.Fatalf("buildIAPaymentRequest() error: %v", err)
	}

	if apiURL != consts.ApiPostURL {
		t.Fatalf("apiURL mismatch: want %q, got %q", consts.ApiPostURL, apiURL)
	}
	if apiReq.Action != platon.ActionCodeAPPLEPAY.String() {
		t.Fatalf("action mismatch: want %q, got %q", platon.ActionCodeAPPLEPAY.String(), apiReq.Action)
	}
	if apiReq.HashType != platon.HashTypeApplePay {
		t.Fatalf("hash type mismatch: want %q, got %q", platon.HashTypeApplePay, apiReq.HashType)
	}
	if apiReq.PaymentToken == nil || *apiReq.PaymentToken == "" {
		t.Fatalf("payment_token must be set for Apple Pay")
	}

	if _, err := apiReq.SignAndPrepare(); err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}
}

func TestBuildIAPaymentRequest_ApplePay_WithSplitRules(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	containerJSON := `{"token":{"foo":"bar"}}`
	containerB64 := base64.StdEncoding.EncodeToString([]byte(containerJSON))

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			AppleContainer: &containerB64,
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
			SplitRules: []SplitRule{
				{SubmerchantIdentification: "submerchant_01", Amount: 100},
			},
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	apiReq, _, err := c.buildIAPaymentRequest(req, false)
	if err != nil {
		t.Fatalf("buildIAPaymentRequest() error: %v", err)
	}

	if apiReq.SplitRules["submerchant_01"] != "1.00" {
		t.Fatalf("split_rules[\"submerchant_01\"] mismatch: want 1.00, got %s", apiReq.SplitRules["submerchant_01"])
	}
}

func TestBuildIAPaymentRequest_GooglePay(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	googleTokenJSON := `{"paymentMethodData":{"tokenizationData":{"token":"{\\\"foo\\\":\\\"bar\\\"}"}}}`
	googleTokenB64 := base64.StdEncoding.EncodeToString([]byte(googleTokenJSON))

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			GoogleToken: &googleTokenB64,
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	apiReq, apiURL, err := c.buildIAPaymentRequest(req, false)
	if err != nil {
		t.Fatalf("buildIAPaymentRequest() error: %v", err)
	}

	if apiURL != consts.ApiPostURL {
		t.Fatalf("apiURL mismatch: want %q, got %q", consts.ApiPostURL, apiURL)
	}
	if apiReq.Action != platon.ActionCodeGOOGLEPAY.String() {
		t.Fatalf("action mismatch: want %q, got %q", platon.ActionCodeGOOGLEPAY.String(), apiReq.Action)
	}
	if apiReq.HashType != platon.HashTypeGooglePay {
		t.Fatalf("hash type mismatch: want %q, got %q", platon.HashTypeGooglePay, apiReq.HashType)
	}
	if apiReq.PaymentToken == nil || *apiReq.PaymentToken == "" {
		t.Fatalf("payment_token must be set for Google Pay")
	}

	if _, err := apiReq.SignAndPrepare(); err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}
}

func TestBuildIAPaymentRequest_CardToken(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			Card: &Card{Token: ref("CARD_TOKEN")},
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
		},
	}

	c := &client{}
	apiReq, apiURL, err := c.buildIAPaymentRequest(req, false)
	if err != nil {
		t.Fatalf("buildIAPaymentRequest() error: %v", err)
	}

	if apiURL != consts.ApiPostUnqURL {
		t.Fatalf("apiURL mismatch: want %q, got %q", consts.ApiPostUnqURL, apiURL)
	}
	if apiReq.Action != platon.ActionCodeSALE.String() {
		t.Fatalf("action mismatch: want %q, got %q", platon.ActionCodeSALE.String(), apiReq.Action)
	}
	if apiReq.HashType != platon.HashTypeCardTokenPayment {
		t.Fatalf("hash type mismatch: want %q, got %q", platon.HashTypeCardTokenPayment, apiReq.HashType)
	}
	if apiReq.ReqToken != nil || apiReq.RecurringInit != nil {
		t.Fatalf("req_token/recurring_init must not be set for CARD_TOKEN payment")
	}

	if _, err := apiReq.SignAndPrepare(); err != nil {
		t.Fatalf("SignAndPrepare() error: %v", err)
	}
}

func TestBuildIAPaymentRequest_CardToken_WithMetadataExtFields(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			Card: &Card{Token: ref("CARD_TOKEN")},
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
			Metadata: map[string]string{
				"ext1":  " merchant-core ",
				"ext2":  "   ",
				"ext4":  "wallet-topup",
				"ext10": "v1",
			},
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
		},
	}

	c := &client{}
	apiReq, _, err := c.buildIAPaymentRequest(req, false)
	if err != nil {
		t.Fatalf("buildIAPaymentRequest() error: %v", err)
	}

	if apiReq.Ext1 == nil || *apiReq.Ext1 != "merchant-core" {
		t.Fatalf("ext1 mismatch: got %#v", apiReq.Ext1)
	}
	if apiReq.Ext2 != nil {
		t.Fatalf("ext2 must be nil for blank metadata value, got %#v", apiReq.Ext2)
	}
	if apiReq.Ext4 == nil || *apiReq.Ext4 != "wallet-topup" {
		t.Fatalf("ext4 mismatch: got %#v", apiReq.Ext4)
	}
	if apiReq.Ext10 == nil || *apiReq.Ext10 != "v1" {
		t.Fatalf("ext10 mismatch: got %#v", apiReq.Ext10)
	}
}

func TestBuildIAPaymentRequest_CardPAN_IsNotSupported(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			Card: &Card{
				Pan:             ref("4111111111111111"),
				ExpirationMonth: ref("01"),
				ExpirationYear:  ref("2026"),
				Cvv2:            ref("123"),
			},
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "desc",
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	if _, _, err := c.buildIAPaymentRequest(req, false); err == nil {
		t.Fatalf("buildIAPaymentRequest() expected error for PAN payment, got nil")
	}
}

func TestBuildIAPaymentRequest_CardPAN_WithSplitRules_IsNotSupported(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			Card: &Card{
				Pan:             ref("4111111111111111"),
				ExpirationMonth: ref("01"),
				ExpirationYear:  ref("2026"),
				Cvv2:            ref("123"),
			},
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      10000,
			Currency:    currency.UAH,
			Description: "desc",
			SplitRules: []SplitRule{
				{SubmerchantIdentification: "submerchant_01", Amount: 2500},
				{SubmerchantIdentification: "submerchant_02", Amount: 7500},
			},
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	if _, _, err := c.buildIAPaymentRequest(req, false); err == nil {
		t.Fatalf("buildIAPaymentRequest() expected error for PAN payment, got nil")
	}
}

func TestBuildIAPaymentRequest_CardPAN_SplitRulesTotalExceedsAmount(t *testing.T) {
	merchant := &Merchant{
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    ref("https://example.com/3ds"),
	}

	req := &Request{
		Merchant: merchant,
		PaymentMethod: &PaymentMethod{
			Card: &Card{
				Pan:             ref("4111111111111111"),
				ExpirationMonth: ref("01"),
				ExpirationYear:  ref("2026"),
				Cvv2:            ref("123"),
			},
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      1000,
			Currency:    currency.UAH,
			Description: "desc",
			SplitRules: []SplitRule{
				{SubmerchantIdentification: "submerchant_01", Amount: 600},
				{SubmerchantIdentification: "submerchant_02", Amount: 500},
			},
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
			Phone: ref("380631234567"),
		},
	}

	c := &client{}
	if _, _, err := c.buildIAPaymentRequest(req, false); err == nil {
		t.Fatalf("expected split rules validation error, got nil")
	}
}

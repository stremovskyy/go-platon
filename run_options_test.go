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
	"testing"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/internal/utils"
	"github.com/stremovskyy/go-platon/platon"
)

func TestPayment_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	var (
		gotEndpoint string
		gotPayload  any
	)

	resp, err := cl.Payment(
		&Request{
			Merchant: &Merchant{
				MerchantKey: "clientKey",
				SecretKey:   "secret123",
				TermsURL:    utils.Ref("https://merchant.example/3ds"),
			},
			PaymentData: &PaymentData{
				PaymentID:   utils.Ref("order-1"),
				Amount:      100,
				Currency:    currency.UAH,
				Description: "dry-run",
			},
			PaymentMethod: &PaymentMethod{
				Card: &Card{
					Token: utils.Ref("CARD_TOKEN"),
				},
			},
			PersonalData: &PersonalData{
				Email: utils.Ref("payer@example.com"),
			},
		}, DryRun(
			func(endpoint string, payload any) {
				gotEndpoint = endpoint
				gotPayload = payload
			},
		),
	)

	if err != nil {
		t.Fatalf("Payment() dry run error: %v", err)
	}
	if resp != nil {
		t.Fatalf("Payment() dry run response: expected nil, got %#v", resp)
	}
	if gotEndpoint != consts.ApiPostUnqURL {
		t.Fatalf("endpoint mismatch: want %q, got %q", consts.ApiPostUnqURL, gotEndpoint)
	}

	req, ok := gotPayload.(*platon.Request)
	if !ok {
		t.Fatalf("payload type mismatch: got %T", gotPayload)
	}
	if req.Action != platon.ActionCodeSALE.String() {
		t.Fatalf("action mismatch: want %q, got %q", platon.ActionCodeSALE.String(), req.Action)
	}
}

func TestVerification_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	var (
		gotEndpoint string
		gotPayload  any
	)

	resp, err := cl.Verification(
		&Request{
			Merchant: &Merchant{
				MerchantKey:     "clientKey",
				SecretKey:       "secret123",
				SuccessRedirect: "https://merchant.example/success",
			},
			PaymentData: &PaymentData{
				PaymentID:   utils.Ref("order-1"),
				Currency:    currency.UAH,
				Description: "verify",
			},
		}, DryRun(
			func(endpoint string, payload any) {
				gotEndpoint = endpoint
				gotPayload = payload
			},
		),
	)

	if err != nil {
		t.Fatalf("Verification() dry run error: %v", err)
	}
	if resp != nil {
		t.Fatalf("Verification() dry run result: expected nil, got %#v", resp)
	}
	if gotEndpoint != consts.ApiPaymentAuthURL {
		t.Fatalf("endpoint mismatch: want %q, got %q", consts.ApiPaymentAuthURL, gotEndpoint)
	}

	form, ok := gotPayload.(*platon.ClientServerVerificationForm)
	if !ok {
		t.Fatalf("payload type mismatch: got %T", gotPayload)
	}
	if form.Endpoint != consts.ApiPaymentAuthURL {
		t.Fatalf("form endpoint mismatch: want %q, got %q", consts.ApiPaymentAuthURL, form.Endpoint)
	}
}

func TestVerificationLink_DryRun(t *testing.T) {
	cl := NewDefaultClient()

	var gotEndpoint string
	_, err := cl.VerificationLink(
		&Request{
			Merchant: &Merchant{
				MerchantKey:     "clientKey",
				SecretKey:       "secret123",
				SuccessRedirect: "https://merchant.example/success",
			},
			PaymentData: &PaymentData{
				PaymentID:   utils.Ref("order-1"),
				Currency:    currency.UAH,
				Description: "verify",
			},
		}, DryRun(
			func(endpoint string, _ any) {
				gotEndpoint = endpoint
			},
		),
	)

	if err != nil {
		t.Fatalf("VerificationLink() dry run error: %v", err)
	}
	if gotEndpoint != consts.ApiPaymentAuthURL {
		t.Fatalf("endpoint mismatch: want %q, got %q", consts.ApiPaymentAuthURL, gotEndpoint)
	}
}

func TestDryRun_DefaultHandler_NilPlatonRequestPayload(t *testing.T) {
	opts := collectRunOptions([]RunOption{DryRun()})
	var req *platon.Request

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("default dry-run handler panicked on nil *platon.Request payload: %v", r)
		}
	}()

	opts.handleDryRun(consts.ApiPostUnqURL, req)
}

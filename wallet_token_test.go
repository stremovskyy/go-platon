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
)

func walletRef(value string) *string {
	return &value
}

func TestRequest_GetApplePayToken_RawToken(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`

	req := &Request{
		PaymentMethod: &PaymentMethod{
			ApplePayToken: walletRef(token),
		},
	}

	got, err := req.GetApplePayToken()
	if err != nil {
		t.Fatalf("GetApplePayToken() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetApplePayToken() mismatch: want %q, got %#v", token, got)
	}
	if !req.IsApplePay() || !req.IsMobile() {
		t.Fatalf("expected Apple Pay request to be detected as mobile wallet")
	}
}

func TestRequest_GetApplePayToken_FullPaymentObject(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`
	payment := `{"token":` + token + `,"billingContact":{"givenName":"Ada"}}`

	req := &Request{
		PaymentMethod: &PaymentMethod{
			ApplePayPayment: walletRef(payment),
		},
	}

	got, err := req.GetApplePayToken()
	if err != nil {
		t.Fatalf("GetApplePayToken() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetApplePayToken() mismatch: want %q, got %#v", token, got)
	}
}

func TestRequest_GetApplePayToken_LegacyBase64(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`
	payment := `{"token":` + token + `}`
	legacy := base64.StdEncoding.EncodeToString([]byte(payment))

	req := &Request{
		PaymentMethod: &PaymentMethod{
			AppleContainer: walletRef(legacy),
		},
	}

	got, err := req.GetAppleContainer()
	if err != nil {
		t.Fatalf("GetAppleContainer() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetAppleContainer() mismatch: want %q, got %#v", token, got)
	}
}

func TestRequest_GetGooglePayToken_RawToken(t *testing.T) {
	token := `{"protocolVersion":"ECv2","signature":"sig","signedMessage":"{\"encryptedMessage\":\"abc\",\"ephemeralPublicKey\":\"def\",\"tag\":\"ghi\"}"}`

	req := &Request{
		PaymentMethod: &PaymentMethod{
			GooglePayToken: walletRef(token),
		},
	}

	got, err := req.GetGooglePayToken()
	if err != nil {
		t.Fatalf("GetGooglePayToken() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetGooglePayToken() mismatch: want %q, got %#v", token, got)
	}
	if !req.IsGooglePay() || !req.IsMobile() {
		t.Fatalf("expected Google Pay request to be detected as mobile wallet")
	}
}

func TestRequest_GetGooglePayToken_FullPaymentData(t *testing.T) {
	token := `{"protocolVersion":"ECv2","signature":"sig","signedMessage":"{\"encryptedMessage\":\"abc\",\"ephemeralPublicKey\":\"def\",\"tag\":\"ghi\"}"}`
	paymentData := `{"apiVersion":2,"paymentMethodData":{"tokenizationData":{"type":"PAYMENT_GATEWAY","token":"` +
		"{\\\"protocolVersion\\\":\\\"ECv2\\\",\\\"signature\\\":\\\"sig\\\",\\\"signedMessage\\\":\\\"{\\\\\\\"encryptedMessage\\\\\\\":\\\\\\\"abc\\\\\\\",\\\\\\\"ephemeralPublicKey\\\\\\\":\\\\\\\"def\\\\\\\",\\\\\\\"tag\\\\\\\":\\\\\\\"ghi\\\\\\\"}\\\"}" +
		`"}}}`

	req := &Request{
		PaymentMethod: &PaymentMethod{
			GooglePayPaymentData: walletRef(paymentData),
		},
	}

	got, err := req.GetGooglePayToken()
	if err != nil {
		t.Fatalf("GetGooglePayToken() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetGooglePayToken() mismatch: want %q, got %#v", token, got)
	}
}

func TestRequest_GetGooglePayToken_LegacyBase64(t *testing.T) {
	token := `{"protocolVersion":"ECv2","signature":"sig","signedMessage":"{\"encryptedMessage\":\"abc\",\"ephemeralPublicKey\":\"def\",\"tag\":\"ghi\"}"}`
	paymentData := `{"paymentMethodData":{"tokenizationData":{"token":"` +
		"{\\\"protocolVersion\\\":\\\"ECv2\\\",\\\"signature\\\":\\\"sig\\\",\\\"signedMessage\\\":\\\"{\\\\\\\"encryptedMessage\\\\\\\":\\\\\\\"abc\\\\\\\",\\\\\\\"ephemeralPublicKey\\\\\\\":\\\\\\\"def\\\\\\\",\\\\\\\"tag\\\\\\\":\\\\\\\"ghi\\\\\\\"}\\\"}" +
		`"}}}`
	legacy := base64.StdEncoding.EncodeToString([]byte(paymentData))

	req := &Request{
		PaymentMethod: &PaymentMethod{
			GoogleToken: walletRef(legacy),
		},
	}

	got, err := req.GetGoogleToken()
	if err != nil {
		t.Fatalf("GetGoogleToken() error: %v", err)
	}
	if got == nil || *got != token {
		t.Fatalf("GetGoogleToken() mismatch: want %q, got %#v", token, got)
	}
}

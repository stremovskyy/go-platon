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

func TestNewApplePayMethod_FromToken(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`

	method, err := NewApplePayMethod(token)
	if err != nil {
		t.Fatalf("NewApplePayMethod() error: %v", err)
	}
	if method == nil || method.ApplePayToken == nil || *method.ApplePayToken != token {
		t.Fatalf("NewApplePayMethod() mismatch: got %#v", method)
	}
}

func TestNewApplePayMethod_FromLegacyBase64(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`
	payment := `{"token":` + token + `}`
	legacy := base64.StdEncoding.EncodeToString([]byte(payment))

	method, err := NewApplePayMethod(legacy)
	if err != nil {
		t.Fatalf("NewApplePayMethod() error: %v", err)
	}
	if method == nil || method.ApplePayToken == nil || *method.ApplePayToken != token {
		t.Fatalf("NewApplePayMethod() mismatch: got %#v", method)
	}
}

func TestNewGooglePayMethod_FromPaymentData(t *testing.T) {
	token := `{"protocolVersion":"ECv2","signature":"sig","signedMessage":"{\"encryptedMessage\":\"abc\",\"ephemeralPublicKey\":\"def\",\"tag\":\"ghi\"}"}`
	paymentData := `{"apiVersion":2,"paymentMethodData":{"tokenizationData":{"type":"PAYMENT_GATEWAY","token":"` +
		"{\\\"protocolVersion\\\":\\\"ECv2\\\",\\\"signature\\\":\\\"sig\\\",\\\"signedMessage\\\":\\\"{\\\\\\\"encryptedMessage\\\\\\\":\\\\\\\"abc\\\\\\\",\\\\\\\"ephemeralPublicKey\\\\\\\":\\\\\\\"def\\\\\\\",\\\\\\\"tag\\\\\\\":\\\\\\\"ghi\\\\\\\"}\\\"}" +
		`"}}}`

	method, err := NewGooglePayMethod(paymentData)
	if err != nil {
		t.Fatalf("NewGooglePayMethod() error: %v", err)
	}
	if method == nil || method.GooglePayToken == nil || *method.GooglePayToken != token {
		t.Fatalf("NewGooglePayMethod() mismatch: got %#v", method)
	}
}

func TestNewGooglePayMethod_RejectsEmpty(t *testing.T) {
	if _, err := NewGooglePayMethod("   "); err == nil {
		t.Fatalf("NewGooglePayMethod() expected error for empty payload")
	}
}

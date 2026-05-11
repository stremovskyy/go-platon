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
	"encoding/json"
	"testing"
)

func assertBase64JSON(t *testing.T, got *string, wantJSON string) {
	t.Helper()

	if got == nil {
		t.Fatal("expected payment_token to be set")
	}

	decoded, err := base64.StdEncoding.DecodeString(*got)
	if err != nil {
		t.Fatalf("payment_token must be base64: %v", err)
	}

	var gotObject any
	if err := json.Unmarshal(decoded, &gotObject); err != nil {
		t.Fatalf("decoded payment_token must be JSON: %v", err)
	}

	var wantObject any
	if err := json.Unmarshal([]byte(wantJSON), &wantObject); err != nil {
		t.Fatalf("want JSON is invalid: %v", err)
	}

	gotCanonical, err := json.Marshal(gotObject)
	if err != nil {
		t.Fatalf("marshal got JSON: %v", err)
	}
	wantCanonical, err := json.Marshal(wantObject)
	if err != nil {
		t.Fatalf("marshal want JSON: %v", err)
	}

	if string(gotCanonical) != string(wantCanonical) {
		t.Fatalf("decoded payment_token mismatch: want %s, got %s", wantCanonical, gotCanonical)
	}
}

func TestNewApplePayMethod_FromToken(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`

	method, err := NewApplePayMethod(token)
	if err != nil {
		t.Fatalf("NewApplePayMethod() error: %v", err)
	}
	if method == nil {
		t.Fatalf("NewApplePayMethod() mismatch: got %#v", method)
	}
	assertBase64JSON(t, method.ApplePayToken, token)
}

func TestNewApplePayMethod_FromLegacyBase64(t *testing.T) {
	token := `{"paymentData":{"version":"EC_v1","data":"abc"},"paymentMethod":{"network":"Visa"},"transactionIdentifier":"tx-1"}`
	payment := `{"token":` + token + `}`
	legacy := base64.StdEncoding.EncodeToString([]byte(payment))

	method, err := NewApplePayMethod(legacy)
	if err != nil {
		t.Fatalf("NewApplePayMethod() error: %v", err)
	}
	if method == nil {
		t.Fatalf("NewApplePayMethod() mismatch: got %#v", method)
	}
	assertBase64JSON(t, method.ApplePayToken, token)
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
	if method == nil {
		t.Fatalf("NewGooglePayMethod() mismatch: got %#v", method)
	}
	assertBase64JSON(t, method.GooglePayToken, token)
}

func TestNewGooglePayMethod_RejectsEmpty(t *testing.T) {
	if _, err := NewGooglePayMethod("   "); err == nil {
		t.Fatalf("NewGooglePayMethod() expected error for empty payload")
	}
}

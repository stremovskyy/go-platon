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
	"io"
	"net/http"
	"strings"
	"testing"
)

type splitRoundTripFunc func(*http.Request) (*http.Response, error)

func (f splitRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSubmerchantAvailableForSplit_LockedStatusReturnsFalse(t *testing.T) {
	client := NewClient(
		WithClient(
			&http.Client{
				Transport: splitRoundTripFunc(
					func(_ *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body: io.NopCloser(
								strings.NewReader(`{"status":"SUCCESS","action":"GET_SUBMERCHANT","submerchant_id":"123456789","submerchant_id_status":"LOCKED","hash":"abc123"}`),
							),
						}, nil
					},
				),
			},
		),
	)

	submerchantID := "123456789"
	req := &Request{
		Merchant: &Merchant{
			MerchantKey: "CLIENT_KEY",
			SecretKey:   "CLIENT_PASS",
		},
		PaymentData: &PaymentData{
			SubmerchantID: &submerchantID,
		},
	}

	enabled, err := client.SubmerchantAvailableForSplit(req)
	if err != nil {
		t.Fatalf("SubmerchantAvailableForSplit() error: %v", err)
	}
	if enabled {
		t.Fatalf("expected LOCKED submerchant to be unavailable")
	}
}

func TestSubmerchantAvailableForSplit_FailedStatusReturnsError(t *testing.T) {
	client := NewClient(
		WithClient(
			&http.Client{
				Transport: splitRoundTripFunc(
					func(_ *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type": []string{"application/json"},
							},
							Body: io.NopCloser(strings.NewReader(`{"status":"FAILED"}`)),
						}, nil
					},
				),
			},
		),
	)

	submerchantID := "123456789"
	req := &Request{
		Merchant: &Merchant{
			MerchantKey: "CLIENT_KEY",
			SecretKey:   "CLIENT_PASS",
		},
		PaymentData: &PaymentData{
			SubmerchantID: &submerchantID,
		},
	}

	_, err := client.SubmerchantAvailableForSplit(req)
	if err == nil {
		t.Fatalf("expected error for FAILED status, got nil")
	}
	if !strings.Contains(err.Error(), "status=FAILED") {
		t.Fatalf("expected FAILED status in error, got %q", err.Error())
	}
}

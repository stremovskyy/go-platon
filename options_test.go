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

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClient_WithClient_UsesProvidedHTTPClient(t *testing.T) {
	called := false
	httpClient := &http.Client{
		Timeout: 0,
		Transport: roundTripperFunc(
			func(req *http.Request) (*http.Response, error) {
				called = true

				if req.URL.String() != consts.ApiPostUnqURL {
					t.Fatalf("url mismatch: want %q, got %q", consts.ApiPostUnqURL, req.URL.String())
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"result":"ACCEPTED"}`)),
				}, nil
			},
		),
	}

	cl := NewClient(WithClient(httpClient))

	req := &Request{
		Merchant: &Merchant{
			MerchantKey: "clientKey",
			SecretKey:   "secret123",
			TermsURL:    ref("https://merchant.example/3ds"),
		},
		PaymentData: &PaymentData{
			PaymentID:   ref("order-1"),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "one-click payment",
		},
		PaymentMethod: &PaymentMethod{
			Card: &Card{
				Token: ref("TOKEN123"),
			},
		},
		PersonalData: &PersonalData{
			Email: ref("payer@example.com"),
		},
	}

	resp, err := cl.Payment(req)
	if err != nil {
		t.Fatalf("Payment() error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if !called {
		t.Fatalf("custom HTTP client transport was not called")
	}
}

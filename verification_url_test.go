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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stremovskyy/go-platon/platon"
)

func TestResolveClientServerVerificationURL_UsesLocationHeader(t *testing.T) {
	wantURL := "https://secure.platononline.com/payment/purchase?token=ABC123"

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("method mismatch: want %q, got %q", http.MethodPost, r.Method)
				}
				if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
					t.Fatalf("content-type mismatch: want application/x-www-form-urlencoded, got %q", got)
				}

				w.Header().Set("Location", wantURL)
				w.WriteHeader(http.StatusFound)
			},
		),
	)
	defer server.Close()

	form := &platon.ClientServerVerificationForm{
		Method:   http.MethodPost,
		Endpoint: server.URL,
		Fields: map[string]string{
			"payment": "CC",
			"key":     "client",
			"url":     "https://merchant.example/success",
			"data":    "payload",
			"sign":    "signature",
		},
	}

	urlResult, err := resolveClientServerVerificationURL(form)
	if err != nil {
		t.Fatalf("resolveClientServerVerificationURL() error: %v", err)
	}
	if urlResult.String() != wantURL {
		t.Fatalf("URL mismatch: want %q, got %q", wantURL, urlResult.String())
	}
}

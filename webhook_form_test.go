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
)

const webhookFormPayload = "id=47097-87770-07123&order=47097-87309-6110&status=SALE&card=411111%2A%2A%2A%2A1111&description=test&amount=0.40&currency=UAH&email=&date=2026-02-13+10%3A32%3A57&ip=250.137.176.130&sign=582d658d7d422e76b2639fac131d093e"

func TestParseWebhookForm(t *testing.T) {
	form, err := ParseWebhookForm([]byte(webhookFormPayload))
	if err != nil {
		t.Fatalf("ParseWebhookForm() error: %v", err)
	}

	if form.Order != "47097-87309-6110" {
		t.Fatalf("order mismatch: got %q", form.Order)
	}
	if form.Status != "SALE" {
		t.Fatalf("status mismatch: got %q", form.Status)
	}
	if form.Card != "411111****1111" {
		t.Fatalf("card mismatch: got %q", form.Card)
	}
}

/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
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

package main

import (
	"fmt"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/examples/internal/config"
)

func main() {
	cfg := config.MustLoad()

	payload := "id=47097-87770-07123&order=47097-87309-6110&status=SALE&card=411111%2A%2A%2A%2A1111&description=%D0%9F%D0%BE%D0%BF%D0%BE%D0%B2%D0%BD%D0%B5%D0%BD%D0%BD%D1%8F+%D0%B1%D0%B0%D0%BB%D0%B0%D0%BD%D1%81%D1%83+%D0%B2%D0%BE%D0%B4%D1%96%D1%8F+%28Platon+split+one+receiver%29&amount=0.40&currency=UAH&name=+&phone=&email=&date=2026-02-13+10%3A32%3A57&ip=250.137.176.130&sign=582d658d7d422e76b2639fac131d093e&rc_id=47097-87770-07123&rc_token=fa0500fb3f4869247b4c5532eaf799bc&issuing_bank=JPMORGAN+CHASE+BANK%2C+N.A.&ext4=&cardholder_email=&brand=VISA&terminal="
	payerEmail := "payer@example.com"

	form, err := go_platon.ParseWebhookForm([]byte(payload))
	if err != nil {
		fmt.Println("webhook parse error:", err)
		return
	}

	ok, err := form.VerifySign(cfg.SecretKey, payerEmail)
	if err != nil {
		fmt.Println("signature verify error:", err)
		return
	}

	fmt.Printf("order=%s status=%s amount=%s currency=%s sign_valid=%t recurrent_token=%s\n", form.Order, form.Status, form.Amount, form.Currency, ok, form.RCToken)
}

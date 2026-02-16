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

	// Platon sends callbacks to a single URL. Use ext fields to route internally.
	payload := "id=47123-08562-28823&order=47123-08266-9485&status=SALE&card=411111%2A%2A%2A%2A1111&description=Simple+verification+example&amount=0.40&currency=UAH&name=+&phone=&email=&date=2026-02-16+08%3A34%3A16&ip=248.245.244.245&sign=b8a167daec9c8510eda2f313f5e893fd&rc_id=47123-08562-28823&rc_token=d62fc9813c21a035d2b65e30e79ba995&issuing_bank=JPMORGAN+CHASE+BANK%2C+N.A.&card_token=a0b520f81ddd1a087ba83506bcb957d472b7abd5383c90e7e0b56aa3fc271583&ext4=wallet-topup&cardholder_email=&brand=VISA&terminal="
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

	target := "default-handler"
	switch form.Ext4 {
	case "wallet-topup":
		target = "wallet-handler"
	case "order-payment":
		target = "orders-handler"
	}

	fmt.Printf(
		"order=%s status=%s amount=%s currency=%s sign_valid=%t recurrent_token=%s ext4=%s route=%s\n",
		form.Order,
		form.Status,
		form.Amount,
		form.Currency,
		ok,
		form.RCToken,
		form.Ext4,
		target,
	)
}

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

package main

import (
	"fmt"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/examples/demo"
	"github.com/stremovskyy/go-platon/examples/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// Platon sends callbacks to a single URL. Use ext fields to route internally.
	payload := "id=47123-08562-28823&order=396bbff2-ce6e-45f8-8559-3e9540cf3808&status=SALE&card=411111%2A%2A%2A%2A1111&description=%D0%9F%D0%BE%D0%BF%D0%BE%D0%B2%D0%BD%D0%B5%D0%BD%D0%BD%D1%8F+%D0%B1%D0%B0%D0%BB%D0%B0%D0%BD%D1%81%D1%83+%D0%B2%D0%BE%D0%B4%D1%96%D1%8F+%28Platon+split+one+receiver%29&amount=1.00&currency=UAH&name=+&phone=%2B380000000000&email=no-reply%40example.com&date=2026-02-16+08%3A34%3A16&ip=127.0.0.1&sign=b8a167daec9c8510eda2f313f5e893fd&rc_id=47123-08562-28823&rc_token=d62fc9813c21a035d2b65e30e79ba995&issuing_bank=JPMORGAN+CHASE+BANK%2C+N.A.&card_token=35f5f6306f9baa5bb9b58803b7edf64d421d890f1b68e0454c3d45724a342694&ext4=payment%3Atest&ext5=%5Boid%3A396bbff2-ce6e-45f8-8559-3e9540cf3808%5D&cardholder_email=&brand=VISA&terminal="
	payerEmail := demo.PayerEmail

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
	case "payment:test":
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

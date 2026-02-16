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
	"github.com/stremovskyy/go-platon/examples/internal/demo"
	"github.com/stremovskyy/go-platon/log"
)

func main() {
	cfg := config.MustLoad()
	var client go_platon.Platon = go_platon.NewDefaultClient()
	client.SetLogLevel(log.LevelDebug)

	merchant := &go_platon.Merchant{
		MerchantKey: cfg.MerchantKey,
		SecretKey:   cfg.SecretKey,
	}

	req := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				Pan: ref(demo.CardNumber),
			},
		},
		PaymentData: &go_platon.PaymentData{
			PlatonTransID: ref("632508054"),
			Amount:        100,
			Metadata: map[string]string{
				"immediately": "Y",
			},
		},
		PersonalData: &go_platon.PersonalData{
			Email: ref(demo.PayerEmail),
		},
	}

	resp, err := client.Refund(req)
	if err != nil {
		fmt.Println("refund error:", err)
		return
	}

	resp.PrettyPrint()
}

func ref(value string) *string {
	return &value
}

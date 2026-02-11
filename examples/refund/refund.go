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
	"os"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/examples/internal/config"
	"github.com/stremovskyy/go-platon/internal/utils"
	"github.com/stremovskyy/go-platon/log"
)

func main() {
	cfg := config.MustLoad()
	client := go_platon.NewDefaultClient()

	merchant := &go_platon.Merchant{
		Name:        cfg.MerchantName,
		MerchantID:  cfg.MerchantID,
		MerchantKey: cfg.MerchantKey,
		SecretKey:   cfg.SecretKey,
	}

	refundRequest := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				// CREDITVOID signature uses first 6 + last 4 of the PAN.
				Pan: utils.Ref(cfg.CardNumber),
			},
		},
		PaymentData: &go_platon.PaymentData{
			PlatonTransID: utils.Ref("632508054"),
			Amount:        100,
			Metadata: map[string]string{
				// Optional: send `immediately=Y` (fast refund).
				"immediately": "Y",
			},
		},
		PersonalData: &go_platon.PersonalData{
			Email: utils.Ref(cfg.PayerEmail),
		},
	}

	client.SetLogLevel(log.LevelDebug)

	refundResponse, err := client.Refund(refundRequest)
	if err != nil {
		fmt.Println(err)

		if refundResponse != nil {
			refundResponse.PrettyPrint()
		}

		os.Exit(1)
	}

	refundResponse.PrettyPrint()
}

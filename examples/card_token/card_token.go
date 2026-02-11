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

	"github.com/google/uuid"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/examples/internal/config"
	"github.com/stremovskyy/go-platon/internal/utils"
	"github.com/stremovskyy/go-platon/log"
)

func main() {
	cfg := config.MustLoad()
	client := go_platon.NewDefaultClient()

	merchant := &go_platon.Merchant{
		Name:            cfg.MerchantName,
		MerchantID:      cfg.MerchantID,
		MerchantKey:     cfg.MerchantKey,
		SecretKey:       cfg.SecretKey,
		SuccessRedirect: cfg.SuccessRedirect,
		FailRedirect:    cfg.FailRedirect,
		TermsURL:        utils.Ref("https://google.com"),
	}

	orderID := uuid.New().String()

	req := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				Token: utils.Ref(cfg.CardToken),
			},
		},
		PaymentData: &go_platon.PaymentData{
			PaymentID:   utils.Ref(orderID),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "One-click token payment: " + orderID,
		},
		PersonalData: &go_platon.PersonalData{
			Email: utils.Ref(cfg.PayerEmail),
		},
	}

	client.SetLogLevel(log.LevelDebug)

	resp, err := client.Payment(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp.PrettyPrint()
}

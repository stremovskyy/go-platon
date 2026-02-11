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

	client.SetLogLevel(log.LevelDebug)

	submerchantID := "12345678"

	enabled, err := client.SubmerchantAvailableForSplit(&go_platon.Request{
		Merchant: merchant,
		PaymentData: &go_platon.PaymentData{
			SubmerchantID: &submerchantID,
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !enabled {
		fmt.Printf("submerchant %s is not enabled for split\n", submerchantID)
		return
	}

	orderID := uuid.New().String()
	paymentRequest := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				Pan:             utils.Ref(cfg.CardNumber),
				ExpirationMonth: utils.Ref(cfg.CardMonth),
				ExpirationYear:  utils.Ref(cfg.CardYear),
				Cvv2:            utils.Ref(cfg.CardCVV),
			},
		},
		PaymentData: &go_platon.PaymentData{
			PaymentID:   &orderID,
			Amount:      300,
			Currency:    currency.UAH,
			Description: "Split payment: " + orderID,
			SplitRules: []go_platon.SplitRule{
				{SubmerchantIdentification: "12345678", Amount: 100},
				{SubmerchantIdentification: "87654321", Amount: 200},
			},
		},
		PersonalData: &go_platon.PersonalData{
			Email: utils.Ref(cfg.PayerEmail),
			Phone: utils.Ref("380631234567"),
		},
	}

	paymentResponse, err := client.Payment(paymentRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paymentResponse.PrettyPrint()
}

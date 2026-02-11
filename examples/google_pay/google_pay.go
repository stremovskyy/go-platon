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
		Login:           cfg.Login,
		MerchantKey:     cfg.MerchantKey,
		SecretKey:       cfg.SecretKey,
		SuccessRedirect: cfg.SuccessRedirect,
		FailRedirect:    cfg.FailRedirect,
		TermsURL:        utils.Ref("https://google.com"),
	}

	uuidString := uuid.New().String()

	holdRequest := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			GoogleToken: utils.Ref(cfg.GoogleToken),
		},
		PaymentData: &go_platon.PaymentData{
			PaymentID:   utils.Ref(uuidString),
			Amount:      100,
			Currency:    currency.UAH,
			Description: "Test payment: " + uuidString,
			IsMobile:    true,
		},
		PersonalData: &go_platon.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
			Email:     utils.Ref(cfg.PayerEmail),
			Phone:     utils.Ref("380631234567"),
		},
	}

	client.SetLogLevel(log.LevelDebug)

	holdResponse, err := client.Hold(holdRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Payment: %s result=%v error=%v\n", uuidString, holdResponse.Result, holdResponse.ErrorMessage)
}

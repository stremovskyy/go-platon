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

	"github.com/google/uuid"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/examples/demo"
	"github.com/stremovskyy/go-platon/examples/internal/config"
	"github.com/stremovskyy/go-platon/log"
)

func main() {
	cfg := config.MustLoad()
	var client = go_platon.NewDefaultClient()
	client.SetLogLevel(log.LevelDebug)

	merchant := &go_platon.Merchant{
		MerchantID:      cfg.MerchantID,
		MerchantKey:     demo.ClientKey,
		SecretKey:       cfg.SecretKey,
		SuccessRedirect: cfg.SuccessRedirect,
		FailRedirect:    cfg.FailRedirect,
		ClientIP:        ref(demo.PayerIP),
		TermsURL:        ref(demo.TermsURL3DS),
	}

	req := &go_platon.Request{
		Merchant: merchant,
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				Token: ref(demo.CardToken),
			},
		},
		PaymentData: &go_platon.PaymentData{
			PaymentID:   ref(uuid.New().String()),
			Amount:      demo.AmountMinor,
			Currency:    currency.UAH,
			Description: demo.Description,
			Metadata: map[string]string{
				"ext4": demo.Ext4,
				"ext5": demo.Ext5,
			},
		},
		PersonalData: &go_platon.PersonalData{
			Email: ref(demo.PayerEmail),
			Phone: ref(demo.PayerPhone),
		},
	}

	resp, err := client.Hold(req)
	if err != nil {
		fmt.Println("hold error:", err)
		return
	}

	resp.PrettyPrint()
}

func ref(value string) *string {
	return &value
}

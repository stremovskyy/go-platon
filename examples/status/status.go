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
	"os"
	"strings"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/examples/internal/config"
	"github.com/stremovskyy/go-platon/log"
)

func main() {
	cfg := config.MustLoad()
	var client = go_platon.NewDefaultClient()
	client.SetLogLevel(log.LevelDebug)

	merchant := &go_platon.Merchant{
		MerchantKey: cfg.MerchantKey,
		SecretKey:   cfg.SecretKey,
	}

	req, err := buildStatusRequest(merchant)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	resp, err := client.Status(req)
	if err != nil {
		fmt.Println("status error:", err)
		return
	}

	resp.PrettyPrint()
}

func ref(value string) *string {
	return &value
}

func buildStatusRequest(merchant *go_platon.Merchant) (*go_platon.Request, error) {
	transID := env("PLATON_STATUS_TRANS_ID")
	orderID := env("PLATON_STATUS_ORDER_ID")
	statusFlow := strings.ToLower(env("PLATON_STATUS_FLOW"))

	switch {
	case transID != "":
		req := &go_platon.Request{
			Merchant: merchant,
			PaymentData: &go_platon.PaymentData{
				PlatonTransID: ref(transID),
			},
		}

		if email := env("PLATON_STATUS_PAYER_EMAIL"); email != "" {
			req.PersonalData = &go_platon.PersonalData{
				Email: ref(email),
			}
		}

		return req, nil

	case orderID != "":
		req := &go_platon.Request{
			Merchant: merchant,
			PaymentData: &go_platon.PaymentData{
				PaymentID: ref(orderID),
			},
		}
		if statusFlow == "a2c" {
			req.PaymentData.Metadata = map[string]string{
				"platon_flow": "a2c",
			}
		}
		return req, nil

	default:
		return nil, fmt.Errorf("set PLATON_STATUS_ORDER_ID for GET_TRANS_STATUS_BY_ORDER, or PLATON_STATUS_TRANS_ID for GET_TRANS_STATUS")
	}
}

func env(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

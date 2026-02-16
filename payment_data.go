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

package go_platon

import "github.com/stremovskyy/go-platon/currency"

// PaymentData represents the data related to a payment transaction.
type PaymentData struct {
	// PlatonPaymentID is the unique identifier for the Platon payment.
	//
	// Deprecated: Platon trans_id can be non-numeric (e.g. contain hyphens). Prefer PlatonTransID.
	PlatonPaymentID *int64
	// PlatonTransID is the Platon transaction identifier (trans_id) used for GET_TRANS_STATUS/CAPTURE/CREDITVOID.
	PlatonTransID *string
	// PaymentID is the unique identifier for the payment.
	PaymentID *string
	// Amount is the amount of the payment in the smallest unit of the currency.
	Amount int
	// Currency is the currency code of the payment.
	Currency currency.Code
	// Description is a brief description of the payment.
	Description string
	// IsMobile indicates whether the payment was made from a mobile device.
	IsMobile bool
	// SplitRules defines optional split payouts to sub-merchants.
	// Amount is specified in minor units.
	SplitRules []SplitRule
	// SubmerchantID is used by GET_SUBMERCHANT request.
	SubmerchantID *string
	// RelatedIds is a list of related payment IDs.
	RelatedIds []int64
	// Metadata is a map of additional data.
	// Supported integration keys:
	// - ext1..ext10: passed to Platon request fields with the same names.
	// - immediately: for Refund, "Y"/"true"/"1" enables fast refund mode.
	// - platon_flow: for Status, value "a2c" switches to A2C status endpoint.
	Metadata map[string]string
}

// SplitRule defines amount distribution to a specific sub-merchant.
type SplitRule struct {
	SubmerchantIdentification string
	Amount                    int
}

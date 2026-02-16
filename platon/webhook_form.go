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

package platon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// WebhookForm represents Platon callback payload sent as
// application/x-www-form-urlencoded.
type WebhookForm struct {
	ID              string
	Order           string
	Status          string
	Card            string
	Description     string
	Amount          string
	Currency        string
	Name            string
	Phone           string
	Email           string
	Date            string
	IP              string
	Sign            string
	RCID            string
	RCToken         string
	IssuingBank     string
	Ext1            string
	Ext2            string
	Ext3            string
	Ext4            string
	Ext5            string
	Ext6            string
	Ext7            string
	Ext8            string
	Ext9            string
	Ext10           string
	CardholderEmail string
	Brand           string
	Terminal        string
}

// ParseWebhookForm parses Platon callback payload sent as
// application/x-www-form-urlencoded.
func ParseWebhookForm(data []byte) (*WebhookForm, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("webhook form payload is empty")
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("cannot parse webhook form payload: %w", err)
	}

	return ParseWebhookValues(values), nil
}

// ParseWebhookValues maps decoded form fields into a WebhookForm model.
func ParseWebhookValues(values url.Values) *WebhookForm {
	if values == nil {
		return &WebhookForm{}
	}

	return &WebhookForm{
		ID:              strings.TrimSpace(values.Get("id")),
		Order:           strings.TrimSpace(values.Get("order")),
		Status:          strings.TrimSpace(values.Get("status")),
		Card:            strings.TrimSpace(values.Get("card")),
		Description:     strings.TrimSpace(values.Get("description")),
		Amount:          strings.TrimSpace(values.Get("amount")),
		Currency:        strings.TrimSpace(values.Get("currency")),
		Name:            values.Get("name"),
		Phone:           strings.TrimSpace(values.Get("phone")),
		Email:           strings.TrimSpace(values.Get("email")),
		Date:            strings.TrimSpace(values.Get("date")),
		IP:              strings.TrimSpace(values.Get("ip")),
		Sign:            strings.TrimSpace(values.Get("sign")),
		RCID:            strings.TrimSpace(values.Get("rc_id")),
		RCToken:         strings.TrimSpace(values.Get("rc_token")),
		IssuingBank:     strings.TrimSpace(values.Get("issuing_bank")),
		Ext1:            strings.TrimSpace(values.Get("ext1")),
		Ext2:            strings.TrimSpace(values.Get("ext2")),
		Ext3:            strings.TrimSpace(values.Get("ext3")),
		Ext4:            strings.TrimSpace(values.Get("ext4")),
		Ext5:            strings.TrimSpace(values.Get("ext5")),
		Ext6:            strings.TrimSpace(values.Get("ext6")),
		Ext7:            strings.TrimSpace(values.Get("ext7")),
		Ext8:            strings.TrimSpace(values.Get("ext8")),
		Ext9:            strings.TrimSpace(values.Get("ext9")),
		Ext10:           strings.TrimSpace(values.Get("ext10")),
		CardholderEmail: strings.TrimSpace(values.Get("cardholder_email")),
		Brand:           strings.TrimSpace(values.Get("brand")),
		Terminal:        strings.TrimSpace(values.Get("terminal")),
	}
}

// ExpectedSign computes the callback signature based on Platon docs:
// md5(strtoupper(strrev(email)+pass+order+strrev(first6+last4)+strrev(status))).
//
// Email from callback may be empty. In that case, pass the email from your
// original payment request via payerEmailOverride.
func (f *WebhookForm) ExpectedSign(secret string, payerEmailOverride string) (string, error) {
	if f == nil {
		return "", fmt.Errorf("webhook form is nil")
	}

	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", fmt.Errorf("secret is required")
	}
	order := strings.TrimSpace(f.Order)
	if order == "" {
		return "", fmt.Errorf("order is required")
	}
	status := strings.TrimSpace(f.Status)
	if status == "" {
		return "", fmt.Errorf("status is required")
	}
	if f.Card == "" {
		return "", fmt.Errorf("card is required")
	}

	card, err := webhookCardSignSource(f.Card)
	if err != nil {
		return "", err
	}

	payerEmail := strings.TrimSpace(payerEmailOverride)
	if payerEmail == "" {
		payerEmail = f.Email
	}

	raw := reverseString(payerEmail) +
		secret +
		order +
		reverseString(card) +
		reverseString(status)

	hash := md5.Sum([]byte(strings.ToUpper(raw)))
	return hex.EncodeToString(hash[:]), nil
}

// VerifySign validates callback signature against callback `sign` field.
func (f *WebhookForm) VerifySign(secret string, payerEmailOverride string) (bool, error) {
	if f == nil {
		return false, fmt.Errorf("webhook form is nil")
	}
	if f.Sign == "" {
		return false, fmt.Errorf("sign is required")
	}

	expected, err := f.ExpectedSign(secret, payerEmailOverride)
	if err != nil {
		return false, err
	}

	return strings.EqualFold(f.Sign, expected), nil
}

func webhookCardSignSource(card string) (string, error) {
	normalized := strings.ReplaceAll(strings.TrimSpace(card), " ", "")
	if len(normalized) < 10 {
		return "", fmt.Errorf("card value is too short to build signature")
	}

	return normalized[:6] + normalized[len(normalized)-4:], nil
}

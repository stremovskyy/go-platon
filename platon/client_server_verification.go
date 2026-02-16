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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	clientServerVerificationPaymentCode = "CC"
	clientServerVerificationFormID      = "verify"
	clientServerVerificationReqToken    = "Y"
	clientServerVerificationRecurring   = "Y"
	clientServerVerificationMethod      = "POST"
)

// ClientServerVerificationForm contains endpoint and form fields for browser-side
// verification submission (Client-Server flow).
type ClientServerVerificationForm struct {
	Method   string
	Endpoint string
	Fields   map[string]string
}

// ClientServerVerificationParams holds normalized values required to build a
// signed Client-Server verification form.
type ClientServerVerificationParams struct {
	ClientKey   string
	Secret      string
	RedirectURL string
	Description string
	Currency    string
	OrderID     *string
	Metadata    map[string]string
}

type clientServerVerificationData struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Currency    string `json:"currency"`
	Recurring   string `json:"recurring"`
	Order       string `json:"order,omitempty"`
	Ext1        string `json:"ext1,omitempty"`
	Ext2        string `json:"ext2,omitempty"`
	Ext3        string `json:"ext3,omitempty"`
	Ext4        string `json:"ext4,omitempty"`
	Ext5        string `json:"ext5,omitempty"`
	Ext6        string `json:"ext6,omitempty"`
	Ext7        string `json:"ext7,omitempty"`
	Ext8        string `json:"ext8,omitempty"`
	Ext9        string `json:"ext9,omitempty"`
	Ext10       string `json:"ext10,omitempty"`
}

// BuildClientServerVerificationForm builds a signed form payload for
// Client-Server card verification.
func BuildClientServerVerificationForm(params ClientServerVerificationParams, endpoint string) (*ClientServerVerificationForm, error) {
	clientKey := strings.TrimSpace(params.ClientKey)
	if clientKey == "" {
		return nil, fmt.Errorf("verification: merchant client_key is required")
	}

	secret := strings.TrimSpace(params.Secret)
	if secret == "" {
		return nil, fmt.Errorf("verification: merchant secret key is required")
	}

	redirectURL := strings.TrimSpace(params.RedirectURL)
	if redirectURL == "" {
		return nil, fmt.Errorf("verification: success redirect URL is required")
	}

	description := strings.TrimSpace(params.Description)
	if description == "" {
		return nil, fmt.Errorf("verification: order_description is required")
	}

	orderCurrency := strings.TrimSpace(params.Currency)
	if orderCurrency == "" {
		return nil, fmt.Errorf("verification: order_currency is required")
	}

	apiEndpoint := strings.TrimSpace(endpoint)
	if apiEndpoint == "" {
		return nil, fmt.Errorf("verification: endpoint is required")
	}

	data := clientServerVerificationData{
		Amount:      VerifyNoAmount.String(),
		Description: description,
		Currency:    orderCurrency,
		Recurring:   clientServerVerificationRecurring,
	}
	if params.OrderID != nil && strings.TrimSpace(*params.OrderID) != "" {
		data.Order = strings.TrimSpace(*params.OrderID)
	}

	data.Ext1 = metadataValue(params.Metadata, "ext1")
	data.Ext2 = metadataValue(params.Metadata, "ext2")
	data.Ext3 = metadataValue(params.Metadata, "ext3")
	data.Ext4 = metadataValue(params.Metadata, "ext4")
	data.Ext5 = metadataValue(params.Metadata, "ext5")
	data.Ext6 = metadataValue(params.Metadata, "ext6")
	data.Ext7 = metadataValue(params.Metadata, "ext7")
	data.Ext8 = metadataValue(params.Metadata, "ext8")
	data.Ext9 = metadataValue(params.Metadata, "ext9")
	data.Ext10 = metadataValue(params.Metadata, "ext10")

	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("verification: cannot encode data payload: %w", err)
	}
	encodedData := base64.StdEncoding.EncodeToString(rawData)

	sign := signClientServerVerification(clientKey, clientServerVerificationPaymentCode, encodedData, redirectURL, secret)

	form := &ClientServerVerificationForm{
		Method:   clientServerVerificationMethod,
		Endpoint: apiEndpoint,
		Fields: map[string]string{
			"payment":   clientServerVerificationPaymentCode,
			"key":       clientKey,
			"url":       redirectURL,
			"data":      encodedData,
			"formid":    clientServerVerificationFormID,
			"req_token": clientServerVerificationReqToken,
			"sign":      sign,
		},
	}

	// Some Platon installations propagate callback ext fields only when they are
	// sent as top-level form fields (not only inside JSON "data").
	setNonEmptyFormField(form.Fields, "ext1", data.Ext1)
	setNonEmptyFormField(form.Fields, "ext2", data.Ext2)
	setNonEmptyFormField(form.Fields, "ext3", data.Ext3)
	setNonEmptyFormField(form.Fields, "ext4", data.Ext4)
	setNonEmptyFormField(form.Fields, "ext5", data.Ext5)
	setNonEmptyFormField(form.Fields, "ext6", data.Ext6)
	setNonEmptyFormField(form.Fields, "ext7", data.Ext7)
	setNonEmptyFormField(form.Fields, "ext8", data.Ext8)
	setNonEmptyFormField(form.Fields, "ext9", data.Ext9)
	setNonEmptyFormField(form.Fields, "ext10", data.Ext10)

	return form, nil
}

func setNonEmptyFormField(fields map[string]string, key string, value string) {
	if fields == nil {
		return
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}

	fields[key] = trimmed
}

func signClientServerVerification(clientKey string, payment string, data string, redirectURL string, secret string) string {
	raw := reverseString(clientKey) +
		reverseString(payment) +
		reverseString(data) +
		reverseString(redirectURL) +
		reverseString(secret)

	hash := md5.Sum([]byte(strings.ToUpper(raw)))
	return hex.EncodeToString(hash[:])
}

func metadataValue(metadata map[string]string, key string) string {
	if metadata == nil {
		return ""
	}

	return strings.TrimSpace(metadata[key])
}

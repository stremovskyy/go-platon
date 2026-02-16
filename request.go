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

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/platon"
)

type Request struct {
	Merchant      *Merchant
	PersonalData  *PersonalData
	PaymentData   *PaymentData
	PaymentMethod *PaymentMethod
}

// BuildClientServerVerificationForm builds signed browser form fields for
// Client-Server card verification (`/payment/auth`).
func BuildClientServerVerificationForm(request *Request) (*platon.ClientServerVerificationForm, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}
	if request.Merchant == nil {
		return nil, fmt.Errorf("verification: merchant is required for client-server flow")
	}

	redirectURL := strings.TrimSpace(request.GetSuccessRedirect())
	if redirectURL == "" {
		redirectURL = strings.TrimSpace(request.GetFailRedirect())
	}

	return platon.BuildClientServerVerificationForm(
		platon.ClientServerVerificationParams{
			ClientKey:   request.GetMerchantKey(),
			Secret:      request.Merchant.SecretKey,
			RedirectURL: redirectURL,
			Description: request.GetDescription(),
			Currency:    request.GetCurrency().String(),
			OrderID:     request.GetPaymentID(),
			Metadata:    request.GetMetadata(),
		},
		consts.ApiPaymentAuthURL,
	)
}

func (r *Request) GetAuth() *platon.Auth {
	if r == nil {
		return &platon.Auth{
			Key:    "EMPTY_KEY",
			Secret: "EMPTY_SECRET",
		}
	}

	if r.Merchant == nil {
		return &platon.Auth{
			Key:    "EMPTY_KEY",
			Secret: "EMPTY_SECRET",
		}
	}

	return &platon.Auth{
		Key:    r.Merchant.MerchantKey,
		Secret: r.Merchant.SecretKey,
	}
}
func (r *Request) GetSuccessRedirect() string {
	if r == nil {
		return ""
	}

	if r.Merchant == nil {
		return ""
	}
	return r.Merchant.SuccessRedirect
}

func (r *Request) GetFailRedirect() string {
	if r == nil {
		return ""
	}

	if r.Merchant == nil {
		return ""
	}
	return r.Merchant.FailRedirect
}

func (r *Request) GetPlatonPaymentID() int64 {
	if r == nil {
		return 0
	}

	if r.PaymentData == nil || r.PaymentData.PlatonPaymentID == nil {
		return 0
	}

	return *r.PaymentData.PlatonPaymentID
}

// GetPlatonTransID returns Platon `trans_id` as string.
// It prefers PaymentData.PlatonTransID and falls back to formatting PaymentData.PlatonPaymentID.
func (r *Request) GetPlatonTransID() *string {
	if r == nil {
		return nil
	}

	if r.PaymentData == nil {
		return nil
	}
	if r.PaymentData.PlatonTransID != nil && *r.PaymentData.PlatonTransID != "" {
		return r.PaymentData.PlatonTransID
	}
	if r.PaymentData.PlatonPaymentID != nil {
		s := fmt.Sprintf("%d", *r.PaymentData.PlatonPaymentID)
		return &s
	}
	return nil
}

func (r *Request) GetCardToken() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Token
}

func (r *Request) GetCardPan() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Pan
}

func (r *Request) GetPaymentID() *string {
	if r == nil {
		return nil
	}

	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.PaymentID
}

func (r *Request) GetPayerEmail() *string {
	if r == nil {
		return nil
	}

	if r.PersonalData == nil {
		return nil
	}

	return r.PersonalData.Email
}

func (r *Request) GetPayerPhone() *string {
	if r == nil {
		return nil
	}

	if r.PersonalData == nil {
		return nil
	}

	return r.PersonalData.Phone
}

func (r *Request) SetRedirects(successURL string, failURL string) {
	if r == nil {
		return
	}

	if r.Merchant == nil {
		r.Merchant = &Merchant{}
	}

	r.Merchant.SuccessRedirect = successURL
	r.Merchant.FailRedirect = failURL
}

func (r *Request) GetAmount() float32 {
	if r == nil {
		return 0
	}

	if r.PaymentData == nil {
		return 0
	}

	return float32(r.PaymentData.Amount) / 100
}

func (r *Request) GetDescription() string {
	if r == nil {
		return ""
	}

	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Description
}

func (r *Request) GetCurrency() currency.Code {
	if r == nil {
		return ""
	}

	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Currency

}

func (r *Request) IsMobile() bool {
	if r == nil {
		return false
	}

	if r.PaymentData == nil {
		return false
	}

	if r.PaymentData.IsMobile {
		return true
	}
	if r.PaymentMethod == nil {
		return false
	}
	if r.PaymentMethod.AppleContainer != nil && *r.PaymentMethod.AppleContainer != "" {
		return true
	}
	if r.PaymentMethod.GoogleToken != nil && *r.PaymentMethod.GoogleToken != "" {
		return true
	}
	return false
}

func (r *Request) GetAppleContainer() (*string, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if r.PaymentMethod == nil || r.PaymentMethod.AppleContainer == nil {
		return nil, fmt.Errorf("Apple Container is not set")
	}
	if *r.PaymentMethod.AppleContainer == "" {
		return nil, fmt.Errorf("Apple Container is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(*r.PaymentMethod.AppleContainer)
	if err != nil {
		return nil, fmt.Errorf("cannot decode Apple Container: %w", err)
	}

	var token map[string]interface{}
	if errr := json.Unmarshal(decoded, &token); errr != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", errr)
	}

	outputJSON, err := json.Marshal(token["token"])
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	outputBase64 := base64.StdEncoding.EncodeToString(outputJSON)
	return &outputBase64, nil
}

func (r *Request) IsApplePay() bool {
	if r == nil {
		return false
	}

	return r.PaymentMethod != nil && r.PaymentMethod.AppleContainer != nil && *r.PaymentMethod.AppleContainer != ""
}

func (r *Request) GetGoogleToken() (*string, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if r.PaymentMethod == nil || r.PaymentMethod.GoogleToken == nil {
		return nil, fmt.Errorf("Google Token is not set")
	}
	if *r.PaymentMethod.GoogleToken == "" {
		return nil, fmt.Errorf("Google Token is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(*r.PaymentMethod.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("cannot decode Google Token: %w", err)
	}

	var data struct {
		PaymentMethodData struct {
			TokenizationData struct {
				Token string `json:"token"`
			} `json:"tokenizationData"`
		} `json:"paymentMethodData"`
	}

	if errr := json.Unmarshal(decoded, &data); errr != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", errr)
	}

	unescapedToken, err := strconv.Unquote(fmt.Sprintf("%q", data.PaymentMethodData.TokenizationData.Token))
	if err != nil {
		return nil, fmt.Errorf("unquote error: %w", err)
	}

	outputBase64 := base64.StdEncoding.EncodeToString([]byte(unescapedToken))
	return &outputBase64, nil
}

func (r *Request) GetTrackingData() *int64 {
	if r == nil {
		return nil
	}

	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.PlatonPaymentID
}

func (r *Request) GetSplitRules() (platon.SplitRules, error) {
	if r == nil {
		return nil, nil
	}

	if r.PaymentData == nil || len(r.PaymentData.SplitRules) == 0 {
		return nil, nil
	}
	if r.PaymentData.Amount <= 0 {
		return nil, fmt.Errorf("amount (minor units) must be > 0 when split rules are provided")
	}

	result := make(platon.SplitRules, len(r.PaymentData.SplitRules))
	totalMinorUnits := 0

	for idx, rule := range r.PaymentData.SplitRules {
		identification := strings.TrimSpace(rule.SubmerchantIdentification)
		if identification == "" {
			return nil, fmt.Errorf("split_rules[%d]: submerchant identification is required", idx)
		}
		if rule.Amount <= 0 {
			return nil, fmt.Errorf("split_rules[%d]: amount (minor units) must be > 0", idx)
		}

		totalMinorUnits += rule.Amount
		if totalMinorUnits > r.PaymentData.Amount {
			return nil, fmt.Errorf("split rules total exceeds amount (%d > %d minor units)", totalMinorUnits, r.PaymentData.Amount)
		}

		if _, exists := result[identification]; exists {
			return nil, fmt.Errorf("split_rules[%d]: duplicate submerchant identification %q", idx, identification)
		}

		result[identification] = fmt.Sprintf("%.2f", float64(rule.Amount)/100)
	}

	if totalMinorUnits != r.PaymentData.Amount {
		return nil, fmt.Errorf("split rules total must equal amount (%d != %d minor units)", totalMinorUnits, r.PaymentData.Amount)
	}

	return result, nil
}

func (r *Request) GetSubmerchantID() *string {
	if r == nil {
		return nil
	}

	if r.PaymentData == nil || r.PaymentData.SubmerchantID == nil {
		return nil
	}

	id := strings.TrimSpace(*r.PaymentData.SubmerchantID)
	if id == "" {
		return nil
	}

	return &id
}

func (r *Request) GetReceiverTIN() *string {
	if r == nil {
		return nil
	}

	if r.PersonalData == nil {
		return nil
	}

	return r.PersonalData.TaxID
}

func (r *Request) GetRelatedIDs() []int64 {
	if r == nil {
		return nil
	}

	if r.PaymentData == nil || r.PaymentData.RelatedIds == nil {
		return nil
	}

	return r.PaymentData.RelatedIds
}

func (r *Request) GetMetadata() map[string]string {
	if r == nil {
		return map[string]string{}
	}

	if r.PaymentData == nil {
		return map[string]string{}
	}

	return r.PaymentData.Metadata
}

func (r *Request) GetMerchantKey() string {
	if r == nil {
		return ""
	}

	if r.Merchant == nil {
		return ""
	}

	return r.Merchant.MerchantKey
}

func (r *Request) GetClientIP() *string {
	if r == nil {
		return nil
	}

	if r.Merchant == nil {
		return nil
	}

	return r.Merchant.ClientIP
}

func (r *Request) GetTermsURL() *string {
	if r == nil {
		return nil
	}

	if r.Merchant == nil {
		return nil
	}

	return r.Merchant.TermsURL
}

func (r *Request) GetCardNumber() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Pan
}

func (r *Request) GetCardExpMonth() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil || r.PaymentMethod.Card.ExpirationMonth == nil {
		return nil
	}

	return r.PaymentMethod.Card.ExpirationMonth
}

func (r *Request) GetCardExpYear() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil || r.PaymentMethod.Card.ExpirationYear == nil {
		return nil
	}

	return r.PaymentMethod.Card.ExpirationYear
}

func (r *Request) GetCardCvv2() *string {
	if r == nil {
		return nil
	}

	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil || r.PaymentMethod.Card.Cvv2 == nil {
		return nil
	}

	return r.PaymentMethod.Card.Cvv2
}

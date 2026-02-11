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

package go_platon

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/stremovskyy/go-platon/consts"
	"github.com/stremovskyy/go-platon/internal/http"
	"github.com/stremovskyy/go-platon/log"
	"github.com/stremovskyy/go-platon/platon"
	"github.com/stremovskyy/recorder"
)

type client struct {
	platonClient *http.Client
}

func (c *client) SetLogLevel(levelDebug log.Level) {
	log.SetLevel(levelDebug)
}

func NewDefaultClient() Platon {
	return &client{
		platonClient: http.NewClient(http.DefaultOptions()),
	}
}

func NewClientWithRecorder(rec recorder.Recorder) Platon {
	return &client{
		platonClient: http.NewClient(http.DefaultOptions()).WithRecorder(rec),
	}
}

func (c *client) Verification(request *Request) (*platon.Result, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}
	createTokenRequest := platon.NewRequest(platon.ActionCodeSALE).WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithChannelNoAmountVerification().
		WithOrderID(request.GetPaymentID()).
		WithOrderAmount(platon.VerifyNoAmount.String()).
		ForCurrency(request.GetCurrency()).
		WithDescription(request.GetDescription()).
		WithPayerIP(request.GetClientIP()).
		WithTermsURL(request.GetTermsURL()).
		WithCardNumber(request.GetCardNumber()).
		WithCardExpMonth(request.GetCardExpMonth()).
		WithCardExpYear(request.GetCardExpYear()).
		WithCardCvv2(request.GetCardCvv2()).
		WithPayerEmail(request.GetPayerEmail()).
		WithPayerPhone(request.GetPayerPhone()).
		WithReqToken(true).
		WithRecurringInitFlag(true).
		UseAsync().
		SignForAction(platon.HashTypeVerification)

	response, err := c.platonClient.Api(createTokenRequest, consts.ApiVerifyURL)
	if err != nil {
		return nil, fmt.Errorf("verification API call: %w", err)
	}

	return response.Result, nil
}

func (c *client) Status(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	transID := request.GetPlatonTransID()
	if transID == nil || *transID == "" {
		return nil, fmt.Errorf("status: trans_id is required (set PaymentData.PlatonTransID or PaymentData.PlatonPaymentID)")
	}

	statusRequest := platon.NewRequest(platon.ActionCodeGetTransStatus).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithTransID(transID).
		WithHashEmail(request.GetPayerEmail()).
		SignForAction(platon.HashTypeGetTransStatus)

	pan := request.GetCardNumber()
	if pan == nil || *pan == "" {
		return nil, fmt.Errorf("status: card_number is required to build signature (only first 6 and last 4 are used)")
	}
	cardHashPart, err := cardHashPartFromPAN(*pan)
	if err != nil {
		return nil, fmt.Errorf("status: cannot derive card hash part from card_number: %w", err)
	}
	statusRequest.WithCardHashPart(&cardHashPart)

	return c.platonClient.Api(statusRequest, consts.ApiGetTransStatus)
}

func (c *client) SubmerchantAvailableForSplit(request *Request) (bool, error) {
	if request == nil {
		return false, platon.ErrRequestIsNil
	}
	if request.GetMerchantKey() == "" {
		return false, fmt.Errorf("split availability: merchant client_key is required")
	}
	submerchantID := request.GetSubmerchantID()
	if submerchantID == nil || *submerchantID == "" {
		return false, fmt.Errorf("split availability: submerchant_id is required")
	}

	apiRequest := platon.NewRequest(platon.ActionCodeGetSubmerchant).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithSubmerchantID(submerchantID).
		SignForAction(platon.HashTypeGetSubmerchant)

	response, err := c.platonClient.Api(apiRequest, consts.ApiGetSubmerchant)
	if err != nil {
		return false, fmt.Errorf("split availability API call: %w", err)
	}

	status, ok := response.SubmerchantIDStatus()
	if !ok {
		return false, fmt.Errorf("split availability: response does not contain submerchant_id_status")
	}

	switch strings.ToUpper(status) {
	case "ENABLED":
		return true, nil
	case "DISABLED":
		return false, nil
	default:
		return false, fmt.Errorf("split availability: unknown submerchant_id_status %q", status)
	}
}

func (c *client) Payment(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	apiRequest, apiURL, err := c.buildIAPaymentRequest(request, false)
	if err != nil {
		return nil, err
	}

	response, err := c.platonClient.Api(apiRequest, apiURL)
	if err != nil {
		return nil, fmt.Errorf("payment API call: %w", err)
	}

	return response, nil
}

func (c *client) Hold(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	apiRequest, apiURL, err := c.buildIAPaymentRequest(request, true)
	if err != nil {
		return nil, err
	}

	response, err := c.platonClient.Api(apiRequest, apiURL)
	if err != nil {
		return nil, fmt.Errorf("hold API call: %w", err)
	}

	return response, nil
}

func (c *client) buildIAPaymentRequest(request *Request, hold bool) (*platon.Request, string, error) {
	if request == nil {
		return nil, "", platon.ErrRequestIsNil
	}
	if request.PaymentData == nil {
		return nil, "", fmt.Errorf("payment: PaymentData is nil")
	}
	if request.GetMerchantKey() == "" {
		return nil, "", fmt.Errorf("payment: merchant client_key is required")
	}
	if request.GetPaymentID() == nil || *request.GetPaymentID() == "" {
		return nil, "", fmt.Errorf("payment: order_id (PaymentData.PaymentID) is required")
	}
	if request.GetCurrency() == "" {
		return nil, "", fmt.Errorf("payment: order_currency is required")
	}
	if request.GetDescription() == "" {
		return nil, "", fmt.Errorf("payment: order_description is required")
	}
	splitRules, err := request.GetSplitRules()
	if err != nil {
		return nil, "", fmt.Errorf("payment: invalid split rules: %w", err)
	}

	common := func(action platon.ActionCode) *platon.Request {
		base := platon.NewRequest(action).
			WithAuth(request.GetAuth()).
			WithClientKey(request.GetMerchantKey()).
			WithOrderID(request.GetPaymentID()).
			WithOrderAmountMinorUnits(request.PaymentData.Amount).
			ForCurrency(request.GetCurrency()).
			WithDescription(request.GetDescription()).
			WithPayerIP(request.GetClientIP()).
			WithTermsURL(request.GetTermsURL()).
			WithPayerEmail(request.GetPayerEmail()).
			WithPayerPhone(request.GetPayerPhone()).
			WithPayerFirstName(nil).
			WithPayerLastName(nil)

		if request.PersonalData != nil {
			base.WithPayerFirstName(request.PersonalData.FirstName).
				WithPayerLastName(request.PersonalData.LastName)
		}

		if hold {
			base.WithHoldAuth()
		}

		return base
	}

	// Mobile payments.
	if request.IsApplePay() {
		container, err := request.GetAppleContainer()
		if err != nil {
			return nil, "", fmt.Errorf("payment: cannot get Apple Pay container: %w", err)
		}
		apiRequest := common(platon.ActionCodeAPPLEPAY).
			WithPaymentToken(container).
			WithSplitRules(splitRules).
			SignForAction(platon.HashTypeApplePay)
		return apiRequest, consts.ApiPostURL, nil
	}

	if request.PaymentMethod != nil && request.PaymentMethod.GoogleToken != nil {
		token, err := request.GetGoogleToken()
		if err != nil {
			return nil, "", fmt.Errorf("payment: cannot get Google Pay token: %w", err)
		}
		apiRequest := common(platon.ActionCodeGOOGLEPAY).
			WithPaymentToken(token).
			WithSplitRules(splitRules).
			SignForAction(platon.HashTypeGooglePay)
		return apiRequest, consts.ApiPostURL, nil
	}

	// One-click by CARD_TOKEN.
	if token := request.GetCardToken(); token != nil && *token != "" {
		apiRequest := common(platon.ActionCodeSALE).
			WithCardToken(token).
			WithSplitRules(splitRules).
			SignForAction(platon.HashTypeCardTokenPayment)
		return apiRequest, consts.ApiPostUnqURL, nil
	}

	// Classic PAN payment.
	if pan := request.GetCardNumber(); pan != nil && *pan != "" {
		apiRequest := common(platon.ActionCodeSALE).
			WithReqToken(false).
			WithRecurringInitFlag(false).
			WithCardNumber(pan).
			WithCardExpMonth(request.GetCardExpMonth()).
			WithCardExpYear(request.GetCardExpYear()).
			WithCardCvv2(request.GetCardCvv2()).
			WithSplitRules(splitRules).
			SignForAction(platon.HashTypeCardPayment)

		return apiRequest, consts.ApiPostUnqURL, nil
	}

	return nil, "", fmt.Errorf("payment: unsupported payment method (expected card PAN, CARD_TOKEN, Apple Pay, or Google Pay data)")
}

func (c *client) Capture(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("capture: %w", platon.ErrRequestIsNil)
	}

	transID := request.GetPlatonTransID()
	if transID == nil || *transID == "" {
		return nil, fmt.Errorf("capture: trans_id is required (set PaymentData.PlatonTransID or PaymentData.PlatonPaymentID)")
	}
	if request.GetMerchantKey() == "" {
		return nil, fmt.Errorf("capture: merchant client_key is required")
	}
	if request.PaymentData == nil {
		return nil, fmt.Errorf("capture: PaymentData is nil")
	}
	if request.PaymentData.Amount <= 0 {
		return nil, fmt.Errorf("capture: PaymentData.Amount (minor units) must be > 0")
	}
	splitRules, err := request.GetSplitRules()
	if err != nil {
		return nil, fmt.Errorf("capture: invalid split rules: %w", err)
	}

	pan := request.GetCardNumber()
	if pan == nil || *pan == "" {
		return nil, fmt.Errorf("capture: card_number is required to build signature (only first 6 and last 4 are used)")
	}
	cardHashPart, err := cardHashPartFromPAN(*pan)
	if err != nil {
		return nil, fmt.Errorf("capture: cannot derive card hash part from card_number: %w", err)
	}

	apiRequest := platon.NewRequest(platon.ActionCodeCAPTURE).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithTransID(transID).
		WithAmountMinorUnits(request.PaymentData.Amount).
		WithSplitRules(splitRules).
		WithHashEmail(request.GetPayerEmail()).
		WithCardHashPart(&cardHashPart).
		SignForAction(platon.HashTypeCapture)

	return c.platonClient.Api(apiRequest, consts.ApiPostUnqURL)
}

func (c *client) Refund(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("refund: %w", platon.ErrRequestIsNil)
	}

	transID := request.GetPlatonTransID()
	if transID == nil || *transID == "" {
		return nil, fmt.Errorf("refund: trans_id is required (set PaymentData.PlatonTransID or PaymentData.PlatonPaymentID)")
	}
	if request.GetMerchantKey() == "" {
		return nil, fmt.Errorf("refund: merchant client_key is required")
	}
	if request.PaymentData == nil {
		return nil, fmt.Errorf("refund: PaymentData is nil")
	}
	if request.PaymentData.Amount <= 0 {
		return nil, fmt.Errorf("refund: PaymentData.Amount (minor units) must be > 0")
	}
	splitRules, err := request.GetSplitRules()
	if err != nil {
		return nil, fmt.Errorf("refund: invalid split rules: %w", err)
	}

	pan := request.GetCardNumber()
	if pan == nil || *pan == "" {
		return nil, fmt.Errorf("refund: card_number is required to build signature (only first 6 and last 4 are used)")
	}
	cardHashPart, err := cardHashPartFromPAN(*pan)
	if err != nil {
		return nil, fmt.Errorf("refund: cannot derive card hash part from card_number: %w", err)
	}

	apiRequest := platon.NewRequest(platon.ActionCodeCREDITVOID).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithTransID(transID).
		WithAmountMinorUnits(request.PaymentData.Amount).
		WithSplitRules(splitRules).
		WithHashEmail(request.GetPayerEmail()).
		WithCardHashPart(&cardHashPart)

	// Optional fast refund flag. If user sets PaymentData.Metadata["immediately"] to "Y"/"true"/"1",
	// send `immediately=Y` as per IA docs.
	if request.PaymentData.Metadata != nil {
		switch strings.ToUpper(strings.TrimSpace(request.PaymentData.Metadata["immediately"])) {
		case "Y", "YES", "TRUE", "1":
			apiRequest.WithImmediately(true)
		}
	}

	apiRequest.SignForAction(platon.HashTypeCreditVoid)
	return c.platonClient.Api(apiRequest, consts.ApiPostUnqURL)
}

func (c *client) Credit(request *Request) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("credit: %w", platon.ErrRequestIsNil)
	}

	return nil, fmt.Errorf("credit: %w", platon.ErrNotImplemented)
}

func cardHashPartFromPAN(pan string) (string, error) {
	digits := digitsOnly(pan)
	if len(digits) < 10 {
		return "", fmt.Errorf("card_number must contain at least 10 digits (got %q)", pan)
	}
	return digits[:6] + digits[len(digits)-4:], nil
}

func digitsOnly(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

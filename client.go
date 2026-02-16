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
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/stremovskyy/go-platon/consts"
	internalhttp "github.com/stremovskyy/go-platon/internal/http"
	"github.com/stremovskyy/go-platon/log"
	"github.com/stremovskyy/go-platon/platon"
	"github.com/stremovskyy/recorder"
)

type client struct {
	platonClient *internalhttp.Client
}

var _ Platon = (*client)(nil)

const (
	platonMetaFlow = "platon_flow"
	platonFlowA2C  = "a2c"

	defaultA2CFirstName = "Payer"
	defaultA2CLastName  = "Cardholder"
	defaultA2CAddress   = "N/A"
	defaultA2CCountry   = "UA"
	defaultA2CState     = "UA"
	defaultA2CCity      = "Kyiv"
	defaultA2CZip       = "00000"
)

func (c *client) SetLogLevel(levelDebug log.Level) {
	log.SetLevel(levelDebug)
}

func NewDefaultClient() Platon {
	return NewClient()
}

func NewClientWithRecorder(rec recorder.Recorder) Platon {
	return NewClient(WithRecorder(rec))
}

func (c *client) Verification(request *Request, runOpts ...RunOption) (*url.URL, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	form, err := BuildClientServerVerificationForm(request)
	if err != nil {
		return nil, err
	}

	opts := collectRunOptions(runOpts)
	if opts.isDryRun() {
		opts.handleDryRun(consts.ApiPaymentAuthURL, form)
		return nil, nil
	}

	return resolveClientServerVerificationURL(form)
}

func (c *client) VerificationLink(request *Request, runOpts ...RunOption) (*url.URL, error) {
	return c.Verification(request, runOpts...)
}

func (c *client) Status(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	opts := collectRunOptions(runOpts)

	orderID := request.GetPaymentID()
	if orderID == nil || strings.TrimSpace(*orderID) == "" {
		return nil, fmt.Errorf("status: order_id is required (set PaymentData.PaymentID)")
	}

	isA2C := isA2CStatusRequest(request)
	statusHashType := platon.HashTypeGetTransStatusByOrder
	statusURL := consts.ApiGetTransStatus
	if isA2C {
		statusURL = consts.ApiP2PUnqURL
	}

	statusRequest := platon.NewRequest(platon.ActionCodeGetTransStatusByOrder).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithOrderID(orderID).
		SignForAction(statusHashType)

	if opts.isDryRun() {
		opts.handleDryRun(statusURL, statusRequest)
		return nil, nil
	}

	return c.platonClient.Api(statusRequest, statusURL)
}

func (c *client) SubmerchantAvailableForSplit(request *Request, runOpts ...RunOption) (bool, error) {
	if request == nil {
		return false, platon.ErrRequestIsNil
	}

	opts := collectRunOptions(runOpts)

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

	if opts.isDryRun() {
		opts.handleDryRun(consts.ApiGetSubmerchant, apiRequest)
		return false, nil
	}

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

func (c *client) Payment(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	opts := collectRunOptions(runOpts)

	apiRequest, apiURL, err := c.buildIAPaymentRequest(request, false)
	if err != nil {
		return nil, err
	}

	if opts.isDryRun() {
		opts.handleDryRun(apiURL, apiRequest)
		return nil, nil
	}

	response, err := c.platonClient.Api(apiRequest, apiURL)
	if err != nil {
		return nil, fmt.Errorf("payment API call: %w", err)
	}

	return response, nil
}

func (c *client) Hold(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, platon.ErrRequestIsNil
	}

	opts := collectRunOptions(runOpts)

	apiRequest, apiURL, err := c.buildIAPaymentRequest(request, true)
	if err != nil {
		return nil, err
	}

	if opts.isDryRun() {
		opts.handleDryRun(apiURL, apiRequest)
		return nil, nil
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
			WithPayerPhone(request.GetPayerPhone())

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

	return nil, "", fmt.Errorf("payment: unsupported payment method (expected CARD_TOKEN, Apple Pay, or Google Pay data)")
}

func (c *client) Capture(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("capture: %w", platon.ErrRequestIsNil)
	}

	opts := collectRunOptions(runOpts)

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

	apiRequest := platon.NewRequest(platon.ActionCodeCAPTURE).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithTransID(transID).
		WithAmountMinorUnits(request.PaymentData.Amount).
		WithSplitRules(splitRules).
		WithHashEmail(request.GetPayerEmail()).
		SignForAction(platon.HashTypeCapture)

	if opts.isDryRun() {
		opts.handleDryRun(consts.ApiPostUnqURL, apiRequest)
		return nil, nil
	}

	return c.platonClient.Api(apiRequest, consts.ApiPostUnqURL)
}

func (c *client) Refund(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("refund: %w", platon.ErrRequestIsNil)
	}

	opts := collectRunOptions(runOpts)

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

	apiRequest := platon.NewRequest(platon.ActionCodeCREDITVOID).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithTransID(transID).
		WithAmountMinorUnits(request.PaymentData.Amount).
		WithSplitRules(splitRules).
		WithHashEmail(request.GetPayerEmail())

	// Optional fast refund flag. If user sets PaymentData.Metadata["immediately"] to "Y"/"true"/"1",
	// send `immediately=Y` as per IA docs.
	if request.PaymentData.Metadata != nil {
		switch strings.ToUpper(strings.TrimSpace(request.PaymentData.Metadata["immediately"])) {
		case "Y", "YES", "TRUE", "1":
			apiRequest.WithImmediately(true)
		}
	}

	apiRequest.SignForAction(platon.HashTypeCreditVoid)

	if opts.isDryRun() {
		opts.handleDryRun(consts.ApiPostUnqURL, apiRequest)
		return nil, nil
	}

	return c.platonClient.Api(apiRequest, consts.ApiPostUnqURL)
}

func (c *client) Credit(request *Request, runOpts ...RunOption) (*platon.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("credit: %w", platon.ErrRequestIsNil)
	}

	opts := collectRunOptions(runOpts)
	if request.GetMerchantKey() == "" {
		return nil, fmt.Errorf("credit: merchant client_key is required")
	}
	if request.PaymentData == nil {
		return nil, fmt.Errorf("credit: PaymentData is nil")
	}
	if request.GetPaymentID() == nil || *request.GetPaymentID() == "" {
		return nil, fmt.Errorf("credit: order_id (PaymentData.PaymentID) is required")
	}
	if request.PaymentData.Amount <= 0 {
		return nil, fmt.Errorf("credit: PaymentData.Amount (minor units) must be > 0")
	}
	if request.GetCurrency() == "" {
		return nil, fmt.Errorf("credit: order_currency is required")
	}
	if request.GetDescription() == "" {
		return nil, fmt.Errorf("credit: order_description is required")
	}

	if splitRules, err := request.GetSplitRules(); err != nil {
		return nil, fmt.Errorf("credit: invalid split rules: %w", err)
	} else if len(splitRules) > 0 {
		return nil, fmt.Errorf("credit: split rules are not supported for CREDIT2CARD")
	}

	a2cPayer := resolveA2CPayerData(request)
	apiRequest := platon.NewRequest(platon.ActionCodeCREDIT2CARD).
		WithAuth(request.GetAuth()).
		WithClientKey(request.GetMerchantKey()).
		WithOrderID(request.GetPaymentID()).
		WithAmountMinorUnits(request.PaymentData.Amount).
		ForCurrency(request.GetCurrency()).
		WithDescription(request.GetDescription()).
		WithPayerFirstName(a2cPayer.FirstName).
		WithPayerLastName(a2cPayer.LastName).
		WithPayerAddress(a2cPayer.Address).
		WithPayerCountry(a2cPayer.Country).
		WithPayerState(a2cPayer.State).
		WithPayerCity(a2cPayer.City).
		WithPayerZip(a2cPayer.Zip).
		WithPayerEmail(request.GetPayerEmail()).
		WithPayerPhone(request.GetPayerPhone())

	if token := request.GetCardToken(); token != nil && *token != "" {
		apiRequest.WithCardToken(token).SignForAction(platon.HashTypeCredit2CardToken)
	} else {
		return nil, fmt.Errorf("credit: card_token is required")
	}

	if opts.isDryRun() {
		opts.handleDryRun(consts.ApiP2PUnqURL, apiRequest)
		return nil, nil
	}

	return c.platonClient.Api(apiRequest, consts.ApiP2PUnqURL)
}

// ParseWebhookXML parses legacy XML webhook payload.
//
// Deprecated: Platon production callbacks use application/x-www-form-urlencoded.
// Use go_platon.ParseWebhookForm for callback parsing and signature verification.
func (c *client) ParseWebhookXML(data []byte) (*platon.Payment, error) {
	return platon.ParsePaymentXML(data)
}

func isA2CStatusRequest(request *Request) bool {
	if request == nil {
		return false
	}

	metadata := request.GetMetadata()
	if metadata != nil {
		if flow, ok := metadata[platonMetaFlow]; ok && strings.EqualFold(strings.TrimSpace(flow), platonFlowA2C) {
			return true
		}
	}

	return false
}

type a2cPayerData struct {
	FirstName *string
	LastName  *string
	Address   *string
	Country   *string
	State     *string
	City      *string
	Zip       *string
}

func resolveA2CPayerData(request *Request) a2cPayerData {
	metadata := request.GetMetadata()

	firstName := firstNonEmptyPointer(
		pointerStringFromPersonalData(request, func(data *PersonalData) *string { return data.FirstName }),
		stringPointerFromMetadata(metadata, "payer_first_name"),
		stringRef(defaultA2CFirstName),
	)
	lastName := firstNonEmptyPointer(
		pointerStringFromPersonalData(request, func(data *PersonalData) *string { return data.LastName }),
		stringPointerFromMetadata(metadata, "payer_last_name"),
		stringRef(defaultA2CLastName),
	)
	address := firstNonEmptyPointer(
		stringPointerFromMetadata(metadata, "payer_address"),
		stringRef(defaultA2CAddress),
	)
	country := normalizeTwoLetterValue(firstNonEmptyPointer(
		stringPointerFromMetadata(metadata, "payer_country"),
		stringRef(defaultA2CCountry),
	), defaultA2CCountry)
	state := normalizeTwoLetterValue(firstNonEmptyPointer(
		stringPointerFromMetadata(metadata, "payer_state"),
		stringPointerFromMetadata(metadata, "payer_country"),
		stringRef(defaultA2CState),
	), defaultA2CState)
	city := firstNonEmptyPointer(
		stringPointerFromMetadata(metadata, "payer_city"),
		stringRef(defaultA2CCity),
	)
	zip := firstNonEmptyPointer(
		stringPointerFromMetadata(metadata, "payer_zip"),
		stringRef(defaultA2CZip),
	)

	return a2cPayerData{
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
		Country:   country,
		State:     state,
		City:      city,
		Zip:       zip,
	}
}

func pointerStringFromPersonalData(request *Request, getter func(*PersonalData) *string) *string {
	if request == nil || request.PersonalData == nil || getter == nil {
		return nil
	}

	return getter(request.PersonalData)
}

func stringPointerFromMetadata(metadata map[string]string, key string) *string {
	if metadata == nil {
		return nil
	}
	value, ok := metadata[key]
	if !ok {
		return nil
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func firstNonEmptyPointer(values ...*string) *string {
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}
		return &trimmed
	}
	return nil
}

func normalizeTwoLetterValue(value *string, fallback string) *string {
	if value == nil {
		return &fallback
	}

	normalized := strings.ToUpper(strings.TrimSpace(*value))
	if normalized == "" {
		return &fallback
	}
	if len(normalized) > 2 {
		normalized = normalized[:2]
	}
	return &normalized
}

func stringRef(value string) *string {
	return &value
}

func resolveClientServerVerificationURL(form *platon.ClientServerVerificationForm) (*url.URL, error) {
	logger := log.NewLogger("Platon Verification: ")

	if form == nil {
		err := fmt.Errorf("verification form is nil")
		logger.Error("%v", err)
		return nil, err
	}

	values := url.Values{}
	for key, value := range form.Fields {
		values.Set(key, value)
	}
	encodedForm := values.Encode()
	logger.Debug("Endpoint: %s", form.Endpoint)
	logger.Debug("Fields count: %d", len(form.Fields))
	logger.Debug(
		"Request (%s):\n%s",
		internalhttp.FormURLEncodedContentType,
		internalhttp.PrettyPrintFormURLEncodedBody(encodedForm),
	)

	req, err := http.NewRequest(http.MethodPost, form.Endpoint, strings.NewReader(encodedForm))
	if err != nil {
		err = fmt.Errorf("cannot build verification request: %w", err)
		logger.Error("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", internalhttp.FormURLEncodedContentType)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("verification request failed: %w", err)
		logger.Error("%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	logger.Debug("Response status: %d", resp.StatusCode)

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		err = fmt.Errorf("cannot read verification response body: %w", err)
		logger.Error("%v", err)
		return nil, err
	}
	logger.Debug("Response body size: %d bytes", len(body))
	if len(body) == 0 {
		logger.Debug("Response: <empty>")
	} else if internalhttp.IsFormURLEncodedContentType(resp.Header.Get("Content-Type")) {
		logger.Debug(
			"Response (%s):\n%s",
			internalhttp.FormURLEncodedContentType,
			truncateVerificationBodyForLog([]byte(internalhttp.PrettyPrintFormURLEncodedBody(string(body)))),
		)
	} else {
		logger.Debug("Response: %s", truncateVerificationBodyForLog(body))
	}

	if location := strings.TrimSpace(resp.Header.Get("Location")); location != "" {
		logger.Debug("Response location: %s", location)
		return parsePurchaseURL(location)
	}

	absRe := regexp.MustCompile(`https://secure\.platononline\.com/payment/purchase\?token=[A-Za-z0-9]+`)
	if match := absRe.Find(body); match != nil {
		return parsePurchaseURL(string(match))
	}

	relRe := regexp.MustCompile(`/payment/purchase\?token=[A-Za-z0-9]+`)
	if match := relRe.Find(body); match != nil {
		return parsePurchaseURL("https://secure.platononline.com" + string(match))
	}

	errMsg := fmt.Sprintf("verification purchase URL was not returned (status=%d)", resp.StatusCode)
	if bytes.Contains(bytes.ToLower(body), []byte("<title>error")) {
		errMsg += "; gateway returned error page (check merchant key, secret/signature, and callback URL)"
	}

	logger.Error("%s", errMsg)
	return nil, errors.New(errMsg)
}

func truncateVerificationBodyForLog(raw []byte) string {
	const max = 512
	if len(raw) <= max {
		return string(raw)
	}
	return string(raw[:max]) + "...(truncated)"
}

func parsePurchaseURL(raw string) (*url.URL, error) {
	parsedURL, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse verification URL %q: %w", raw, err)
	}
	if !parsedURL.IsAbs() {
		return nil, fmt.Errorf("verification URL is not absolute: %q", raw)
	}
	return parsedURL, nil
}

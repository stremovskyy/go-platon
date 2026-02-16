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
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/stremovskyy/go-platon/log"
)

var orderAmountRe = regexp.MustCompile("^[0-9]+\\.[0-9]{2}$")

// Request represents the main payment request structure
type Request struct {
	Action           string  `json:"action" validate:"omitempty,oneof=SALE GET_TRANS_STATUS GET_TRANS_STATUS_BY_ORDER APPLEPAY GOOGLEPAY CAPTURE CREDITVOID CREDIT2CARD GET_SUBMERCHANT"`
	ClientKey        string  `json:"client_key" validate:"required"`
	Hash             string  `json:"hash,omitempty" validate:"omitempty,len=32"`
	ChannelId        string  `json:"channel_id,omitempty" validate:"omitempty,max=255"`
	PayerIp          *string `json:"payer_ip,omitempty" validate:"omitempty,ipv4"`
	TermUrl3ds       *string `json:"term_url_3ds,omitempty" validate:"omitempty,max=1024,url"`
	OrderID          *string `json:"order_id,omitempty" validate:"omitempty,max=255"`
	OrderAmount      string  `json:"order_amount,omitempty" validate:"omitempty"`
	OrderCurrency    string  `json:"order_currency,omitempty" validate:"omitempty,alpha,len=3"`
	SubmerchantID    *string `json:"submerchant_id,omitempty" validate:"omitempty,max=255"`
	OrderDescription *string `json:"order_description,omitempty" validate:"omitempty,max=1024"`

	// Apple Pay / Google Pay request payload (base64 string, formatted per IA docs).
	PaymentToken *string `json:"payment_token,omitempty" validate:"omitempty"`

	PayerEmail     *string `json:"payer_email,omitempty" validate:"omitempty,email,max=256"`
	PayerPhone     *string `json:"payer_phone,omitempty" validate:"omitempty,numeric,startswith=380,max=32"`
	PayerFirstName *string `json:"payer_first_name,omitempty" validate:"omitempty,max=32"`
	PayerLastName  *string `json:"payer_last_name,omitempty" validate:"omitempty,max=32"`
	PayerAddress   *string `json:"payer_address,omitempty" validate:"omitempty,max=256"`
	PayerCountry   *string `json:"payer_country,omitempty" validate:"omitempty,max=2"`
	PayerState     *string `json:"payer_state,omitempty" validate:"omitempty,max=2"`
	PayerCity      *string `json:"payer_city,omitempty" validate:"omitempty,max=32"`
	PayerZip       *string `json:"payer_zip,omitempty" validate:"omitempty,max=32"`
	CustomerWallet *string `json:"customer_wallet,omitempty" validate:"omitempty,max=255"`
	CardNumber     *string `json:"card_number,omitempty" validate:"omitempty,numeric,len=16"`
	CardExpMonth   *string `json:"card_exp_month,omitempty" validate:"omitempty,numeric,len=2"`
	CardExpYear    *string `json:"card_exp_year,omitempty" validate:"omitempty,numeric,len=4"`
	CardCvv2       *string `json:"card_cvv2,omitempty" validate:"omitempty,numeric,len=3"`
	CardToken      *string `json:"card_token,omitempty" validate:"omitempty,max=32"`

	// "auth" parameter: Y to create HOLD (preauth), N for normal SALE.
	AuthFlag *string `json:"auth,omitempty" validate:"omitempty,oneof=Y N"`

	// Recurring payment: first transaction trans_id.
	RecurringFirstTransID *string `json:"recurring_first_trans_id,omitempty" validate:"omitempty,max=32"`

	// GET_TRANS_STATUS request trans_id.
	TransId *string `json:"trans_id,omitempty" validate:"omitempty,max=32"`

	// CAPTURE / CREDITVOID amount.
	Amount string `json:"amount,omitempty" validate:"omitempty"`

	// CREDITVOID: fast refund flag.
	Immediately *string `json:"immediately,omitempty" validate:"omitempty,oneof=Y"`

	ReqToken      *string `json:"req_token,omitempty" validate:"omitempty,oneof=Y N"`
	RecurringInit *string `json:"recurring_init,omitempty" validate:"omitempty,oneof=Y N"`
	Async         *string `json:"async,omitempty" validate:"omitempty,oneof=Y N"`

	Ext1  *string `json:"ext1,omitempty" validate:"omitempty,max=1024"`
	Ext2  *string `json:"ext2,omitempty" validate:"omitempty,max=1024"`
	Ext3  *string `json:"ext3,omitempty" validate:"omitempty,max=1024"`
	Ext4  *string `json:"ext4,omitempty" validate:"omitempty,max=1024"`
	Ext5  *string `json:"ext5,omitempty" validate:"omitempty,max=1024"`
	Ext6  *string `json:"ext6,omitempty" validate:"omitempty,max=1024"`
	Ext7  *string `json:"ext7,omitempty" validate:"omitempty,max=1024"`
	Ext8  *string `json:"ext8,omitempty" validate:"omitempty,max=1024"`
	Ext9  *string `json:"ext9,omitempty" validate:"omitempty,max=1024"`
	Ext10 *string `json:"ext10,omitempty" validate:"omitempty,max=1024"`

	// Optional split distribution rules for SALE/CAPTURE/CREDITVOID.
	SplitRules SplitRules `json:"split_rules,omitempty" validate:"omitempty"`

	// HashEmail is an internal helper for signature generation for CAPTURE/CREDITVOID/GET_TRANS_STATUS.
	// Per IA docs, it is not sent to Platon and may be empty if not specified in the initial payment.
	HashEmail *string `json:"-"`

	Auth     *Auth    `json:"-"`
	HashType HashType `json:"-"`
}

// NewPaymentRequest creates a new validated payment request
func (r *Request) SignAndPrepare() (*Request, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	var sign string
	var err error

	switch r.HashType {
	case HashTypeVerification:
		sign, err = r.generateCardPanSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeCardPayment:
		sign, err = r.generateCardPanSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeCardTokenPayment:
		sign, err = r.generateCardTokenSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeApplePay, HashTypeGooglePay:
		sign, err = r.generatePaymentTokenSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeRecurring:
		sign, err = r.generateRecurringSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeGetTransStatus, HashTypeCapture, HashTypeCreditVoid:
		sign, err = r.generateTransIDSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeGetTransStatusByOrder:
		sign, err = r.generateGetTransStatusByOrderSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeGetSubmerchant:
		sign, err = r.generateGetSubmerchantSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeCredit2Card:
		sign, err = r.generateCredit2CardSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	case HashTypeCredit2CardToken:
		sign, err = r.generateCredit2CardTokenSignature()
		if err != nil {
			return nil, fmt.Errorf("signature generation failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown hash type: %s", r.HashType)
	}

	r.Hash = sign

	if err := r.validateByHashType(); err != nil {
		return nil, err
	}

	// Validate request
	if err := validator.New().Struct(r); err != nil {
		return nil, fmt.Errorf("internal request validation failed: %w", err)
	}

	return r, nil
}

func (r *Request) SignForAction(t HashType) *Request {
	if r == nil {
		return nil
	}

	r.HashType = t

	return r
}

func (r *Request) generateSignature(signArray []string) (string, error) {
	// Create a logger instance with a custom prefix.
	logger := log.NewLogger("PlatonSignature")

	logger.All("Generating signature with property keys: %v", signArray)

	var concatenated string

	for _, key := range signArray {
		var value string
		switch key {
		case "key":
			if r.Auth == nil {
				return "", fmt.Errorf("auth is nil; cannot retrieve key")
			}
			value = r.Auth.Key
		case "pass":
			if r.Auth == nil {
				return "", fmt.Errorf("auth is nil; cannot retrieve secret")
			}
			value = r.Auth.Secret
		default:
			fieldValue, err := getFieldValueByJSONTag(r, key)
			if err != nil {
				return "", err
			}
			value = fieldValue
		}

		// Reverse the string value.
		reversed := reverseString(value)

		logger.All("Key '%s': original='%s', reversed='%s'", key, value, reversed)

		concatenated += reversed
	}

	// Log the concatenated reversed string.
	logger.All("Concatenated reversed string: %s", concatenated)

	// Convert to uppercase.
	upperConcatenated := strings.ToUpper(concatenated)
	logger.All("Uppercased string: %s", upperConcatenated)

	// Compute the MD5 hash.
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateCardPanSignature() (string, error) {
	// Create a logger instance with a custom prefix
	logger := log.NewLogger("CardPanSignature")
	logger.All("Generating signature for payment request")

	// Validate required fields for hash generation
	if r.PayerEmail == nil {
		return "", fmt.Errorf("payer_email is required for signature generation")
	}
	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.CardNumber == nil {
		return "", fmt.Errorf("card_number is required for signature generation")
	}

	// Extract card number first 6 and last 4 digits
	cardNumber := *r.CardNumber
	if len(cardNumber) < 10 {
		return "", fmt.Errorf("card_number is too short")
	}
	cardFirst6 := cardNumber[0:6]
	cardLast4 := cardNumber[len(cardNumber)-4:]

	// Reverse strings according to PHP implementation
	reversedEmail := reverseString(*r.PayerEmail)
	reversedCard := reverseString(cardFirst6 + cardLast4)

	// Log the components
	logger.All("Components: email='%s', card='%s'", reversedEmail, reversedCard)

	// Concatenate according to PHP implementation:
	// strrev(email) + client_pass + strrev(first6+last4)
	concatenated := reversedEmail + r.Auth.Secret + reversedCard

	// Convert to uppercase
	upperConcatenated := strings.ToUpper(concatenated)
	logger.All("Uppercased concatenated string: %s", upperConcatenated)

	// Compute the MD5 hash
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateCardTokenSignature() (string, error) {
	logger := log.NewLogger("CardTokenSignature")
	logger.All("Generating signature for card_token request")

	if r.PayerEmail == nil {
		return "", fmt.Errorf("payer_email is required for signature generation")
	}
	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.CardToken == nil || *r.CardToken == "" {
		return "", fmt.Errorf("card_token is required for signature generation")
	}

	reversedEmail := reverseString(*r.PayerEmail)
	reversedToken := reverseString(*r.CardToken)
	concatenated := reversedEmail + r.Auth.Secret + reversedToken

	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generatePaymentTokenSignature() (string, error) {
	logger := log.NewLogger("PaymentTokenSignature")
	logger.All("Generating signature for payment_token request")

	if r.PayerEmail == nil {
		return "", fmt.Errorf("payer_email is required for signature generation")
	}
	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.PaymentToken == nil || *r.PaymentToken == "" {
		return "", fmt.Errorf("payment_token is required for signature generation")
	}

	reversedEmail := reverseString(*r.PayerEmail)
	reversedToken := reverseString(*r.PaymentToken)
	concatenated := reversedEmail + r.Auth.Secret + reversedToken

	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateRecurringSignature() (string, error) {
	// Per IA docs, recurring payments by token use the same signature as one-click by CARD_TOKEN.
	return r.generateCardTokenSignature()
}

func (r *Request) generateTransIDSignature() (string, error) {
	logger := log.NewLogger("TransIDSignature")
	logger.All("Generating signature for trans_id based request")

	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.TransId == nil || *r.TransId == "" {
		return "", fmt.Errorf("trans_id is required for signature generation")
	}

	// "email" used in signature per IA docs. It is not sent to Platon and may be empty.
	email := ""
	if r.HashEmail != nil {
		email = *r.HashEmail
	} else if r.PayerEmail != nil {
		// Backward-compatible fallback if caller provided payer_email only.
		email = *r.PayerEmail
	}

	reversedEmail := reverseString(email)
	concatenated := reversedEmail + r.Auth.Secret + *r.TransId

	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateGetTransStatusByOrderSignature() (string, error) {
	logger := log.NewLogger("GetTransStatusByOrderSignature")
	logger.All("Generating signature for GET_TRANS_STATUS_BY_ORDER request")

	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.OrderID == nil || *r.OrderID == "" {
		return "", fmt.Errorf("order_id is required for signature generation")
	}

	concatenated := *r.OrderID + r.Auth.Secret
	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateGetSubmerchantSignature() (string, error) {
	logger := log.NewLogger("GetSubmerchantSignature")
	logger.All("Generating signature for GET_SUBMERCHANT request")

	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.SubmerchantID == nil || *r.SubmerchantID == "" {
		return "", fmt.Errorf("submerchant_id is required for signature generation")
	}

	// Per IA docs:
	// md5(strtoupper(client_pass + submerchant_id))
	concatenated := r.Auth.Secret + *r.SubmerchantID
	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateCredit2CardSignature() (string, error) {
	logger := log.NewLogger("Credit2CardSignature")
	logger.All("Generating signature for CREDIT2CARD request by PAN")

	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.CardNumber == nil || *r.CardNumber == "" {
		return "", fmt.Errorf("card_number is required for signature generation")
	}

	cardNumber := *r.CardNumber
	if len(cardNumber) < 10 {
		return "", fmt.Errorf("card_number is too short")
	}
	cardHashPart := cardNumber[0:6] + cardNumber[len(cardNumber)-4:]

	reversedCardHash := reverseString(cardHashPart)
	concatenated := r.Auth.Secret + reversedCardHash
	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) generateCredit2CardTokenSignature() (string, error) {
	logger := log.NewLogger("Credit2CardTokenSignature")
	logger.All("Generating signature for CREDIT2CARD request by card token")

	if r.Auth == nil || r.Auth.Secret == "" {
		return "", fmt.Errorf("Auth secret is required for signature generation")
	}
	if r.CardToken == nil || *r.CardToken == "" {
		return "", fmt.Errorf("card_token is required for signature generation")
	}

	reversedToken := reverseString(*r.CardToken)
	concatenated := r.Auth.Secret + reversedToken
	upperConcatenated := strings.ToUpper(concatenated)
	hash := md5.Sum([]byte(upperConcatenated))
	signature := hex.EncodeToString(hash[:])
	logger.All("Generated MD5 signature: %s", signature)

	return signature, nil
}

func (r *Request) ToMap() map[string]interface{} {
	if r == nil {
		return map[string]interface{}{}
	}

	requestMap := make(map[string]interface{})

	v := reflect.ValueOf(*r)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the json tag name
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// Split the tag to get the name part (before any comma)
		tagName := strings.Split(tag, ",")[0]

		fieldValue := v.Field(i)

		// Skip empty values
		if fieldValue.IsZero() {
			continue
		}

		// Handle pointer types
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		// Add to map using the json tag name
		requestMap[tagName] = fieldValue.Interface()
	}

	return requestMap
}

func (r *Request) validateByHashType() error {
	switch r.HashType {
	case HashTypeVerification:
		// Per IA docs, verification requests must explicitly request tokenization + recurring init.
		if r.ReqToken == nil {
			r.ReqToken = refString("Y")
		}
		if r.RecurringInit == nil {
			r.RecurringInit = refString("Y")
		}

		if r.Action != ActionCodeSALE.String() {
			return fmt.Errorf("verification: action must be %s", ActionCodeSALE.String())
		}
		if r.ChannelId != "VERIFY_ZERO" {
			return fmt.Errorf("verification: channel_id must be VERIFY_ZERO")
		}
		if r.OrderAmount != VerifyNoAmount.String() {
			return fmt.Errorf("verification: order_amount must be %s", VerifyNoAmount.String())
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("verification: order_id is required")
		}
		if len(*r.OrderID) > 32 {
			return fmt.Errorf("verification: order_id must be <= 32 characters")
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("verification: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("verification: order_description is required")
		}
		if len(*r.OrderDescription) > 255 {
			return fmt.Errorf("verification: order_description must be <= 255 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("verification: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("verification: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 255 {
			return fmt.Errorf("verification: term_url_3ds must be <= 255 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("verification: payer_email is required")
		}
		if r.PayerPhone == nil || *r.PayerPhone == "" {
			return fmt.Errorf("verification: payer_phone is required")
		}
		if r.CardNumber == nil || *r.CardNumber == "" {
			return fmt.Errorf("verification: card_number is required")
		}
		if r.CardExpMonth == nil || *r.CardExpMonth == "" {
			return fmt.Errorf("verification: card_exp_month is required")
		}
		if r.CardExpYear == nil || *r.CardExpYear == "" {
			return fmt.Errorf("verification: card_exp_year is required")
		}
		if r.CardCvv2 == nil || *r.CardCvv2 == "" {
			return fmt.Errorf("verification: card_cvv2 is required")
		}
		if r.ReqToken == nil || *r.ReqToken == "" {
			return fmt.Errorf("verification: req_token is required")
		}
		if *r.ReqToken != "Y" {
			return fmt.Errorf("verification: req_token must be Y")
		}
		if r.RecurringInit == nil || *r.RecurringInit == "" {
			return fmt.Errorf("verification: recurring_init is required")
		}
		if *r.RecurringInit != "Y" {
			return fmt.Errorf("verification: recurring_init must be Y")
		}

	case HashTypeCardPayment:
		// Per IA docs, card payments require req_token/recurring_init flags to be explicitly present (Y/N).
		if r.ReqToken == nil {
			r.ReqToken = refString("N")
		}
		if r.RecurringInit == nil {
			r.RecurringInit = refString("N")
		}

		if r.Action != ActionCodeSALE.String() {
			return fmt.Errorf("card_payment: action must be %s", ActionCodeSALE.String())
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("card_payment: order_id is required")
		}
		if len(*r.OrderID) > 32 {
			return fmt.Errorf("card_payment: order_id must be <= 32 characters")
		}
		if r.OrderAmount == "" {
			return fmt.Errorf("card_payment: order_amount is required")
		}
		if !orderAmountRe.MatchString(r.OrderAmount) {
			return fmt.Errorf("card_payment: order_amount must match %q (got %q)", orderAmountRe.String(), r.OrderAmount)
		}
		if v, err := parseOrderAmountMinorUnits(r.OrderAmount); err != nil || v <= 0 {
			return fmt.Errorf("card_payment: order_amount must be > 0 (got %q)", r.OrderAmount)
		}
		if err := validateSplitRules(r.SplitRules, r.OrderAmount, "card_payment"); err != nil {
			return err
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("card_payment: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("card_payment: order_description is required")
		}
		if len(*r.OrderDescription) > 255 {
			return fmt.Errorf("card_payment: order_description must be <= 255 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("card_payment: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("card_payment: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 255 {
			return fmt.Errorf("card_payment: term_url_3ds must be <= 255 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("card_payment: payer_email is required")
		}
		if r.PayerPhone == nil || *r.PayerPhone == "" {
			return fmt.Errorf("card_payment: payer_phone is required")
		}
		if r.CardNumber == nil || *r.CardNumber == "" {
			return fmt.Errorf("card_payment: card_number is required")
		}
		if r.CardExpMonth == nil || *r.CardExpMonth == "" {
			return fmt.Errorf("card_payment: card_exp_month is required")
		}
		if r.CardExpYear == nil || *r.CardExpYear == "" {
			return fmt.Errorf("card_payment: card_exp_year is required")
		}
		if r.CardCvv2 == nil || *r.CardCvv2 == "" {
			return fmt.Errorf("card_payment: card_cvv2 is required")
		}
		if r.ReqToken == nil || *r.ReqToken == "" {
			return fmt.Errorf("card_payment: req_token is required")
		}
		if r.RecurringInit == nil || *r.RecurringInit == "" {
			return fmt.Errorf("card_payment: recurring_init is required")
		}

	case HashTypeCardTokenPayment:
		if r.Action != ActionCodeSALE.String() {
			return fmt.Errorf("card_token_payment: action must be %s", ActionCodeSALE.String())
		}
		if r.CardToken == nil || *r.CardToken == "" {
			return fmt.Errorf("card_token_payment: card_token is required")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("card_token_payment: order_id is required")
		}
		if len(*r.OrderID) > 32 {
			return fmt.Errorf("card_token_payment: order_id must be <= 32 characters")
		}
		if r.OrderAmount == "" {
			return fmt.Errorf("card_token_payment: order_amount is required")
		}
		if !orderAmountRe.MatchString(r.OrderAmount) {
			return fmt.Errorf("card_token_payment: order_amount must match %q (got %q)", orderAmountRe.String(), r.OrderAmount)
		}
		if v, err := parseOrderAmountMinorUnits(r.OrderAmount); err != nil || v <= 0 {
			return fmt.Errorf("card_token_payment: order_amount must be > 0 (got %q)", r.OrderAmount)
		}
		if err := validateSplitRules(r.SplitRules, r.OrderAmount, "card_token_payment"); err != nil {
			return err
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("card_token_payment: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("card_token_payment: order_description is required")
		}
		if len(*r.OrderDescription) > 255 {
			return fmt.Errorf("card_token_payment: order_description must be <= 255 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("card_token_payment: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("card_token_payment: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 255 {
			return fmt.Errorf("card_token_payment: term_url_3ds must be <= 255 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("card_token_payment: payer_email is required")
		}

	case HashTypeApplePay:
		if r.Action != ActionCodeAPPLEPAY.String() {
			return fmt.Errorf("apple_pay: action must be %s", ActionCodeAPPLEPAY.String())
		}
		if r.PaymentToken == nil || *r.PaymentToken == "" {
			return fmt.Errorf("apple_pay: payment_token is required")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("apple_pay: order_id is required")
		}
		if len(*r.OrderID) > 255 {
			return fmt.Errorf("apple_pay: order_id must be <= 255 characters")
		}
		if r.OrderAmount == "" {
			return fmt.Errorf("apple_pay: order_amount is required")
		}
		if !orderAmountRe.MatchString(r.OrderAmount) {
			return fmt.Errorf("apple_pay: order_amount must match %q (got %q)", orderAmountRe.String(), r.OrderAmount)
		}
		if v, err := parseOrderAmountMinorUnits(r.OrderAmount); err != nil || v <= 0 {
			return fmt.Errorf("apple_pay: order_amount must be > 0 (got %q)", r.OrderAmount)
		}
		if err := validateSplitRules(r.SplitRules, r.OrderAmount, "apple_pay"); err != nil {
			return err
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("apple_pay: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("apple_pay: order_description is required")
		}
		if len(*r.OrderDescription) > 1024 {
			return fmt.Errorf("apple_pay: order_description must be <= 1024 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("apple_pay: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("apple_pay: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 1024 {
			return fmt.Errorf("apple_pay: term_url_3ds must be <= 1024 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("apple_pay: payer_email is required")
		}
		if r.PayerPhone == nil || *r.PayerPhone == "" {
			return fmt.Errorf("apple_pay: payer_phone is required")
		}

	case HashTypeGooglePay:
		if r.Action != ActionCodeGOOGLEPAY.String() {
			return fmt.Errorf("google_pay: action must be %s", ActionCodeGOOGLEPAY.String())
		}
		if r.PaymentToken == nil || *r.PaymentToken == "" {
			return fmt.Errorf("google_pay: payment_token is required")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("google_pay: order_id is required")
		}
		if len(*r.OrderID) > 255 {
			return fmt.Errorf("google_pay: order_id must be <= 255 characters")
		}
		if r.OrderAmount == "" {
			return fmt.Errorf("google_pay: order_amount is required")
		}
		if !orderAmountRe.MatchString(r.OrderAmount) {
			return fmt.Errorf("google_pay: order_amount must match %q (got %q)", orderAmountRe.String(), r.OrderAmount)
		}
		if v, err := parseOrderAmountMinorUnits(r.OrderAmount); err != nil || v <= 0 {
			return fmt.Errorf("google_pay: order_amount must be > 0 (got %q)", r.OrderAmount)
		}
		if err := validateSplitRules(r.SplitRules, r.OrderAmount, "google_pay"); err != nil {
			return err
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("google_pay: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("google_pay: order_description is required")
		}
		if len(*r.OrderDescription) > 255 {
			return fmt.Errorf("google_pay: order_description must be <= 255 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("google_pay: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("google_pay: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 255 {
			return fmt.Errorf("google_pay: term_url_3ds must be <= 255 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("google_pay: payer_email is required")
		}
		if r.PayerPhone == nil || *r.PayerPhone == "" {
			return fmt.Errorf("google_pay: payer_phone is required")
		}

	case HashTypeRecurring:
		if r.Action != ActionCodeSALE.String() {
			return fmt.Errorf("recurring: action must be %s", ActionCodeSALE.String())
		}
		if r.CardToken == nil || *r.CardToken == "" {
			return fmt.Errorf("recurring: card_token is required")
		}
		if r.Ext3 == nil || *r.Ext3 != "recurring" {
			return fmt.Errorf("recurring: ext3 must be \"recurring\"")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("recurring: order_id is required")
		}
		if len(*r.OrderID) > 32 {
			return fmt.Errorf("recurring: order_id must be <= 32 characters")
		}
		if r.OrderAmount == "" {
			return fmt.Errorf("recurring: order_amount is required")
		}
		if !orderAmountRe.MatchString(r.OrderAmount) {
			return fmt.Errorf("recurring: order_amount must match %q (got %q)", orderAmountRe.String(), r.OrderAmount)
		}
		if v, err := parseOrderAmountMinorUnits(r.OrderAmount); err != nil || v <= 0 {
			return fmt.Errorf("recurring: order_amount must be > 0 (got %q)", r.OrderAmount)
		}
		if err := validateSplitRules(r.SplitRules, r.OrderAmount, "recurring"); err != nil {
			return err
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("recurring: order_currency is required")
		}
		if r.OrderDescription == nil || *r.OrderDescription == "" {
			return fmt.Errorf("recurring: order_description is required")
		}
		if len(*r.OrderDescription) > 255 {
			return fmt.Errorf("recurring: order_description must be <= 255 characters")
		}
		if r.PayerIp == nil || *r.PayerIp == "" {
			return fmt.Errorf("recurring: payer_ip is required")
		}
		if r.TermUrl3ds == nil || *r.TermUrl3ds == "" {
			return fmt.Errorf("recurring: term_url_3ds is required")
		}
		if len(*r.TermUrl3ds) > 255 {
			return fmt.Errorf("recurring: term_url_3ds must be <= 255 characters")
		}
		if r.PayerEmail == nil || *r.PayerEmail == "" {
			return fmt.Errorf("recurring: payer_email is required")
		}

	case HashTypeGetTransStatus:
		if r.Action != ActionCodeGetTransStatus.String() {
			return fmt.Errorf("get_trans_status: action must be %s", ActionCodeGetTransStatus.String())
		}
		if r.TransId == nil || *r.TransId == "" {
			return fmt.Errorf("get_trans_status: trans_id is required")
		}

	case HashTypeGetTransStatusByOrder:
		if r.Action != ActionCodeGetTransStatusByOrder.String() {
			return fmt.Errorf("get_trans_status_by_order: action must be %s", ActionCodeGetTransStatusByOrder.String())
		}
		if r.OrderID == nil || strings.TrimSpace(*r.OrderID) == "" {
			return fmt.Errorf("get_trans_status_by_order: order_id is required")
		}

	case HashTypeCapture:
		if r.Action != ActionCodeCAPTURE.String() {
			return fmt.Errorf("capture: action must be %s", ActionCodeCAPTURE.String())
		}
		if r.TransId == nil || *r.TransId == "" {
			return fmt.Errorf("capture: trans_id is required")
		}
		if r.Amount == "" {
			return fmt.Errorf("capture: amount is required")
		}
		if !orderAmountRe.MatchString(r.Amount) {
			return fmt.Errorf("capture: amount must match %q (got %q)", orderAmountRe.String(), r.Amount)
		}
		if v, err := parseOrderAmountMinorUnits(r.Amount); err != nil || v <= 0 {
			return fmt.Errorf("capture: amount must be > 0 (got %q)", r.Amount)
		}
		if err := validateSplitRules(r.SplitRules, r.Amount, "capture"); err != nil {
			return err
		}

	case HashTypeCreditVoid:
		if r.Action != ActionCodeCREDITVOID.String() {
			return fmt.Errorf("creditvoid: action must be %s", ActionCodeCREDITVOID.String())
		}
		if r.TransId == nil || *r.TransId == "" {
			return fmt.Errorf("creditvoid: trans_id is required")
		}
		if r.Amount == "" {
			return fmt.Errorf("creditvoid: amount is required")
		}
		if !orderAmountRe.MatchString(r.Amount) {
			return fmt.Errorf("creditvoid: amount must match %q (got %q)", orderAmountRe.String(), r.Amount)
		}
		if v, err := parseOrderAmountMinorUnits(r.Amount); err != nil || v <= 0 {
			return fmt.Errorf("creditvoid: amount must be > 0 (got %q)", r.Amount)
		}
		if err := validateSplitRules(r.SplitRules, r.Amount, "creditvoid"); err != nil {
			return err
		}

	case HashTypeCredit2Card:
		if r.Action != ActionCodeCREDIT2CARD.String() {
			return fmt.Errorf("credit2card: action must be %s", ActionCodeCREDIT2CARD.String())
		}
		if r.CardNumber == nil || *r.CardNumber == "" {
			return fmt.Errorf("credit2card: card_number is required")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("credit2card: order_id is required")
		}
		if r.Amount == "" {
			return fmt.Errorf("credit2card: amount is required")
		}
		if !orderAmountRe.MatchString(r.Amount) {
			return fmt.Errorf("credit2card: amount must match %q (got %q)", orderAmountRe.String(), r.Amount)
		}
		if v, err := parseOrderAmountMinorUnits(r.Amount); err != nil || v <= 0 {
			return fmt.Errorf("credit2card: amount must be > 0 (got %q)", r.Amount)
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("credit2card: order_currency is required")
		}
		if r.OrderDescription == nil || strings.TrimSpace(*r.OrderDescription) == "" {
			return fmt.Errorf("credit2card: order_description is required")
		}
		if r.PayerFirstName == nil || strings.TrimSpace(*r.PayerFirstName) == "" {
			return fmt.Errorf("credit2card: payer_first_name is required")
		}
		if r.PayerLastName == nil || strings.TrimSpace(*r.PayerLastName) == "" {
			return fmt.Errorf("credit2card: payer_last_name is required")
		}
		if r.PayerAddress == nil || strings.TrimSpace(*r.PayerAddress) == "" {
			return fmt.Errorf("credit2card: payer_address is required")
		}
		if r.PayerCountry == nil || strings.TrimSpace(*r.PayerCountry) == "" {
			return fmt.Errorf("credit2card: payer_country is required")
		}
		if r.PayerState == nil || strings.TrimSpace(*r.PayerState) == "" {
			return fmt.Errorf("credit2card: payer_state is required")
		}
		if r.PayerCity == nil || strings.TrimSpace(*r.PayerCity) == "" {
			return fmt.Errorf("credit2card: payer_city is required")
		}
		if r.PayerZip == nil || strings.TrimSpace(*r.PayerZip) == "" {
			return fmt.Errorf("credit2card: payer_zip is required")
		}
		if len(r.SplitRules) > 0 {
			return fmt.Errorf("credit2card: split_rules are not allowed")
		}

	case HashTypeCredit2CardToken:
		if r.Action != ActionCodeCREDIT2CARD.String() {
			return fmt.Errorf("credit2card_token: action must be %s", ActionCodeCREDIT2CARD.String())
		}
		if r.CardToken == nil || *r.CardToken == "" {
			return fmt.Errorf("credit2card_token: card_token is required")
		}
		if r.OrderID == nil || *r.OrderID == "" {
			return fmt.Errorf("credit2card_token: order_id is required")
		}
		if r.Amount == "" {
			return fmt.Errorf("credit2card_token: amount is required")
		}
		if !orderAmountRe.MatchString(r.Amount) {
			return fmt.Errorf("credit2card_token: amount must match %q (got %q)", orderAmountRe.String(), r.Amount)
		}
		if v, err := parseOrderAmountMinorUnits(r.Amount); err != nil || v <= 0 {
			return fmt.Errorf("credit2card_token: amount must be > 0 (got %q)", r.Amount)
		}
		if r.OrderCurrency == "" {
			return fmt.Errorf("credit2card_token: order_currency is required")
		}
		if r.OrderDescription == nil || strings.TrimSpace(*r.OrderDescription) == "" {
			return fmt.Errorf("credit2card_token: order_description is required")
		}
		if r.PayerFirstName == nil || strings.TrimSpace(*r.PayerFirstName) == "" {
			return fmt.Errorf("credit2card_token: payer_first_name is required")
		}
		if r.PayerLastName == nil || strings.TrimSpace(*r.PayerLastName) == "" {
			return fmt.Errorf("credit2card_token: payer_last_name is required")
		}
		if r.PayerAddress == nil || strings.TrimSpace(*r.PayerAddress) == "" {
			return fmt.Errorf("credit2card_token: payer_address is required")
		}
		if r.PayerCountry == nil || strings.TrimSpace(*r.PayerCountry) == "" {
			return fmt.Errorf("credit2card_token: payer_country is required")
		}
		if r.PayerState == nil || strings.TrimSpace(*r.PayerState) == "" {
			return fmt.Errorf("credit2card_token: payer_state is required")
		}
		if r.PayerCity == nil || strings.TrimSpace(*r.PayerCity) == "" {
			return fmt.Errorf("credit2card_token: payer_city is required")
		}
		if r.PayerZip == nil || strings.TrimSpace(*r.PayerZip) == "" {
			return fmt.Errorf("credit2card_token: payer_zip is required")
		}
		if len(r.SplitRules) > 0 {
			return fmt.Errorf("credit2card_token: split_rules are not allowed")
		}

	case HashTypeGetSubmerchant:
		if r.Action != ActionCodeGetSubmerchant.String() {
			return fmt.Errorf("get_submerchant: action must be %s", ActionCodeGetSubmerchant.String())
		}
		if r.SubmerchantID == nil || strings.TrimSpace(*r.SubmerchantID) == "" {
			return fmt.Errorf("get_submerchant: submerchant_id is required")
		}
		if len(r.SplitRules) > 0 {
			return fmt.Errorf("get_submerchant: split_rules are not allowed")
		}
	}

	return nil
}

func refString(value string) *string {
	return &value
}

func validateSplitRules(rules SplitRules, totalAmount string, context string) error {
	if len(rules) == 0 {
		return nil
	}
	if totalAmount == "" {
		return fmt.Errorf("%s: amount is required when split_rules are provided", context)
	}

	totalMinorUnits, err := parseOrderAmountMinorUnits(totalAmount)
	if err != nil || totalMinorUnits <= 0 {
		return fmt.Errorf("%s: invalid amount %q for split_rules", context, totalAmount)
	}

	splitMinorUnits := 0
	for submerchantID, amount := range rules {
		if strings.TrimSpace(submerchantID) == "" {
			return fmt.Errorf("%s: split_rules key (submerchant_id) is required", context)
		}

		if !orderAmountRe.MatchString(amount) {
			return fmt.Errorf("%s: split_rules[%q] amount must match %q (got %q)", context, submerchantID, orderAmountRe.String(), amount)
		}
		minorUnits, parseErr := parseOrderAmountMinorUnits(amount)
		if parseErr != nil || minorUnits <= 0 {
			return fmt.Errorf("%s: split_rules[%q] amount must be > 0 (got %q)", context, submerchantID, amount)
		}
		splitMinorUnits += minorUnits
	}
	if splitMinorUnits != totalMinorUnits {
		return fmt.Errorf("%s: split rules total must equal amount (%d != %d minor units)", context, splitMinorUnits, totalMinorUnits)
	}

	return nil
}

func parseOrderAmountMinorUnits(amount string) (int, error) {
	parts := strings.SplitN(amount, ".", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid amount format")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major amount")
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minor amount")
	}
	return major*100 + minor, nil
}

// getFieldValueByJSONTag uses reflection to search for a struct field whose "json" tag (or field name)
// matches the provided key. It returns the field's string representation.
func getFieldValueByJSONTag(obj interface{}, key string) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("object is nil")
	}

	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return "", fmt.Errorf("object is invalid")
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", fmt.Errorf("object pointer is nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("object must be a struct or pointer to struct")
	}

	t := v.Type()

	// Iterate over all fields looking for a matching json tag.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			// The tag may include options (e.g. "payment,omitempty"), so split by comma.
			parts := strings.Split(tag, ",")
			if parts[0] == key {
				return fmt.Sprintf("%v", v.Field(i).Interface()), nil
			}
		}
		// Fallback: check if the lower-cased field name matches the key.
		if strings.ToLower(field.Name) == key {
			return fmt.Sprintf("%v", v.Field(i).Interface()), nil
		}
	}
	return "", fmt.Errorf("field with json tag or name '%s' not found", key)
}

// reverseString returns the reversed version of the provided string.
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

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

type Lang string

func (l *Lang) String() string {
	if l == nil {
		return ""
	}

	return string(*l)
}

const (
	LangUk Lang = "UK"
	LangEn Lang = "EN"
)

type CurrencyKey string

const (
	CurrencyUAH CurrencyKey = "UAH"
	CurrencyUSD CurrencyKey = "USD"
	CurrencyEUR CurrencyKey = "EUR"
)

// type FormID string
//
// func (a FormID) String() string {
// 	return string(a)
// }
//
// // FormID types for Request
// const (
// 	FormIDVerify FormID = "verify"
//
// 	ActionGetPaymentStatus FormID = "GetPaymentStatus"
// 	ActionDebiting         FormID = "Debiting"
// 	ActionCompletion       FormID = "Completion"
// 	ActionReversal         FormID = "Reversal"
// 	ActionCredit           FormID = "A2CPay"
// 	MobilePaymentCreate    FormID = "PaymentCreate"
// )

type FixedAmount string

func (a FixedAmount) String() string {
	return string(a)
}

const (
	VerifyFixedAmount FixedAmount = "1.00"
	VerifyNoAmount    FixedAmount = "0.40"
)

type ActionCode string

func (a ActionCode) String() string {
	return string(a)
}

const (
	ActionCodeSALE           ActionCode = "SALE"
	ActionCodeGetTransStatus ActionCode = "GET_TRANS_STATUS"
	ActionCodeAPPLEPAY       ActionCode = "APPLEPAY"
	ActionCodeGOOGLEPAY      ActionCode = "GOOGLEPAY"
	ActionCodeCAPTURE        ActionCode = "CAPTURE"
	ActionCodeCREDITVOID     ActionCode = "CREDITVOID"
	ActionCodeGetSubmerchant ActionCode = "GET_SUBMERCHANT"
)

type HashType string

func (h HashType) String() string {
	return string(h)
}

const (
	// HashTypeVerification is used for card verification (IA "Верифікація картки").
	HashTypeVerification HashType = "verification"

	// HashTypeCardPayment is used for IA payments by PAN (card_number + exp + cvv2).
	HashTypeCardPayment HashType = "card_payment"

	// HashTypeCardTokenPayment is used for IA payments by CARD_TOKEN (one-click).
	HashTypeCardTokenPayment HashType = "card_token_payment"

	// HashTypeApplePay is used for IA Apple Pay payments.
	HashTypeApplePay HashType = "apple_pay"

	// HashTypeGooglePay is used for IA Google Pay payments.
	HashTypeGooglePay HashType = "google_pay"

	// HashTypeRecurring is used for IA recurring payments.
	HashTypeRecurring HashType = "recurring"

	// HashTypeGetTransStatus is used for the GET_TRANS_STATUS request.
	HashTypeGetTransStatus HashType = "get_trans_status"

	// HashTypeCapture is used for CAPTURE (confirm HOLD).
	HashTypeCapture HashType = "capture"

	// HashTypeCreditVoid is used for CREDITVOID (refund).
	HashTypeCreditVoid HashType = "creditvoid"

	// HashTypeGetSubmerchant is used for GET_SUBMERCHANT requests.
	HashTypeGetSubmerchant HashType = "get_submerchant"
)

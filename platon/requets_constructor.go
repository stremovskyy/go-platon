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
	"fmt"

	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/internal/utils"
)

func NewRequest(action ActionCode) *Request {
	req := &Request{
		Action: action.String(),
	}

	return req
}

func (r *Request) WithAuth(auth *Auth) *Request {
	if r == nil {
		return nil
	}

	r.Auth = auth

	return r
}

func (r *Request) WithClientKey(clientKey string) *Request {
	if r == nil {
		return nil
	}

	r.ClientKey = clientKey
	return r
}

func (r *Request) WithReqToken(flag bool) *Request {
	if r == nil {
		return nil
	}

	if flag {
		r.ReqToken = utils.Ref("Y")
	} else {
		r.ReqToken = utils.Ref("N")
	}
	return r
}

func (r *Request) WithRecToken() *Request {
	if r == nil {
		return nil
	}

	r.ReqToken = utils.Ref("Y")

	return r
}

func (r *Request) WithRecurringInitFlag(flag bool) *Request {
	if r == nil {
		return nil
	}

	if flag {
		r.RecurringInit = utils.Ref("Y")
	} else {
		r.RecurringInit = utils.Ref("N")
	}
	return r
}

func (r *Request) WithRecurringInit() *Request {
	if r == nil {
		return nil
	}

	r.RecurringInit = utils.Ref("Y")

	return r
}

func (r *Request) WithAsync(flag bool) *Request {
	if r == nil {
		return nil
	}

	if flag {
		r.Async = utils.Ref("Y")
	} else {
		r.Async = utils.Ref("N")
	}
	return r
}

func (r *Request) UseAsync() *Request {
	if r == nil {
		return nil
	}

	r.Async = utils.Ref("Y")

	return r
}

func (r *Request) WithChannelNoAmountVerification() *Request {
	if r == nil {
		return nil
	}

	r.ChannelId = "VERIFY_ZERO"

	return r
}

func (r *Request) WithPayerIP(ip *string) *Request {
	if r == nil {
		return nil
	}

	if ip == nil {
		r.PayerIp = utils.Ref("127.0.0.1")
	} else {
		r.PayerIp = ip
	}

	return r
}

func (r *Request) WithTermsURL(url *string) *Request {
	if r == nil {
		return nil
	}

	r.TermUrl3ds = url

	return r
}

func (r *Request) WithCardNumber(pan *string) *Request {
	if r == nil {
		return nil
	}

	r.CardNumber = pan

	return r
}

func (r *Request) WithCardToken(token *string) *Request {
	if r == nil {
		return nil
	}

	r.CardToken = token

	return r
}

func (r *Request) WithCardExpMonth(month *string) *Request {
	if r == nil {
		return nil
	}

	r.CardExpMonth = month

	return r
}

func (r *Request) WithCardExpYear(year *string) *Request {
	if r == nil {
		return nil
	}

	r.CardExpYear = year

	return r
}

func (r *Request) WithCardCvv2(cvv2 *string) *Request {
	if r == nil {
		return nil
	}

	r.CardCvv2 = cvv2

	return r
}

func (r *Request) WithPayerEmail(email *string) *Request {
	if r == nil {
		return nil
	}

	r.PayerEmail = email

	return r
}

func (r *Request) WithPayerPhone(phone *string) *Request {
	if r == nil {
		return nil
	}

	r.PayerPhone = phone

	return r
}

func (r *Request) WithPayerFirstName(firstName *string) *Request {
	if r == nil {
		return nil
	}

	r.PayerFirstName = firstName
	return r
}

func (r *Request) WithPayerLastName(lastName *string) *Request {
	if r == nil {
		return nil
	}

	r.PayerLastName = lastName
	return r
}

func (r *Request) WithApplePayData(data *string) *Request {
	if r == nil {
		return nil
	}

	// Backward-compatible helper. IA docs use the `payment_token` parameter for Apple Pay.
	r.PaymentToken = data
	return r
}

func (r *Request) WithGooglePayToken(token *string) *Request {
	if r == nil {
		return nil
	}

	// IA docs use the `payment_token` parameter for Google Pay.
	r.PaymentToken = token
	return r
}

func (r *Request) WithPaymentToken(token *string) *Request {
	if r == nil {
		return nil
	}

	r.PaymentToken = token
	return r
}

func (r *Request) WithHoldAuth() *Request {
	if r == nil {
		return nil
	}

	r.AuthFlag = utils.Ref("Y")
	return r
}

func (r *Request) WithVerifyAmount(amount float32) *Request {
	if r == nil {
		return nil
	}

	r.OrderAmount = fmt.Sprintf("%.2f", amount)

	if amount <= 0 {
		r.OrderAmount = VerifyNoAmount.String()
	}

	return r
}

func (r *Request) WithOrderAmountMinorUnits(amount int) *Request {
	if r == nil {
		return nil
	}

	// amount is in minor units (e.g. kopecks); Platon expects a decimal string with 2 digits.
	r.OrderAmount = fmt.Sprintf("%.2f", float64(amount)/100)
	return r
}

func (r *Request) WithOrderAmount(amount string) *Request {
	if r == nil {
		return nil
	}

	r.OrderAmount = amount
	return r
}

func (r *Request) ForCurrency(currency currency.Code) *Request {
	if r == nil {
		return nil
	}

	r.OrderCurrency = currency.String()
	return r
}

func (r *Request) WithSubmerchantID(submerchantID *string) *Request {
	if r == nil {
		return nil
	}

	r.SubmerchantID = submerchantID
	return r
}

func (r *Request) WithDescription(description string) *Request {
	if r == nil {
		return nil
	}

	r.OrderDescription = &description

	return r
}

func (r *Request) WithOrderID(orderID *string) *Request {
	if r == nil {
		return nil
	}

	r.OrderID = orderID

	return r
}

func (r *Request) WithRecurringFirstTransID(transID *string) *Request {
	if r == nil {
		return nil
	}

	r.RecurringFirstTransID = transID
	return r
}

func (r *Request) WithTransID(transID *string) *Request {
	if r == nil {
		return nil
	}

	r.TransId = transID
	return r
}

func (r *Request) WithAmountMinorUnits(amount int) *Request {
	if r == nil {
		return nil
	}

	// amount is in minor units (e.g. kopecks); Platon expects a decimal string with 2 digits.
	r.Amount = fmt.Sprintf("%.2f", float64(amount)/100)
	return r
}

func (r *Request) WithAmount(amount string) *Request {
	if r == nil {
		return nil
	}

	r.Amount = amount
	return r
}

func (r *Request) WithSplitRules(splitRules SplitRules) *Request {
	if r == nil {
		return nil
	}

	if len(splitRules) == 0 {
		r.SplitRules = nil
		return r
	}
	r.SplitRules = splitRules
	return r
}

func (r *Request) WithImmediately(flag bool) *Request {
	if r == nil {
		return nil
	}

	if flag {
		r.Immediately = utils.Ref("Y")
	} else {
		r.Immediately = nil
	}
	return r
}

// WithHashEmail sets the email used for signature generation for CAPTURE/CREDITVOID/GET_TRANS_STATUS.
// This value is not sent to Platon (json:"-").
func (r *Request) WithHashEmail(email *string) *Request {
	if r == nil {
		return nil
	}

	r.HashEmail = email
	return r
}

// WithCardHashPart sets the first6+last4 part used for signature generation for CAPTURE/CREDITVOID/GET_TRANS_STATUS.
// This value is not sent to Platon (json:"-").
func (r *Request) WithCardHashPart(cardHashPart *string) *Request {
	if r == nil {
		return nil
	}

	r.CardHashPart = cardHashPart
	return r
}

func (r *Request) WithExt3(value *string) *Request {
	if r == nil {
		return nil
	}

	r.Ext3 = value
	return r
}

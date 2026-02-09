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

package consts

const (
	Version    = "1.0.0"
	ApiVersion = "1.28"

	baseUrl = "https://secure.platononline.com"

	ApiPaymentURL = baseUrl + "/payment"

	// ApiConfigurationURL is the IA configuration endpoint (e.g. GET_SUBMERCHANT).
	ApiConfigurationURL = baseUrl + "/configuration/"

	// ApiPostURL is the IA endpoint for Apple Pay and Google Pay.
	ApiPostURL = baseUrl + "/post/"

	// ApiPostUnqURL is the IA Server-Server endpoint for card payments, verification, one-click,
	// recurring by token, capture/refund, and status.
	ApiPostUnqURL = baseUrl + "/post-unq/"

	// ApiVerifyURL is the legacy name for the IA Server-Server endpoint (`/post-unq/`).
	// It is used both for card verification and card/token payments.
	ApiVerifyURL = ApiPostUnqURL

	// Backward-compatible aliases (deprecated names).
	ApiApplePayURL    = ApiPostURL
	ApiGooglePayURL   = ApiPostURL
	ApiRecurringURL   = ApiPostUnqURL
	ApiGetTransStatus = ApiPostUnqURL
	ApiGetSubmerchant = ApiConfigurationURL
)

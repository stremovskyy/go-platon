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
	"net/url"

	"github.com/stremovskyy/go-platon/log"
	"github.com/stremovskyy/go-platon/platon"
)

// Platon is the public client interface.
//
// Methods accept optional RunOption values (for example DryRun()).
// Verification executes client-server verification and returns ready-to-use purchase URL.
type Platon interface {
	Verification(request *Request, opts ...RunOption) (*url.URL, error)
	VerificationLink(request *Request, opts ...RunOption) (*url.URL, error)
	Status(request *Request, opts ...RunOption) (*platon.Response, error)
	Payment(request *Request, opts ...RunOption) (*platon.Response, error)
	Hold(request *Request, opts ...RunOption) (*platon.Response, error)
	SubmerchantAvailableForSplit(request *Request, opts ...RunOption) (bool, error)
	Capture(request *Request, opts ...RunOption) (*platon.Response, error)
	Refund(request *Request, opts ...RunOption) (*platon.Response, error)
	Credit(request *Request, opts ...RunOption) (*platon.Response, error)
	// Deprecated: Platon production callbacks use application/x-www-form-urlencoded.
	// Use go_platon.ParseWebhookForm for callback parsing and signature verification.
	ParseWebhookXML(data []byte) (*platon.Payment, error)
	SetLogLevel(levelDebug log.Level)
}

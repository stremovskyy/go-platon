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

package platon

import (
	"strings"
	"testing"
)

func TestUnmarshalJSONResponse_SubmerchantStatus(t *testing.T) {
	raw := []byte(`{"status":"success","response":{"submerchant_id":"12345678","submerchant_id_status":"ENABLED"}}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	status, ok := resp.SubmerchantIDStatus()
	if !ok {
		t.Fatalf("expected submerchant status payload")
	}
	if status != "ENABLED" {
		t.Fatalf("expected ENABLED, got %q", status)
	}
}

func TestUnmarshalJSONResponse_SubmerchantStatusTopLevel(t *testing.T) {
	raw := []byte(`{"status":"SUCCESS","action":"GET_SUBMERCHANT","submerchant_id":"12345678","submerchant_id_status":"ENABLED","hash":"abc123"}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	status, ok := resp.SubmerchantIDStatus()
	if !ok {
		t.Fatalf("expected submerchant status payload")
	}
	if status != "ENABLED" {
		t.Fatalf("expected ENABLED, got %q", status)
	}
	if resp.ResponseData == nil || resp.ResponseData.Hash == nil || *resp.ResponseData.Hash != "abc123" {
		t.Fatalf("expected top-level hash to be mapped into response data")
	}
}

func TestUnmarshalJSONResponse_AllowsNullErrorMessage(t *testing.T) {
	raw := []byte(`{"result":"ACCEPTED","error_message":null}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	if resp.ErrorMessage != "" {
		t.Fatalf("expected empty error message, got %q", resp.ErrorMessage)
	}
	if gotErr := resp.GetError(); gotErr != nil {
		t.Fatalf("expected nil error, got %v", gotErr)
	}
}

func TestUnmarshalJSONResponse_DeclinedReasonReturnsError(t *testing.T) {
	raw := []byte(`{"result":"DECLINED","decline_reason":"102: Token is not active","error_message":null}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	if resp.DeclineReason != "102: Token is not active" {
		t.Fatalf("unexpected decline reason: %q", resp.DeclineReason)
	}

	gotErr := resp.GetError()
	if gotErr == nil {
		t.Fatalf("expected decline error, got nil")
	}
	if !strings.Contains(gotErr.Error(), "102: Token is not active") {
		t.Fatalf("expected decline reason in error, got %q", gotErr.Error())
	}
}

func TestResponse_GetError_DeclinedWithoutReason(t *testing.T) {
	declined := ResultDeclined
	resp := &Response{
		Result: &declined,
	}

	gotErr := resp.GetError()
	if gotErr == nil {
		t.Fatalf("expected decline error, got nil")
	}
	if !strings.Contains(gotErr.Error(), "unknown platon api decline") {
		t.Fatalf("unexpected error: %q", gotErr.Error())
	}
}

func TestUnmarshalJSONResponse_ErrorMessageObject(t *testing.T) {
	raw := []byte(`{"result":"ERROR","error_message":{"field":"Wrong cardholder_email"}}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	if resp.ErrorMessage != `{"field":"Wrong cardholder_email"}` {
		t.Fatalf("unexpected normalized error message: %q", resp.ErrorMessage)
	}

	gotErr := resp.GetError()
	if gotErr == nil {
		t.Fatalf("expected api error, got nil")
	}
	if !strings.Contains(gotErr.Error(), "Wrong cardholder_email") {
		t.Fatalf("expected parsed object in error, got %q", gotErr.Error())
	}
}

func TestUnmarshalJSONResponse_StatusByOrderPayload(t *testing.T) {
	raw := []byte(`{"action":"GET_TRANS_STATUS_BY_ORDER","orders":[{"amount":"101.55","date":"2026-03-19 07:14:31","order_id":"58036844-17ef-48a0-8484-848320f461fe","status":"SETTLED","trans_id":"47390-44719-73004"}],"result":"SUCCESS"}`)

	resp, err := UnmarshalJSONResponse(raw)
	if err != nil {
		t.Fatalf("UnmarshalJSONResponse() error: %v", err)
	}

	if len(resp.Orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(resp.Orders))
	}
	if resp.Status == nil || *resp.Status != "SETTLED" {
		t.Fatalf("expected top-level status to be populated from orders[0], got %v", resp.Status)
	}
	if resp.OrderId == nil || *resp.OrderId != "58036844-17ef-48a0-8484-848320f461fe" {
		t.Fatalf("unexpected order_id: %v", resp.OrderId)
	}
	if resp.TransId == nil || *resp.TransId != "47390-44719-73004" {
		t.Fatalf("unexpected trans_id: %v", resp.TransId)
	}
	if gotErr := resp.GetError(); gotErr != nil {
		t.Fatalf("expected nil error, got %v", gotErr)
	}
}

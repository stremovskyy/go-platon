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

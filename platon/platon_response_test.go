package platon

import "testing"

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

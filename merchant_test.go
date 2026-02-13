package go_platon

import "testing"

func TestMerchant_NilReceiverMethods(t *testing.T) {
	var merchant *Merchant

	if got := merchant.GetMerchantID(); got != nil {
		t.Fatalf("GetMerchantID() mismatch: want nil, got %v", *got)
	}
	if got := merchant.GetMobileLogin(); got != nil {
		t.Fatalf("GetMobileLogin() mismatch: want nil, got %q", *got)
	}
}

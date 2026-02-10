package go_platon

import (
	"math"
	"testing"
)

func TestRequest_GetAmount_UsesMinorUnits(t *testing.T) {
	req := &Request{
		PaymentData: &PaymentData{
			Amount: 199,
		},
	}

	got := req.GetAmount()
	if math.Abs(float64(got)-1.99) > 1e-6 {
		t.Fatalf("GetAmount() mismatch: want 1.99, got %.8f", got)
	}
}

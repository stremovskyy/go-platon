package platon

import "testing"

func TestPayment_GetTransactionByID_ReturnsPointerToUnderlyingElement(t *testing.T) {
	payment := &Payment{
		Transactions: Transactions{
			Transaction: []Transaction{
				{ID: 1001, Desc: "first"},
				{ID: 1002, Desc: "second"},
			},
		},
	}

	tx := payment.GetTransactionByID(1002)
	if tx == nil {
		t.Fatalf("GetTransactionByID() returned nil for existing transaction")
	}

	tx.Desc = "updated"
	if payment.Transactions.Transaction[1].Desc != "updated" {
		t.Fatalf("GetTransactionByID() returned pointer to copy, expected pointer to slice element")
	}
}

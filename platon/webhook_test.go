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

func TestPayment_NilReceiver_IsSafe(t *testing.T) {
	var payment *Payment

	if payment.IsValid() {
		t.Fatalf("nil payment must be invalid")
	}
	if tx := payment.GetTransactionByID(1); tx != nil {
		t.Fatalf("expected nil transaction for nil payment receiver")
	}
	if got := payment.String(); got != "Payment[<nil>]" {
		t.Fatalf("unexpected String() for nil payment: %q", got)
	}
}

func TestTransactions_NilReceiver_IsSafe(t *testing.T) {
	var txs *Transactions

	if got := txs.Len(); got != 0 {
		t.Fatalf("Len() mismatch: want 0, got %d", got)
	}
	if tx := txs.First(); tx != nil {
		t.Fatalf("First() expected nil transaction for nil receiver")
	}
	if tx := txs.Last(); tx != nil {
		t.Fatalf("Last() expected nil transaction for nil receiver")
	}
}

package go_platon

import (
	"testing"
)

const webhookFormPayload = "id=47097-87770-07123&order=47097-87309-6110&status=SALE&card=411111%2A%2A%2A%2A1111&description=test&amount=0.40&currency=UAH&email=&date=2026-02-13+10%3A32%3A57&ip=250.137.176.130&sign=582d658d7d422e76b2639fac131d093e"

func TestParseWebhookForm(t *testing.T) {
	form, err := ParseWebhookForm([]byte(webhookFormPayload))
	if err != nil {
		t.Fatalf("ParseWebhookForm() error: %v", err)
	}

	if form.Order != "47097-87309-6110" {
		t.Fatalf("order mismatch: got %q", form.Order)
	}
	if form.Status != "SALE" {
		t.Fatalf("status mismatch: got %q", form.Status)
	}
	if form.Card != "411111****1111" {
		t.Fatalf("card mismatch: got %q", form.Card)
	}
}

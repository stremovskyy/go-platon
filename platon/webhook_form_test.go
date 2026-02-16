package platon

import (
	"strings"
	"testing"
)

const webhookFormPayload = "id=47097-87770-07123&order=47097-87309-6110&status=SALE&card=411111%2A%2A%2A%2A1111&description=%D0%9F%D0%BE%D0%BF%D0%BE%D0%B2%D0%BD%D0%B5%D0%BD%D0%BD%D1%8F+%D0%B1%D0%B0%D0%BB%D0%B0%D0%BD%D1%81%D1%83+%D0%B2%D0%BE%D0%B4%D1%96%D1%8F+%28Platon+split+one+receiver%29&amount=0.40&currency=UAH&name=+&phone=&email=&date=2026-02-13+10%3A32%3A57&ip=250.137.176.130&sign=582d658d7d422e76b2639fac131d093e&rc_id=47097-87770-07123&rc_token=fa0500fb3f4869247b4c5532eaf799bc&issuing_bank=JPMORGAN+CHASE+BANK%2C+N.A.&ext4=&cardholder_email=&brand=VISA&terminal="

func TestParseWebhookForm(t *testing.T) {
	form, err := ParseWebhookForm([]byte(webhookFormPayload))
	if err != nil {
		t.Fatalf("ParseWebhookForm() error: %v", err)
	}

	if form.ID != "47097-87770-07123" {
		t.Fatalf("id mismatch: got %q", form.ID)
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
	if form.Amount != "0.40" {
		t.Fatalf("amount mismatch: got %q", form.Amount)
	}
	if form.Currency != "UAH" {
		t.Fatalf("currency mismatch: got %q", form.Currency)
	}
	if !strings.Contains(form.Description, "Platon split one receiver") {
		t.Fatalf("description was not decoded correctly: got %q", form.Description)
	}
	if form.Date != "2026-02-13 10:32:57" {
		t.Fatalf("date mismatch: got %q", form.Date)
	}
	if form.Sign != "582d658d7d422e76b2639fac131d093e" {
		t.Fatalf("sign mismatch: got %q", form.Sign)
	}
	if form.RCToken != "fa0500fb3f4869247b4c5532eaf799bc" {
		t.Fatalf("rc_token mismatch: got %q", form.RCToken)
	}
	if form.IssuingBank != "JPMORGAN CHASE BANK, N.A." {
		t.Fatalf("issuing_bank mismatch: got %q", form.IssuingBank)
	}
}

func TestWebhookForm_ExpectedSignAndVerify(t *testing.T) {
	form, err := ParseWebhookForm([]byte(webhookFormPayload))
	if err != nil {
		t.Fatalf("ParseWebhookForm() error: %v", err)
	}

	expected, err := form.ExpectedSign("SECRET", "payer@example.com")
	if err != nil {
		t.Fatalf("ExpectedSign() error: %v", err)
	}
	if expected != "8c089577f40387dd2a0c5f91b1b703c8" {
		t.Fatalf("expected signature mismatch: got %q", expected)
	}

	form.Sign = expected
	ok, err := form.VerifySign("SECRET", "payer@example.com")
	if err != nil {
		t.Fatalf("VerifySign() error: %v", err)
	}
	if !ok {
		t.Fatalf("VerifySign() expected true")
	}

	ok, err = form.VerifySign("WRONG_SECRET", "payer@example.com")
	if err != nil {
		t.Fatalf("VerifySign() with wrong secret error: %v", err)
	}
	if ok {
		t.Fatalf("VerifySign() expected false for wrong secret")
	}
}

func TestWebhookForm_ExpectedSign_UsesCallbackEmailWhenOverrideIsEmpty(t *testing.T) {
	form := &WebhookForm{
		Order:  "order-1",
		Status: "SALE",
		Card:   "411111****1111",
		Email:  "payer@example.com",
	}

	fromOverride, err := form.ExpectedSign("SECRET", "payer@example.com")
	if err != nil {
		t.Fatalf("ExpectedSign() with override error: %v", err)
	}

	fromCallbackEmail, err := form.ExpectedSign("SECRET", "")
	if err != nil {
		t.Fatalf("ExpectedSign() with callback email error: %v", err)
	}

	if fromOverride != fromCallbackEmail {
		t.Fatalf("ExpectedSign() mismatch: override=%q callback=%q", fromOverride, fromCallbackEmail)
	}
}

func TestWebhookCardSignSource_Validation(t *testing.T) {
	if _, err := webhookCardSignSource("1234"); err == nil {
		t.Fatalf("expected error for short card")
	}

	got, err := webhookCardSignSource("411111 **** 1111")
	if err != nil {
		t.Fatalf("webhookCardSignSource() error: %v", err)
	}
	if got != "4111111111" {
		t.Fatalf("card sign source mismatch: got %q", got)
	}
}

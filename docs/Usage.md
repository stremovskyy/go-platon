# GO-Platon Usage Guide

This guide describes the current public API of the library and how it maps to Platon IA (Server-Server) documentation.

## Quick Start (Card PAN Payment)

```go
package main

import (
	"fmt"
	"time"

	go_platon "github.com/stremovskyy/go-platon"
	"github.com/stremovskyy/go-platon/currency"
	"github.com/stremovskyy/go-platon/internal/utils"
)

func main() {
	client := go_platon.NewClient(
		go_platon.WithTimeout(30*time.Second),
	)

	merchant := &go_platon.Merchant{
		MerchantID:  "4767",
		MerchantKey: "CLIENT_KEY",
		SecretKey:   "CLIENT_PASS",
		TermsURL:    utils.Ref("https://example.com/3ds-term"),
	}

	orderID := "order-123"
	req := &go_platon.Request{
		Merchant: merchant,
		PaymentData: &go_platon.PaymentData{
			PaymentID:   utils.Ref(orderID), // Platon "order_id"
			Amount:      100,                // minor units (e.g. 100 -> 1.00 UAH)
			Currency:    currency.UAH,
			Description: "Test payment",
		},
		PaymentMethod: &go_platon.PaymentMethod{
			Card: &go_platon.Card{
				Pan:             utils.Ref("4111111111111111"),
				ExpirationMonth: utils.Ref("01"),
				ExpirationYear:  utils.Ref("2026"),
				Cvv2:            utils.Ref("123"),
			},
		},
		PersonalData: &go_platon.PersonalData{
			Email: utils.Ref("payer@example.com"),
			Phone: utils.Ref("380631234567"),
		},
	}

	resp, err := client.Payment(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result=%v TransID=%v Error=%q\n", resp.Result, resp.TransId, resp.ErrorMessage)
}
```

## One-Click Payment (CARD_TOKEN)

Set `PaymentMethod.Card.Token` instead of PAN/expiry/CVV:

```go
req.PaymentMethod = &go_platon.PaymentMethod{
	Card: &go_platon.Card{
		Token: utils.Ref("CARD_TOKEN"),
	},
}

resp, err := client.Payment(req)
```

Runnable example: `examples/card_token/card_token.go`.

## Apple Pay / Google Pay

- Apple Pay: set `PaymentMethod.AppleContainer` (base64 string of the Apple container).
- Google Pay: set `PaymentMethod.GoogleToken` (base64 string of the Google Pay token).

Then call `client.Payment(req)` or `client.Hold(req)`.

## Card Verification (Client-Server)

Card verification must use Client-Server flow (`/payment/auth`) and be submitted from payer browser.

The simplest way is to call:

```go
verificationURL, err := client.Verification(req)
if err != nil {
	panic(err)
}
fmt.Println(verificationURL.String())
```

`client.VerificationLink(req)` is an alias with the same behavior.

If you need full control over HTML/form rendering, use
`go_platon.BuildClientServerVerificationForm(req)` and submit returned fields manually.

## Webhook Callback (`application/x-www-form-urlencoded`)

Platon uses a single callback URL for all payment flows.
Recommended integration pattern:

- send routing context in `ext1..ext10`
- receive all callbacks on one frontend webhook endpoint
- verify `sign`
- route internally by `ext*` values

`go-platon` maps `PaymentData.Metadata["ext1"]..["ext10"]` to Platon request fields `ext1..ext10`.

```mermaid
flowchart LR
    A["Your backend: create payment/hold/capture/refund"] -->| "ext4=wallet-topup" | B["Platon API"]
    B --> C["Single callback URL: /platon/webhook"]
    C --> D{"Verify sign"}
    D -->| "invalid" | E["Reject callback"]
    D -->| "valid" | F{"ext4 value"}
    F -->| "wallet-topup" | G["Wallet service handler"]
    F -->| "order-payment" | H["Orders service handler"]
    F -->| "other" | I["Default/manual queue"]
```

Example: send route marker in payment request:

```go
req.PaymentData.Metadata = map[string]string{
	"ext4": "wallet-topup",
}
```

Then parse callback payload and route:

```go
body, err := io.ReadAll(r.Body)
if err != nil {
	panic(err)
}

form, err := go_platon.ParseWebhookForm(body)
if err != nil {
	panic(err)
}

ok, err := form.VerifySign("CLIENT_PASS", "payer@example.com")
if err != nil {
	panic(err)
}
if !ok {
	panic("invalid webhook sign")
}

switch form.Ext4 {
case "wallet-topup":
	// forward to wallet callback handler
case "order-payment":
	// forward to order callback handler
default:
	// fallback handler
}
```

## GET_TRANS_STATUS_BY_ORDER

`client.Status(req)` sends `GET_TRANS_STATUS_BY_ORDER`.

Required:

- `PaymentData.PaymentID` (merchant `order_id`)

Signature uses `order_id + client_pass` (uppercase MD5).

## GET_SUBMERCHANT

`client.SubmerchantAvailableForSplit(req)` sends `GET_SUBMERCHANT` to IA `/configuration/` and returns `true` when
`submerchant_id_status=ENABLED`.

Required:

- `Merchant.MerchantKey`
- `Merchant.SecretKey`
- `PaymentData.SubmerchantID`

Example:

```go
enabled, err := client.SubmerchantAvailableForSplit(&go_platon.Request{
	Merchant: merchant,
	PaymentData: &go_platon.PaymentData{
		SubmerchantID: utils.Ref("12345678"),
	},
})
```

## Split Rules (`split_rules`)

For `Payment`/`Hold` (`SALE`, including Apple Pay/Google Pay), `Capture` (`CAPTURE`), and `Refund` (`CREDITVOID`),
you can pass split distribution in `PaymentData.SplitRules`. Amounts are in minor units.

```go
req.PaymentData.Amount = 1500
req.PaymentData.SplitRules = []go_platon.SplitRule{
	{SubmerchantIdentification: "submerchant_01", Amount: 1000}, // 10.00
	{SubmerchantIdentification: "submerchant_02", Amount: 500},  // 5.00
}
```

The total split amount must be equal to `PaymentData.Amount`.
The SDK serializes this as `split_rules={"submerchant_01":"10.00","submerchant_02":"5.00"}`.

## CAPTURE (Confirm HOLD)

`client.Capture(req)` sends a `CAPTURE` request (confirm a HOLD/preauth) to IA `/post-unq/`.

Required:

- `PaymentData.PlatonTransID` (or legacy `PaymentData.PlatonPaymentID`)
- `PaymentData.Amount` (minor units, e.g. 100 -> 1.00)

Optional:

- `PersonalData.Email` (signature-only)

## CREDITVOID (Refund)

`client.Refund(req)` sends a `CREDITVOID` request (refund) to IA `/post-unq/`.

Required:

- `PaymentData.PlatonTransID` (or legacy `PaymentData.PlatonPaymentID`)
- `PaymentData.Amount` (minor units, e.g. 100 -> 1.00)

Optional:

- `PersonalData.Email` (signature-only)
- `PaymentData.Metadata["immediately"]` set to `Y`/`true`/`1` to send `immediately=Y` (fast refund)

## CREDIT2CARD (A2C payout)

`client.Credit(req)` sends an A2C payout request to `/p2p-unq/` with `action=CREDIT2CARD`.

Required:

- `PaymentData.PaymentID` (order_id)
- `PaymentData.Amount` (minor units, e.g. 100 -> 1.00)
- `PaymentData.Currency`
- `PaymentData.Description`
- `PaymentMethod.Card.Token`

Payer identity fields required by A2C (`payer_first_name`, `payer_last_name`, `payer_address`,
`payer_country`, `payer_state`, `payer_city`, `payer_zip`) are taken from request data when provided,
or filled with safe defaults.

## A2C Status

`client.Status(req)` supports A2C status checks over `/p2p-unq/` when
`PaymentData.Metadata["platon_flow"] == "a2c"`.

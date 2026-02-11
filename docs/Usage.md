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

## Card Verification (VERIFY_ZERO)

`client.Verification(req)` sends a verification request with:

- `channel_id=VERIFY_ZERO`
- `order_amount=0.40`
- `req_token=Y` and `recurring_init=Y`

It returns `*platon.Result`.

## GET_TRANS_STATUS

`client.Status(req)` requires:

- `PaymentData.PlatonTransID` (Platon `trans_id`, string) or `PaymentData.PlatonPaymentID` (legacy int64)
- `PaymentMethod.Card.Pan` (only the first 6 and last 4 are used for signature generation)

`PersonalData.Email` is used only for signature generation and is not sent to Platon. If the email was not provided
in the initial payment request, IA docs allow signing with an empty email.

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
- `PaymentMethod.Card.Pan` (only first 6 + last 4 are used for signature)

Optional:

- `PersonalData.Email` (signature-only)

## CREDITVOID (Refund)

`client.Refund(req)` sends a `CREDITVOID` request (refund) to IA `/post-unq/`.

Required:

- `PaymentData.PlatonTransID` (or legacy `PaymentData.PlatonPaymentID`)
- `PaymentData.Amount` (minor units, e.g. 100 -> 1.00)
- `PaymentMethod.Card.Pan` (only first 6 + last 4 are used for signature)

Optional:

- `PersonalData.Email` (signature-only)
- `PaymentData.Metadata["immediately"]` set to `Y`/`true`/`1` to send `immediately=Y` (fast refund)

## Limitations

`Credit` is not implemented yet.

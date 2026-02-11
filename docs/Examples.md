# GO-Platon Examples

Runnable examples are available under `examples/`.
Examples load credentials from environment variables using `examples/internal/config`.
Copy `examples/.env.example` to `examples/.env` and set real values before running examples.

## Card Verification

See `examples/verification/verification.go`.

## Card Payment / Hold

See:

- `examples/payment/payment.go`
- `examples/hold/hold.go`

## One-Click Payment (CARD_TOKEN)

See `examples/card_token/card_token.go`.

## Apple Pay / Google Pay

See:

- `examples/apple_pay/apple_pay.go`
- `examples/google_pay/google_pay.go`

## GET_TRANS_STATUS

See `examples/status/status.go`.

## Split / GET_SUBMERCHANT

See `examples/split/split.go`.

## CAPTURE / CREDITVOID

See:

- `examples/capture/capture.go`
- `examples/refund/refund.go`

# GO-Platon Examples

Runnable examples are available under `examples/`.
Examples load credentials from environment variables using `examples/internal/config`.
Copy `examples/.env.example` to `examples/.env` and set real values before running examples.
Demo card/token/email values used by payment examples are defined in `examples/internal/demo/data.go`.

## Card Verification (Client-Server Form)

See `examples/verification/verification.go` (calls `client.Verification(req)` and prints URL in format
`https://secure.platononline.com/payment/purchase?token=...`).

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

## Webhook Callback (form-urlencoded)

See `examples/webhook/vebhook.go`.

## Split / GET_SUBMERCHANT

See `examples/split/split.go`.

## CAPTURE / CREDITVOID

See:

- `examples/capture/capture.go`
- `examples/refund/refund.go`

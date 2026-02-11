# GO-Platon

[![GoDoc](https://godoc.org/github.com/stremovskyy/go-platon?status.svg)](https://godoc.org/github.com/stremovskyy/go-platon)
[![Go Report Card](https://goreportcard.com/badge/github.com/stremovskyy/go-platon)](https://goreportcard.com/report/github.com/stremovskyy/go-platon)
[![License](https://img.shields.io/github/license/stremovskyy/go-platon)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/stremovskyy/go-platon)](go.mod)

A Go client library for the Platon.ua Payment Gateway API. This library provides a simple and reliable way to integrate Platon.ua payment processing into your Go applications.

## Features

- IA (Server-Server) card payments by PAN or `CARD_TOKEN`
- IA (Server-Server) card verification (`VERIFY_ZERO`, `order_amount=0.40`)
- IA Apple Pay / Google Pay (Server-Server)
- GET_TRANS_STATUS request support
- CAPTURE (confirm HOLD) support
- CREDITVOID (refund) support
- Split payouts (`split_rules`) for SALE/CAPTURE/CREDITVOID
- GET_SUBMERCHANT support
- URL-encoded requests (`application/x-www-form-urlencoded`)
- Structured logging and optional request/response recording

## Installation

```bash
go get github.com/stremovskyy/go-platon
```

## Documentation

 **Complete documentation is available in:**
- [Usage Guide](docs/Usage.md) - Comprehensive guide to using the library
- [Code Examples](docs/Examples.md) - Real-world examples and best practices

For API reference and package documentation:
- [GoDoc Documentation](https://godoc.org/github.com/stremovskyy/go-platon)

## Quick Start

```go
package main

import (
    "github.com/stremovskyy/go-platon"
)

func main() {
    // Create a new client with default options
    client := go_platon.NewDefaultClient()

    // Configure merchant details
    merchant := &go_platon.Merchant{
        Name:        "Your Merchant Name",
        MerchantID:  "your-merchant-id",
        MerchantKey: "your-merchant-key",
        SecretKey:   "your-merchant-secret",
    }

    // Use the client...
}
```

## Configuration Options

The client can be configured with various options:

```go
client := go_platon.NewClient(
    WithTimeout(time.Second * 30),
    WithKeepAlive(time.Second * 30),
)
```

## Error Handling

Most API/validation issues are returned as `error` with context (wrapping `platon.Error` where applicable).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to update tests as appropriate.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Platon.ua for providing the payment gateway service
- Contributors who have helped improve this library

## Support

If you are having problems, please let us know by [raising a new issue](https://github.com/stremovskyy/go-platon/issues/new).

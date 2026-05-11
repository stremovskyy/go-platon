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

package go_platon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func firstNonEmptyString(values ...*string) *string {
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}
		return &trimmed
	}
	return nil
}

func decodeMaybeBase64JSON(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return trimmed
	}

	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}

	for _, encoding := range encodings {
		decoded, err := encoding.DecodeString(trimmed)
		if err != nil {
			continue
		}
		decoded = []byte(strings.TrimSpace(string(decoded)))
		if json.Valid(decoded) {
			return string(decoded)
		}
	}

	return trimmed
}

func unwrapJSONString(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return trimmed
	}

	var unwrapped string
	if err := json.Unmarshal([]byte(trimmed), &unwrapped); err == nil {
		unwrapped = strings.TrimSpace(unwrapped)
		if unwrapped != "" {
			return unwrapped
		}
	}

	return trimmed
}

func normalizeApplePayPayload(raw string) (string, error) {
	normalized := unwrapJSONString(decodeMaybeBase64JSON(raw))
	if !json.Valid([]byte(normalized)) {
		return "", fmt.Errorf("apple pay payload must be valid JSON or base64-encoded JSON")
	}

	var payment struct {
		Token json.RawMessage `json:"token"`
	}
	if err := json.Unmarshal([]byte(normalized), &payment); err == nil {
		token := strings.TrimSpace(string(payment.Token))
		if token != "" && token != "null" {
			return token, nil
		}
	}

	return normalized, nil
}

func normalizeGooglePayPayload(raw string) (string, error) {
	normalized := unwrapJSONString(decodeMaybeBase64JSON(raw))
	if !json.Valid([]byte(normalized)) {
		return "", fmt.Errorf("google pay payload must be valid JSON or base64-encoded JSON")
	}

	var paymentData struct {
		PaymentMethodData struct {
			TokenizationData struct {
				Token string `json:"token"`
			} `json:"tokenizationData"`
		} `json:"paymentMethodData"`
	}
	if err := json.Unmarshal([]byte(normalized), &paymentData); err == nil {
		token := strings.TrimSpace(paymentData.PaymentMethodData.TokenizationData.Token)
		if token != "" {
			return normalizeGooglePayTokenString(token)
		}
	}

	return normalizeGooglePayTokenString(normalized)
}

func normalizeGooglePayTokenString(raw string) (string, error) {
	token := unwrapJSONString(strings.TrimSpace(raw))
	if token == "" {
		return "", fmt.Errorf("google pay tokenizationData.token is empty")
	}
	if json.Valid([]byte(token)) {
		return token, nil
	}

	unescaped, err := strconv.Unquote(fmt.Sprintf("%q", token))
	if err == nil {
		unescaped = unwrapJSONString(strings.TrimSpace(unescaped))
		if json.Valid([]byte(unescaped)) {
			return unescaped, nil
		}
	}

	// Keep the original token for backward compatibility with existing integrations
	// that already pass the exact string returned by Google Pay / frontend code.
	return token, nil
}

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

package http

import (
	"mime"
	"net/url"
	"sort"
	"strings"
)

const FormURLEncodedContentType = "application/x-www-form-urlencoded"

func IsFormURLEncodedContentType(contentType string) bool {
	trimmed := strings.TrimSpace(contentType)
	if trimmed == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(trimmed)
	if err != nil {
		return strings.EqualFold(trimmed, FormURLEncodedContentType)
	}

	return strings.EqualFold(mediaType, FormURLEncodedContentType)
}

// PrettyPrintFormURLEncodedBody formats a URL-encoded body
func PrettyPrintFormURLEncodedBody(raw string) string {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return raw
	}
	if len(values) == 0 {
		return "<empty>"
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	firstLine := true

	for _, key := range keys {
		fieldValues := values[key]
		if len(fieldValues) == 0 {
			if !firstLine {
				b.WriteByte('\n')
			}
			firstLine = false
			b.WriteString(key)
			b.WriteString("=<empty>")
			continue
		}

		for _, value := range fieldValues {
			if !firstLine {
				b.WriteByte('\n')
			}
			firstLine = false

			b.WriteString(key)
			b.WriteByte('=')
			if value == "" {
				b.WriteString("<empty>")
				continue
			}
			b.WriteString(value)
		}
	}

	if b.Len() == 0 {
		return "<empty>"
	}

	return b.String()
}

// FormatBodyForDebug pretty-prints URL-encoded bodies and returns
// raw payload for all other content types.
func FormatBodyForDebug(contentType string, raw []byte) string {
	if len(raw) == 0 {
		return "<empty>"
	}

	text := string(raw)
	if IsFormURLEncodedContentType(contentType) {
		return PrettyPrintFormURLEncodedBody(text)
	}

	return text
}

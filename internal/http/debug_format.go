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

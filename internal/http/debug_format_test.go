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
	"strings"
	"testing"
)

func TestIsFormURLEncodedContentType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "exact",
			contentType: "application/x-www-form-urlencoded",
			want:        true,
		},
		{
			name:        "with parameters",
			contentType: "application/x-www-form-urlencoded; charset=utf-8",
			want:        true,
		},
		{
			name:        "different case",
			contentType: "Application/X-WWW-Form-Urlencoded",
			want:        true,
		},
		{
			name:        "json",
			contentType: "application/json",
			want:        false,
		},
		{
			name:        "empty",
			contentType: "",
			want:        false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()
				if got := IsFormURLEncodedContentType(tc.contentType); got != tc.want {
					t.Fatalf("IsFormURLEncodedContentType(%q) = %v, want %v", tc.contentType, got, tc.want)
				}
			},
		)
	}
}

func TestPrettyPrintFormURLEncodedBody(t *testing.T) {
	t.Parallel()

	raw := "z=last&a=first+value&a=second&empty="
	got := PrettyPrintFormURLEncodedBody(raw)

	lines := strings.Split(got, "\n")
	want := []string{
		"a=first value",
		"a=second",
		"empty=<empty>",
		"z=last",
	}
	if len(lines) != len(want) {
		t.Fatalf("unexpected line count: got %d (%q), want %d", len(lines), got, len(want))
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d mismatch: got %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestFormatBodyForDebug(t *testing.T) {
	t.Parallel()

	t.Run(
		"form-urlencoded", func(t *testing.T) {
			t.Parallel()

			body := []byte("b=2&a=hello+world")
			got := FormatBodyForDebug("application/x-www-form-urlencoded; charset=utf-8", body)
			want := "a=hello world\nb=2"
			if got != want {
				t.Fatalf("FormatBodyForDebug(form) = %q, want %q", got, want)
			}
		},
	)

	t.Run(
		"non-form", func(t *testing.T) {
			t.Parallel()

			body := []byte("{\"ok\":true}")
			got := FormatBodyForDebug("application/json", body)
			if got != string(body) {
				t.Fatalf("FormatBodyForDebug(non-form) = %q, want %q", got, string(body))
			}
		},
	)
}

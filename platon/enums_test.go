package platon

import "testing"

func TestLang_String_NilReceiver(t *testing.T) {
	var lang *Lang

	if got := lang.String(); got != "" {
		t.Fatalf("String() mismatch: want empty string, got %q", got)
	}
}

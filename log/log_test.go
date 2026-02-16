package log

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestAll_UsesDebugThreshold(t *testing.T) {
	previousLevel := getLogLevel()
	t.Cleanup(func() {
		SetLevel(previousLevel)
	})

	logger := NewLogger("test ")

	SetLevel(LevelInfo)
	infoOutput := captureStderr(t, func() {
		logger.All("all-message")
	})
	if infoOutput != "" {
		t.Fatalf("expected no output for LevelInfo, got %q", infoOutput)
	}

	SetLevel(LevelDebug)
	debugOutput := captureStderr(t, func() {
		logger.All("all-message")
	})
	if !strings.Contains(debugOutput, "all-message") {
		t.Fatalf("expected output to contain message, got %q", debugOutput)
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	originalStderr := os.Stderr
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}

	os.Stderr = writePipe
	defer func() {
		os.Stderr = originalStderr
	}()

	fn()

	if err := writePipe.Close(); err != nil {
		t.Fatalf("writePipe.Close() error: %v", err)
	}

	var output bytes.Buffer
	if _, err := io.Copy(&output, readPipe); err != nil {
		t.Fatalf("io.Copy() error: %v", err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatalf("readPipe.Close() error: %v", err)
	}

	return output.String()
}

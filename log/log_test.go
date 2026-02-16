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
	t.Cleanup(
		func() {
			SetLevel(previousLevel)
		},
	)

	logger := NewLogger("test ")

	SetLevel(LevelInfo)
	infoOutput := captureStderr(
		t, func() {
			logger.All("all-message")
		},
	)
	if infoOutput != "" {
		t.Fatalf("expected no output for LevelInfo, got %q", infoOutput)
	}

	SetLevel(LevelDebug)
	debugOutput := captureStderr(
		t, func() {
			logger.All("all-message")
		},
	)
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

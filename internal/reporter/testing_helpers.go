package reporter

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// captureStdout runs fn while capturing stdout, returning the captured output.
// It safely restores stdout even if fn panics.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w

	// Ensure cleanup happens even if fn panics
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	func() {
		defer func() {
			w.Close()
			os.Stdout = old
		}()
		fn()
	}()

	return <-outC
}

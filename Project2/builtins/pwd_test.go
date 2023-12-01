package builtins_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Aryan-Patel5475/CSCE4600/Project2/builtins"
)

func TestPrintWorkingDirectory(t *testing.T) {
	expectedDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get working directory: %v", err)
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Error creating pipe: %v", err)
	}
	os.Stdout = w

	err = builtins.PrintWorkingDirectory()
	if err != nil {
		t.Fatalf("PrintWorkingDirectory() unexpected error: %v", err)
	}

	// Close the write end of the pipe before reading from the read end
	if err := w.Close(); err != nil {
		t.Fatalf("Error closing write pipe: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Error copying from pipe: %v", err)
	}

	output := buf.String()
	if output != expectedDir+"\\n" {
		t.Errorf("PrintWorkingDirectory() = %v, want %v", output, expectedDir)
	}
}

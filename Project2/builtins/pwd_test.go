package builtins_test

import (
	"bytes"
	"github.com/Aryan-Patel5475/CSCE4600/Project2/builtins"
	"os"
	"testing"
)

func TestPrintWorkingDirectory(t *testing.T) {
	// setup
	w := &bytes.Buffer{}

	// testing
	err := builtins.PrintWorkingDirectory(w)
	if err != nil {
		t.Fatalf("PrintWorkingDirectory() unexpected error: %v", err)
	}

	// Check if the output is the current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get working dir")
	}
	if gotWd := w.String(); gotWd != wd+"\n" {
		t.Errorf("PrintWorkingDirectory() = %v, want %v", gotWd, wd)
	}
}

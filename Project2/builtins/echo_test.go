package builtins_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Aryan-Patel5475/CSCE4600/Project2/builtins"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "single word",
			args:     []string{"hello"},
			expected: "hello\n",
		},
		{
			name:     "multiple words",
			args:     []string{"hello", "world"},
			expected: "hello world\n",
		},
		{
			name:     "empty",
			args:     []string{},
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Error creating pipe: %v", err)
			}
			os.Stdout = w

			builtins.Echo(tt.args...)

			_ = w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				t.Fatalf("Error copying from pipe: %v", err)
			}

			output := buf.String()
			if output != tt.expected {
				t.Errorf("Echo() = %v, want %v", output, tt.expected)
			}
		})
	}
}

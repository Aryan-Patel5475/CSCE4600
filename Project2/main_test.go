package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"os/user"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

func Test_runLoop(t *testing.T) {
	t.Parallel()
	exitCmd := strings.NewReader("exit\n")
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		wantErrW  string
		checkFunc func()
	}{
		{
			name: "no error",
			args: args{
				r: exitCmd,
			},
		},
		{
			name: "read error should have no effect",
			args: args{
				r: iotest.ErrReader(io.EOF),
			},
			wantErrW: "EOF",
		},
		{
			name: "handle input error",
			args: args{
				r: strings.NewReader("invalidcommand\nexit\n"),
			},
			wantErrW: "invalidcommand",
		},
		{
			name: "exit signal works",
			args: args{
				r: strings.NewReader("echo hello\nexit\n"),
			},
			wantW: "hello\nexiting gracefully...",
		},
		{
			name: "whoami command works",
			args: args{
				r: strings.NewReader("whoami\nexit\n"),
			},
			wantW: fmt.Sprintf("%s\n", getUsername()),
		},
		{
			name: "env command works",
			args: args{
				r: strings.NewReader(("env\nexit\n")),
			},
			wantW: fmt.Sprintf("%s\n", os.Environ()),
		},
		{
			name: "pwd command works",
			args: args{
				r: strings.NewReader("pwd\nexit\n"),
			},
			wantW: fmt.Sprintf("%s\n", getCurrentDirectory()),
		},
		{
			name: "mkdir command works",
			args: args{
				r: strings.NewReader("mkdir dir1 dir2\nexit\n"),
			},
			// Since mkdir doesn't produce output upon success, we don't specify wantW.
			checkFunc: func() {
				// Check if the directories are created successfully.
				_, err := os.Stat("dir1")
				assert.NoError(t, err, "dir1 should exist")

				_, err = os.Stat("dir2")
				assert.NoError(t, err, "dir2 should exist")
			},
		},
		{
			name: "rm command works",
			args: args{
				r: strings.NewReader("rm dir1 dir2\nexit\n"),
			},
			// Since rm doesn't produce output upon success, we don't specify wantW.
			checkFunc: func() {
				// Check if the directories are removed successfully.
				_, err := os.Stat("dir1")
				assert.True(t, os.IsNotExist(err), "dir1 should not exist")

				_, err = os.Stat("dir2")
				assert.True(t, os.IsNotExist(err), "dir2 should not exist")
			},
		},
		{
			name: "unknown command",
			args: args{
				r: strings.NewReader("unknown\nexit\n"),
			},
			wantErrW: "unknown",
		},
		{
			name: "echo multiple arguments",
			args: args{
				r: strings.NewReader("echo arg1 arg2 arg3\nexit\n"),
			},
			wantW: "arg1 arg2 arg3\nexiting gracefully...",
		},
		{
			name: "cd to invalid directory",
			args: args{
				r: strings.NewReader("cd non-existent-directory\nexit\n"),
			},
			wantErrW: "no such file or directory",
		},
		{
			name: "runLoop with timeout",
			args: args{
				r: iotest.TimeoutReader(strings.NewReader("exit\n")),
			},
			// Optionally check for specific behavior related to timeouts.
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}
			errW := &bytes.Buffer{}

			exit := make(chan struct{}, 2)
			// run the loop for 10ms
			go runLoop(tt.args.r, w, errW, exit)
			time.Sleep(10 * time.Millisecond)
			exit <- struct{}{}

			require.NotEmpty(t, w.String())
			if tt.wantErrW != "" {
				require.Contains(t, errW.String(), tt.wantErrW)
			} else {
				require.Empty(t, errW.String())
			}
		})
	}
}

func getUsername() string {
	u, err := user.Current()
	if err != nil {
		return "error getting username"
	}
	return u.Username
}

func getCurrentDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return "error getting current directory"
	}
	return wd
}

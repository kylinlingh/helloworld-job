package utils

import (
	"bytes"
	"fmt"
	"helloworld/pkg/errors"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// from: /Users/guanhao.lin/.gvm/gos/go1.18.10/src/cmd/go/internal/modfetch/codehost/codehost.go

var dirLock sync.Map

// Run runs the command line in the given directory
// (an empty dir means the current directory).
// It returns the standard output and, for a non-zero exit,
// a *RunError indicating the command, exit status, and standard error.
// Standard error is unavailable for commands that exit successfully.
// Example:
// out, gitErr := Run(r.dir, "git", "ls-remote", "-q", r.remote)
//
//	if gitErr != nil {
//				if rerr, ok := gitErr.(*RunError); ok {
//					if bytes.Contains(rerr.Stderr, []byte("fatal: could not read Username")) {
//						rerr.HelpText = "Confirm the import path was entered correctly.\nIf this is a private repository, see https://golang.org/doc/faq#git_https for additional information."
//					}
//				}
func Run(dir string, cmdline ...any) ([]byte, error) {
	return RunWithStdin(dir, nil, cmdline...)
}

func RunWithStdin(dir string, stdin io.Reader, cmdline ...any) ([]byte, error) {
	if dir != "" {
		muIface, ok := dirLock.Load(dir)
		if !ok {
			muIface, _ = dirLock.LoadOrStore(dir, new(sync.Mutex))
		}
		mu := muIface.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
	}
	// TODO: Impose limits on command output size.
	// TODO: Set environment to get English error messages.
	cmd := StringList(cmdline)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = dir
	c.Stdin = stdin
	c.Stderr = &stderr
	c.Stdout = &stdout
	err := c.Run()
	if err != nil {
		err = &errors.RunError{Cmd: strings.Join(cmd, " ") + " in " + dir, Stderr: stderr.Bytes(), Err: err}
	}
	return stdout.Bytes(), err
}

// StringList flattens its arguments into a single []string.
// Each argument in args must have type string or []string.
func StringList(args ...any) []string {
	var x []string
	for _, arg := range args {
		switch arg := arg.(type) {
		case []string:
			x = append(x, arg...)
		case string:
			x = append(x, arg)
		default:
			panic("stringList: invalid argument of type " + fmt.Sprintf("%T", arg))
		}
	}
	return x
}

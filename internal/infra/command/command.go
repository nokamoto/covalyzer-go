//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock_test.go -package=$GOPACKAGE
package command

import (
	"bytes"
	"log/slog"
	"os/exec"
)

type option func(*exec.Cmd)

type command struct{}

type runner interface {
	run(cmd string, args []string, opts ...option) error
	runO(cmd string, args []string, opts ...option) (*bytes.Buffer, error)
}

func (*command) run(cmd string, args []string, opts ...option) error {
	return run(cmd, args...)(opts...)
}

func (*command) runO(cmd string, args []string, opts ...option) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	opts = append(opts, withStdout(&buf))
	if err := run(cmd, args...)(opts...); err != nil {
		return &buf, err
	}
	return &buf, nil
}

func withDir(dir string) option {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}

func withStdout(buf *bytes.Buffer) option {
	return func(c *exec.Cmd) {
		c.Stdout = newLogWriter(buf)
	}
}

func run(cmd string, args ...string) func(opts ...option) error {
	return func(opts ...option) error {
		c := exec.Command(cmd, args...)
		c.Stdout = newLogWriter(nil)
		c.Stderr = newErrorLogWriter()
		for _, opt := range opts {
			opt(c)
		}
		if err := c.Run(); err != nil {
			slog.Debug("failed to run", "cmd", cmd, "args", args, "error", err)
			return err
		}
		return nil
	}
}

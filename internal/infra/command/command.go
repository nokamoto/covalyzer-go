package command

import (
	"bytes"
	"log/slog"
	"os/exec"
)

type Option func(*exec.Cmd)

// WithDir sets the working directory for the command.
func WithDir(dir string) Option {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}

// WithStdout sets the stdout for the command to capture the output to a buffer.
func WithStdout(buf *bytes.Buffer) Option {
	return func(c *exec.Cmd) {
		c.Stdout = newLogWriter(buf)
	}
}

func Run(cmd string, args ...string) func(opts ...Option) error {
	return func(opts ...Option) error {
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

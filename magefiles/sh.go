package main

import "github.com/magefile/mage/sh"

type shell struct {
	cmd  string
	args []string
}

type shells []shell

func do(cmd string, args ...string) shells {
	return shells{{cmd, args}}
}

func (ss shells) then(cmd string, args ...string) shells {
	return append(ss, do(cmd, args...)...)
}

func (ss shells) run() error {
	for _, s := range ss {
		if err := sh.Run(s.cmd, s.args...); err != nil {
			return err
		}
	}
	return nil
}

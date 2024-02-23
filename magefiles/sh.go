package main

import (
	"github.com/magefile/mage/sh"
)

type shell struct {
	cmd  string
	args []string
	v    bool
}

type shells []shell

func do(cmd string, args ...string) shells {
	return shells{{cmd, args, false}}
}

func doV(cmd string, args ...string) shells {
	return shells{{cmd, args, true}}
}

func (ss shells) then(cmd string, args ...string) shells {
	return append(ss, do(cmd, args...)...)
}

func (ss shells) thenV(cmd string, args ...string) shells {
	return append(ss, doV(cmd, args...)...)
}

func (ss shells) run() error {
	for _, s := range ss {
		f := sh.Run
		if s.v {
			f = sh.RunV
		}
		if err := f(s.cmd, s.args...); err != nil {
			return err
		}
	}
	return nil
}

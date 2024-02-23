package main

import (
	"github.com/magefile/mage/sh"
)

type shell struct {
	cmd  string
	args []string
	v    bool
	env  map[string]string
}

type shells []shell

func do(cmd string, args ...string) shells {
	return shells{{cmd, args, false, nil}}
}

func doV(cmd string, args ...string) shells {
	return shells{{cmd, args, true, nil}}
}

func doWith(env map[string]string, cmd string, args ...string) shells {
	return shells{{cmd, args, true, env}}
}

func (ss shells) then(cmd string, args ...string) shells {
	return append(ss, do(cmd, args...)...)
}

func (ss shells) thenV(cmd string, args ...string) shells {
	return append(ss, doV(cmd, args...)...)
}

func (ss shells) thenWith(env map[string]string, cmd string, args ...string) shells {
	return append(ss, doWith(env, cmd, args...)...)
}

func (ss shells) run() error {
	for _, s := range ss {
		f := sh.Run
		if s.v {
			f = sh.RunV
		}
		if s.env != nil {
			f = func(cmd string, args ...string) error {
				return sh.RunWith(s.env, cmd, args...)
			}
		}
		if err := f(s.cmd, s.args...); err != nil {
			return err
		}
	}
	return nil
}

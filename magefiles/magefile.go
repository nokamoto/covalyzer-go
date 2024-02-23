//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Build() error {
	if err := sh.Run("go", "install", "golang.org/x/tools/cmd/goimports@latest"); err != nil {
		return err
	}
	if err := sh.Run("goimports", "-w", "."); err != nil {
		return err
	}
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "mod", "tidy")
}

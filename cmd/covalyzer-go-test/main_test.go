package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/nokamoto/covalyzer-go/test"
)

func TestCovalyzerGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CovalyzerGo Suite")
}

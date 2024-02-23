package test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("covalyzer-go", func() {
	It("should be run with config.yaml", func() {
		cmd := exec.Command("covalyzer-go")
		cmd.Stdout = GinkgoWriter
		cmd.Stderr = GinkgoWriter
		cmd.Env = append(os.Environ(), "CONFIG_YAML=../../config.yaml", "DEBUG=true")
		Expect(cmd.Run()).ShouldNot(HaveOccurred())
	})

	It("should be success", func() {
		Expect(1).To(Equal(1))
	})

	It("should be failed", Label("fail"), func() {
		Fail("failed")
	})
})

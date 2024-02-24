package main

import (
	"slices"
	"strings"

	"github.com/nokamoto/covalyzer-go/internal/util/xslices"
)

func Build() error {
	list, err := golist()
	if err != nil {
		return err
	}

	ginkgoExcluded := slices.DeleteFunc(list, func(s string) bool {
		return strings.Contains(s, "covalyzer-go-test")
	})

	return do("go", "install", "golang.org/x/tools/cmd/goimports@latest").
		then("goimports", "-w", ".").
		then("go", "install", "github.com/bufbuild/buf/cmd/buf@v1.29.0").
		then("go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0").
		then("buf", "format", "-w").
		then("buf", "generate").
		then("go", "install", "go.uber.org/mock/mockgen@latest").
		then("go", "generate", "./...").
		then("go", "mod", "download").
		thenV("go", xslices.Concat("test", "-coverprofile=coverage.out", ginkgoExcluded)...).
		then("go", "mod", "tidy").
		run()
}

func Install() error {
	return do("go", "install", "./cmd/covalyzer-go").
		then("go", "install", "github.com/onsi/ginkgo/v2/ginkgo@latest").
		thenWith(map[string]string{"DEBUG": "1"}, "covalyzer-go").
		run()
}

func E2e() error {
	list, err := golist()
	if err != nil {
		return err
	}

	ginkgoOnly := slices.DeleteFunc(list, func(s string) bool {
		return !strings.Contains(s, "covalyzer-go-test")
	})

	return do("go", "install", "./cmd/covalyzer-go").
		thenV("go", xslices.Concat("test", ginkgoOnly, "--ginkgo.label-filter=!fail", "-v")...).
		run()
}

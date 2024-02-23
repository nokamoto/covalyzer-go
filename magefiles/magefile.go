package main

func Build() error {
	return do("go", "install", "golang.org/x/tools/cmd/goimports@latest").
		then("goimports", "-w", ".").
		then("go", "install", "github.com/bufbuild/buf/cmd/buf@v1.29.0").
		then("go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0").
		then("buf", "format", "-w").
		then("buf", "generate").
		then("go", "install", "go.uber.org/mock/mockgen@latest").
		then("go", "generate", "./...").
		then("go", "mod", "download").
		thenV("go", "test", "./...").
		then("go", "mod", "tidy").
		run()
}

func Install() error {
	return do("go", "install", "./cmd/covalyzer-go").
		thenWith(map[string]string{"DEBUG": "1"}, "covalyzer-go").
		run()
}

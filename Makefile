all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64

build-darwin-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/changelog-cli_darwin_amd64 cmd/changelog-cli/main.go

build-darwin-arm64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/changelog-cli_darwin_arm64 cmd/changelog-cli/main.go

build-linux-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/changelog-cli_linux_amd64 cmd/changelog-cli/main.go

build-linux-arm64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/changelog-cli_linux_arm64 cmd/changelog-cli/main.go

build-windows-amd64:
	mkdir -p build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/changelog-cli_windows_amd64.exe cmd/changelog-cli/main.go

build-docker-image:
	docker buildx build --platform linux/amd64,linux/arm64 -t sisimomo/changelog-cli . --push

gox-linux:
	gox -osarch="linux/amd64 linux/arm64" -output="build/changelog-cli_{{.OS}}_{{.Arch}}" ./cmd/changelog-cli

gox-all:
	gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" -output="build/changelog-cli_{{.OS}}_{{.Arch}}" ./cmd/changelog-cli

clean:
	rm -f build/changelog-cli_*
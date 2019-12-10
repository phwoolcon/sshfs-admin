.PHONY: build clean

build:
	@go version
	GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -v \
		-ldflags "-s -w -X main.GinMode=release -X main.Version=v$(shell date +%y.%-m.%-d)" \
		-o build/sshfs-admin-linux-x64 cmd/web/main.go
	upx -9 build/sshfs-admin-linux-x64
clean:
	rm -rf build

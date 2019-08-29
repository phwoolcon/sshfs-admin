build:
	GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -v \
		-ldflags '-s -w -X main.GinMode=release' -o build/sshfs-admin-linux-x64 cmd/web/main.go
	upx -9 build/sshfs-admin-linux-x64
clean:
	rm -rf build

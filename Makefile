all: clean build linux-amd64 linux-arm64 mac-amd64 mac-arm64

clean:
	rm -rf build

build:
	mkdir build

linux-amd64:
	mkdir build/linuxamd64
	env GOOS=linux GOARCH=amd64 go build -o build/linuxamd64/glasnostr main.go
	tar -czf build/linuxamd64/glasnostr-linux-amd64.tar.gz build/linuxamd64/glasnostr
linux-arm64:
	mkdir build/linuxarm64
	env GOOS=linux GOARCH=arm64 go build -o build/linuxarm64/glasnostr main.go
	tar -czf build/linuxarm64/glasnostr-linux-arm64.tar.gz build/linuxarm64/glasnostr
mac-amd64:
	mkdir build/darwinamd64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwinamd64/glasnostr main.go
	tar -czf build/darwinamd64/glasnostr-mac-amd64.tar.gz build/darwinamd64/glasnostr
mac-arm64:
	mkdir build/darwinarm64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwinarm64/glasnostr main.go
	tar -czf build/darwinarm64/glasnostr-mac-arm64.tar.gz build/darwinarm64/glasnostr
local: clean build
	go build -o build/glasnostr main.go
local-install: clean build
	go build -o build/glasnostr main.go
	sudo ln -sf ${PWD}/build/glasnostr /usr/local/bin/glasnostr

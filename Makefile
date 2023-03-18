all: clean build linux-amd64 linux-arm64 mac-amd64 mac-arm64

clean:
	rm -rf build

build:
	mkdir build

linux-amd64: clean build
	mkdir build/linuxamd64
	env GOOS=linux GOARCH=amd64 go build -o build/linuxamd64/glasnostr main.go
	cd build/linuxamd64; tar -czf glasnostr-linux-amd64.tar.gz glasnostr
linux-arm64: clean build
	mkdir build/linuxarm64
	env GOOS=linux GOARCH=arm64 go build -o build/linuxarm64/glasnostr main.go
	cd build/linuxarm64; tar -czf glasnostr-linux-arm64.tar.gz glasnostr
mac-amd64: clean build
	mkdir build/darwinamd64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwinamd64/glasnostr main.go
	cd build/darwinamd64; tar -czf glasnostr-mac-amd64.tar.gz glasnostr
mac-arm64: clean build
	mkdir build/darwinarm64
	env GOOS=darwin GOARCH=amd64 go build -o build/darwinarm64/glasnostr main.go
	cd build/darwinarm64; tar -czf glasnostr-mac-arm64.tar.gz glasnostr
local: clean build
	go build -o build/glasnostr main.go
local-install: clean build
	go build -o build/glasnostr main.go
	sudo ln -sf ${PWD}/build/glasnostr /usr/local/bin/glasnostr

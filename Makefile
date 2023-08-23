export CGO_ENABLED=1
publish = usr/publish
binarys = bin/ice.bin
version = src/version.go
binpack = src/binpack.go
flags = -ldflags "-w -s" -v

all: def
	@date +"%Y-%m-%d %H:%M:%S"
	go build ${flags} -o ${binarys} src/main.go ${version} ${binpack} && ./${binarys} forever restart &>/dev/null

def:
	@ [ -f src/version.go ] || echo "package main" > src/version.go
	@ [ -f src/binpack.go ] || echo "package main" > src/binpack.go

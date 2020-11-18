export GOPROXY=https://goproxy.cn
export GORPIVATE=github.com
export CGO_ENABLED=0

all:
	@echo && date
	go build -v -o bin/ice.bin src/main.go && chmod u+x bin/ice.bin && chmod u+x bin/ice.sh && ./bin/ice.sh restart

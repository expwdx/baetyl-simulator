APP:=bsctl
SRC_FILES:=$(shell find . -type f -name '*.go')

.PHONY: all
all: build

.PHONY: build
build: $(SRC_FILES)
	env GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 go build -o output/$(APP) .

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	@rm -rf output

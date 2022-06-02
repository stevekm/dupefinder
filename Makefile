SHELL:=/bin/bash
.ONESHELL:

# go version go1.17.10 darwin/amd64
# $ go mod init dupefinder
# # gofmt -l -w .
# # # $ go get github.com/google/go-cmp
# go run cmd/main.go .

test:
	set -euo pipefail
	go clean -testcache && \
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

docker-test:
	docker run --workdir $(CURDIR) -v $(CURDIR):$(CURDIR) --rm -ti golang:1.18-alpine ./test.sh

build:
	go build -o ./dupefinder cmd/main.go

# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
GIT_TAG:=$(shell git tag)
build-all:
	mkdir -p build ; \
	for os in darwin linux windows; do \
	for arch in amd64 arm64; do \
	output="build/dupefinder-$(GIT_TAG)-$$os-$$arch" ; \
	if [ "$${os}" == "windows" ]; then output="$${output}.exe"; fi ; \
	echo "building: $$output" ; \
	GOOS=$$os GOARCH=$$arch go build -o "$${output}" cmd/main.go ; \
	done ; \
	done

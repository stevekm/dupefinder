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

SHELL:=/bin/bash
.ONESHELL:

# go version go1.17.10 darwin/amd64
# $ go mod init dupefinder
# # gofmt -l -w .
# # # $ go get github.com/google/go-cmp
# go run cmd/main.go .

format:
	gofmt -l -w .

test:
	set -euo pipefail
	go clean -testcache && \
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

docker-test:
	docker run --workdir $(CURDIR) -v $(CURDIR):$(CURDIR) --rm -ti golang:1.18-alpine ./test.sh

build:
	go build -o ./dupefinder cmd/main.go
.PHONY:build
# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
GIT_TAG:=$(shell git describe --tags)
build-all:
	mkdir -p build ; \
	for os in darwin linux windows; do \
	for arch in amd64 arm64; do \
	output="build/dupefinder-v$(GIT_TAG)-$$os-$$arch" ; \
	if [ "$${os}" == "windows" ]; then output="$${output}.exe"; fi ; \
	echo "building: $$output" ; \
	GOOS=$$os GOARCH=$$arch go build -o "$${output}" cmd/main.go ; \
	done ; \
	done


build-server:
	go build -o ./server srv/server.go && chmod +x ./server

# $ curl "http://localhost:1000/?foo"
run-server:
	go run srv/server.go


# ~~~~~ Set up Benchmark dir ~~~~~ #
# set up a dir with tons of files and some very large duplicate files to test the program against

# https://go.dev/dl/go1.18.3.darwin-amd64.tar.gz
# https://dl.google.com/go/go1.18.3.darwin-amd64.tar.gz
BENCHDIR:=benchmarkdir
GO_TAR:=go1.18.3.darwin-amd64.tar.gz
$(GO_TAR):
	set -e
	wget https://dl.google.com/go/$(GO_TAR)

$(BENCHDIR): $(GO_TAR)
	set -e
	mkdir -p "$(BENCHDIR)"
	tar -C "$(BENCHDIR)" -xf "$(GO_TAR)"
	/bin/cp "$(GO_TAR)" $(BENCHDIR)
	/bin/cp "$(GO_TAR)" $(BENCHDIR)/go/
	/bin/cp "$(GO_TAR)" $(BENCHDIR)/copy2.tar.gz
	for i in $$(seq 1 5); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/bin/foo ; done
	for i in $$(seq 1 5); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/doc/foo2 ; done
	for i in $$(seq 1 10); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/lib/bar ; done
	for i in $$(seq 1 10); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/misc/bar2 ; done
	for i in $$(seq 1 15); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/src/baz ; done
	for i in $$(seq 1 15); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/test/baz2 ; done
	for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/src/buzz ; done
	for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/test/buzz2 ; done
	for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/bin/buzz3 ; done

benchmark-dir: $(BENCHDIR)

BIN:=./dupefinder
BENCHARGS:=--parallel 1
# takes about 10-60s for each iteration on standard NVMe SSD
# on SATA HDD each iteration should take about 2min30s - 3min~ish
benchmark: $(BENCHDIR) $(BIN)
	for i in md5 sha1 sha256 xxhash; do \
	echo ">>> ----- $$i ------" ; \
	for q in $$(seq 1 4 ); do \
	( set  -x ; time $(BIN) $(BENCHARGS) --algo $$i "$(BENCHDIR)" > /dev/null ; ) ; \
	done ; \
	done

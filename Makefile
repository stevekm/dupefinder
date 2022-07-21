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







# ~~~~~ Set up Benchmark dir ~~~~~ #
# set up a dir with tons of files and some very large duplicate files
# install conda to get a lot of files and dirs
# USAGE: $ ./dupefinder conda
# $ find conda/ -type f | wc -l
#    17953
BENCHDIR:=conda
CONDASH:=Miniconda3-py39_4.12.0-MacOSX-arm64.sh
# CONDASH=Miniconda3-py39_4.12.0-MacOSX-x86_64.sh
# CONDASH=Miniconda3-py39_4.12.0-Linux-x86_64.sh
CONDAURL:=https://repo.anaconda.com/miniconda/$(CONDASH)
$(CONDASH):
	wget "$(CONDAURL)"

$(BENCHDIR): $(CONDASH)
	@set -e
	bash "$(CONDASH)" -b -p $(BENCHDIR) || { rm -rf $(BENCHDIR); exit 1; }


# https://go.dev/dl/
go1.18.3.darwin-amd64.tar.gz:
	wget https://go.dev/dl/go1.18.3.darwin-amd64.tar.gz

big-dir-for-benchmarks: $(BENCHDIR) go1.18.3.darwin-amd64.tar.gz
	set -e
	/bin/cp go1.18.3.darwin-amd64.tar.gz $(BENCHDIR)/bin/
	/bin/cp go1.18.3.darwin-amd64.tar.gz $(BENCHDIR)/include/
	for i in $$(seq 1 15); do cat go1.18.3.darwin-amd64.tar.gz >> $(BENCHDIR)/etc/foo ; done
	/bin/cp $(BENCHDIR)/etc/foo $(BENCHDIR)/bin/bar

BIN:=./dupefinder
BENCHARGS:=--parallel 1
# takes about 2-20s for each iteration on standard NVMe SSD
benchmark: $(BENCHDIR) $(BIN)
	for i in md5 sha1 sha256 xxhash; do \
	echo ">>> ----- $$i ------" ; \
	for q in $$(seq 1 4 ); do \
	( set  -x ; time $(BIN) $(BENCHARGS) --algo $$i "$(BENCHDIR)" > /dev/null ; ) ; \
	done ; \
	done

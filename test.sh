#!/bin/sh
# test script for use inside of Docker container; otherwise use the Makefile
set -euo pipefail
go clean -testcache && \
go test -v ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''

MAKEFLAGS += \
	--warn-undefined-variables
SHELL = /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.SUFFIXES:

.PHONY: reality fast

reality:
	go run .
	@echo "Running normal simulation..."
	@echo "latency 1000ms and 50rps"

fast:
	@echo "Running fast simulation..."
	@echo "latency 100ms and 500rps"
	go run . \
  -l 100 \
  -d 100 \
  -w 12 \
  -s 15 \
  -e 0.001 \
  -a 5 \
  -r 500

#!/usr/bin/env bash

# Set up env variables
export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco

# Build executable
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/orinoco.go

# Start sieve server in the background
${PROJECT_DIR}/bin/orinoco pump -c ${CONFIG} &

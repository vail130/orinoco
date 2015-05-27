#!/usr/bin/env bash

TEST=1 ./bin/orinoco sieve &
sieve_pid=$!

go test ./...

kill $sieve_pid

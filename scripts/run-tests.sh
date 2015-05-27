#!/usr/bin/env bash

./bin/orinoco sieve &
sieve_pid=$!

go test ./...

kill $sieve_pid

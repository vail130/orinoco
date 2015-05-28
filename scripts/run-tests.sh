#!/usr/bin/env bash

TEST=1 ./bin/orinoco sieve &
sieve_pid=$!

rm -f pump.log
TEST=1 ./bin/orinoco pump -l pump.log &
pump_pid=$!

go test ./...

kill $pump_pid
kill $sieve_pid

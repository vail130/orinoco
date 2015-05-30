#!/usr/bin/env bash

export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

echo "" > ${PROJECT_DIR}/test.log
find ${PROJECT_DIR} -name "pump.log.*" -exec rm -f {} \;
	
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/orinoco.go

${PROJECT_DIR}/bin/orinoco sieve &
SIEVE_PID=$!

sleep 1

${PROJECT_DIR}/bin/orinoco pump -l ${PROJECT_DIR}/pump.log &
PUMP_PID=$!

sleep 1

cd ${PROJECT_DIR}
/usr/bin/go test ./...

kill $PUMP_PID $SIEVE_PID

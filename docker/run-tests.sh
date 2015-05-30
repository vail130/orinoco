#!/usr/bin/env bash

export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

echo "" > ${PROJECT_DIR}/test.log
find ${PROJECT_DIR} -name "pump.log.*" -exec rm -f {} \;
	
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/orinoco.go

${PROJECT_DIR}/bin/orinoco sieve &
SIEVE_PID=$!

COUNT=0
until nc -v -w 1 -z localhost 9966
do
    if [ "${COUNT}" -gt "3" ]
    then
        echo "Timed out!"
        exit 1
    fi
    COUNT=$((${COUNT} + 1))
    sleep 1
done

${PROJECT_DIR}/bin/orinoco pump -l ${PROJECT_DIR}/pump.log &
PUMP_PID=$!

cd ${PROJECT_DIR}
/usr/bin/go test github.com/vail130/orinoco/stringutils
/usr/bin/go test github.com/vail130/orinoco/sliceutils
/usr/bin/go test github.com/vail130/orinoco/httputils
/usr/bin/go test github.com/vail130/orinoco/sieve
/usr/bin/go test github.com/vail130/orinoco/pump
#/usr/bin/go test github.com/vail130/orinoco/tap
#/usr/bin/go test github.com/vail130/orinoco/litmus

kill $PUMP_PID $SIEVE_PID

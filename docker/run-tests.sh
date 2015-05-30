#!/usr/bin/env bash

export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

rm -f ${PROJECT_DIR}/pump.log
find ${PROJECT_DIR} -name "pump.log.*" -exec rm -f {} \;
rm -f ${PROJECT_DIR}/tap.log
	
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

${PROJECT_DIR}/bin/orinoco tap -l ${PROJECT_DIR}/tap.log &
TAP_PID=$!

${PROJECT_DIR}/bin/orinoco litmus -c ${PROJECT_DIR}/test-fixtures/litmus/test-litmus-config.yml &
LITMUS_PID=$!

cd ${PROJECT_DIR}
/usr/bin/go test github.com/vail130/orinoco/stringutils
/usr/bin/go test github.com/vail130/orinoco/sliceutils
/usr/bin/go test github.com/vail130/orinoco/httputils
/usr/bin/go test github.com/vail130/orinoco/sieve
/usr/bin/go test github.com/vail130/orinoco/pump
/usr/bin/go test github.com/vail130/orinoco/tap
/usr/bin/go test github.com/vail130/orinoco/litmus

kill $PUMP_PID $SIEVE_PID $TAP_PID $LITMUS_PID

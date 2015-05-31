#!/usr/bin/env bash

# Set up env variables
export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

# Remove artifacts of previous test runs
rm -f ${PROJECT_DIR}/pump.log
find ${PROJECT_DIR} -name "pump.log.*" -exec rm -f {} \;
rm -f ${PROJECT_DIR}/tap.log

# Build executable
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/orinoco.go

# Start sieve server in the background
${PROJECT_DIR}/bin/orinoco sieve &
SIEVE_PID=$!

# Wait for sieve to come up on http://localhost:9966
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

# Start tap service in the background
${PROJECT_DIR}/bin/orinoco tap -l ${PROJECT_DIR}/tap.log &
TAP_PID=$!

# Start pump service in the background
${PROJECT_DIR}/bin/orinoco pump -c ${PROJECT_DIR}/test-fixtures/test-pump-config.yml &
PUMP_PID=$!

# Start litmus service in the background
${PROJECT_DIR}/bin/orinoco litmus -c ${PROJECT_DIR}/test-fixtures/test-litmus-config.yml &
LITMUS_PID=$!

# Specify packages to test here
read -r -d '' PACKAGES << EOM
sliceutils
stringutils
httputils
sieve
tap
pump
litmus
EOM

# Run all packages or only one specified in environment.
# Allows testing one package.
cd ${PROJECT_DIR}
for pkg in $PACKAGES; do
	if [ -z ${TEST_PKG} ] || [ "${TEST_PKG}" == "${pkg}" ]; then
		/usr/bin/go test github.com/vail130/orinoco/${pkg}
	fi
done

# Kill child processes
kill $SIEVE_PID $TAP_PID $PUMP_PID $LITMUS_PID
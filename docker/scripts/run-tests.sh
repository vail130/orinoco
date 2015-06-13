#!/usr/bin/env bash

# Unpack arguments
TEST_TYPE=$1

# Set up env variables
export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

# Remove artifacts of previous test runs
mkdir -p ${PROJECT_DIR}/bin
rm -rf ${PROJECT_DIR}/bin/*
mkdir -p ${PROJECT_DIR}/artifacts
rm -rf ${PROJECT_DIR}/artifacts/*

# Build executable
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/main.go

# Start sieve server in the background
if [ "${TEST_TYPE}" == "s3" ]; then
	TEST_CONFIG=${PROJECT_DIR}/test-fixtures/test-s3-config.yml
else
	TEST_CONFIG=${PROJECT_DIR}/test-fixtures/test-config.yml
fi

${PROJECT_DIR}/bin/orinoco $TEST_CONFIG &
ORINOCO_PID=$!

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

# Specify packages to test here
#tap
#pump
read -r -d '' PACKAGES << EOM
sliceutils
stringutils
httputils
sieve
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
kill $ORINOCO_PID

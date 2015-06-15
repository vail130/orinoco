#!/usr/bin/env bash

# Set up env variables
export GOPATH=/go:/go/src/github.com/vail130/orinoco/Godeps/_workspace
export PROJECT_DIR=/go/src/github.com/vail130/orinoco
export TEST=1

# Remove artifacts of previous test runs
mkdir -p ${PROJECT_DIR}/bin
rm -rf ${PROJECT_DIR}/bin/*
mkdir -p ${PROJECT_DIR}/artifacts
rm -rf ${PROJECT_DIR}/artifacts/*

function wait_for_port() {
	COUNT=0
	until nc -v -w 1 -z localhost $1
	do
	    if [ "${COUNT}" -gt "3" ]
	    then
	        echo "Timed out!"
	        exit 1
	    fi
	    COUNT=$((${COUNT} + 1))
	    sleep 1
	done
}

if [ "${TEST_SVC}" == "http" ]; then
	/usr/local/bin/reflect --port 8080 &> ${PROJECT_DIR}/artifacts/reflect.log &
	REFLECT_PID=$!
	wait_for_port 8080
fi

# Build executable
/usr/bin/go build -o ${PROJECT_DIR}/bin/orinoco ${PROJECT_DIR}/main.go

TEST_CONFIG="${PROJECT_DIR}/test-fixtures/test-${TEST_SVC}-config.yml"
${PROJECT_DIR}/bin/orinoco $TEST_CONFIG &
ORINOCO_PID=$!
wait_for_port 9966

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

if [ "${TEST_SVC}" == "http" ]; then
	kill $REFLECT_PID
fi

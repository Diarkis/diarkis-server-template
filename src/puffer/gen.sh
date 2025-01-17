#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

PROJECT_NAME=$(go mod edit -json | jq -r .Module.Path)/puffer/go

PUFFER_BIN=
if [ $(uname) == 'Darwin' ]; then
    PUFFER_BIN=./puffer-mac
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    PUFFER_BIN=./puffer-linux
else
    echo "unsupported platform"
    exit 1
fi

${PUFFER_BIN} . . ${PROJECT_NAME}/puffer/go
go fmt ./go/...

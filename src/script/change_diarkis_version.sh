#!/bin/bash

# This script is used to change the version of Diarkis in the project.
# download diarkis codere from help.diarkis.io and changet the version in the go.mod file

# Check if the script is being run from the root of the project
if [ ! -f go.mod ]; then
    echo "This script must be run from the root of the project"
    exit 1
fi

SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)

VERSION=$1
if [ $# -ne 1 ]; then
    # prompt the user to enter the version of Diarkis to use
    echo -n "Enter the version of Diarkis to use: "
    read -r version
    VERSION=${version}
fi

pushd ${SCRIPT_DIR}/..
    rm -rf "./coderefs/${VERSION}"
    # Check if the version is provided
    wget -qO- "https://docs.diarkis.io/sdk/server_coderef/${VERSION}.tar.gz" | tar xvzf - -C ./coderefs --strip-components=1
    go mod edit -require "github.com/Diarkis/diarkis@${VERSION}"
    go mod edit -replace "github.com/Diarkis/diarkis=./coderefs/${VERSION}"
popd

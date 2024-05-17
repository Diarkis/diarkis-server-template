#!/bin/bash

# This script is used to change the version of Diarkis in the project.
# download diarkis codere from help.diarkis.io and changet the version in the go.mod file

# Check if the script is being run from the root of the project
if [ ! -f go.mod ]; then
    echo "This script must be run from the root of the project"
    exit 1
fi
# prompt the user to enter the version of Diarkis to use
echo -n "Enter the version of Diarkis to use:"
read -r version
rm -rf "./coderefs/$version"
# Check if the version is provided
wget -qO- "https://docs.diarkis.io/sdk/server_coderef/$version.tar.gz" | tar xvzf - -C ./coderefs --strip-components=1
go mod edit -require "github.com/Diarkis/diarkis@$version"
go mod edit -replace "github.com/Diarkis/diarkis=./coderefs/$version"

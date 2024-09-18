#!/bin/bash

function usage {
    cat <<EOF
$(basename ${0}) is a generate diarkis project tool.

Usage:
    $(basename ${0}) projectID builderToken outputPath <moduleName optional>
Sample:
    $(basename ${0}) 012345678 sampleToken /tmp/sample-project
    or
    $(basename ${0}) 012345678 sampleToken /tmp/sample-project github.com/sample-origanization/sample-project

EOF
    exit 1
}

# Accept optional module name
if [ $# -ne 3 -a $# -ne 4 ]; then
    usage
fi

if [ $# -eq 4 ]; then
    module_name=$4
fi

project_id=$1
builder_token=$2
output_path=$3
if [ -z "${module_name}" ]; then
    module_name=$(basename $output_path)
fi

go run ./tools/install.go $project_id $builder_token $output_path
pushd $output_path
    go mod edit -module $module_name
    if [ $(uname) == 'Darwin' ]; then
        find . -type f -name '*.go'  -exec sed -i '' -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" {} \;
        sed -i '' -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" puffer/gen.sh
    elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
        find . -type f -name '*.go'  -exec sed -i -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" {} \;
        sed -i -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" puffer/gen.sh
        echo "Linux"
    else
        echo "Unsupported OS"
        uname -a
        exit 1
    fi
popd

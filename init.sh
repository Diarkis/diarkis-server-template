#!/bin/bash

function usage {
    cat <<EOF
$(basename ${0}) is a generate diarkis project tool.

Usage:
    $(basename ${0}) moduleName builderToken outputPath
Sample:
    $(basename ${0}) github.com/sample-origanization/sample-project sampleToken /tmp/sample-project

EOF
    exit 1
}

if [ $# -ne 3 ]; then
    usage
fi

project_id=$1
builder_token=$2
output_path=$3
module_name=$(basename $output_path)
go run ./tools/install.go $project_id $builder_token $output_path
cd $output_path &&
    go mod edit -module $module_name &&
    find . -type f -name '*.go'  -exec sed -i '' -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" {} \; &&
    sed -i '' -e "s%github.com/Diarkis/diarkis-server-template%$module_name%g" puffer/gen.sh
#!/bin/bash

function usage {
    cat <<EOF
$(basename ${0}) is a generate diarkis project tool.

Usage:
    $(basename ${0}) projectID builderToken outputPath <moduleName optional>
Sample:
    $(basename ${0}) 012345678 sampleToken /tmp/sample-project
    or
    $(basename ${0}) 012345678 sampleToken /tmp/sample-project github.com/sample-organization/sample-project

EOF
    exit 1
}

CURR_PWD=`pwd`

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

# check module name using go mod init
tmp_dir=

cleanup() {
    cd ${CURR_PWD}
    if [ ! -z "$tmp_dir" ]; then
        rm -fr "$tmp_dir"
    fi
}

trap cleanup EXIT

tmp_dir=`mktemp -d -q`
exit_code=$?
if [ $exit_code -ne 0 ]; then
    echo "fail to create temporary directory"
    exit 1
fi

pushd ${tmp_dir}
    go mod init ${module_name}
    if [ $? -ne 0 ]; then
        exit 1
    fi
popd

go run ./tools/install $project_id $builder_token $output_path
go run ./tools/rewrite_import.go $output_path "github.com/Diarkis/diarkis-server-template" "$module_name"
pushd $output_path
    go mod edit -module $module_name
    echo "PROJECT_NAME=$module_name" > puffer/vars.sh
popd

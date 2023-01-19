#!/bin/bash
## made by RHEMS Japan.Co., Ltd.
## usage

if [ $# != 3 ]; then
    echo "Usage: $0 DIARKIS_PROJECT_ID NAMESPACE GCP_ACCOUNT_ID"
    exit 1
fi

DIARKIS_PROJECT_ID=$1
NAMESPACE=$2
GCP_ACCOUNT_ID=$3

DIR_BASE=$(cd $(dirname ${0}); pwd)
CLIENT_DIR=${DIR_BASE}/../overlays/clients/${DIARKIS_PROJECT_ID}
echo ${CLIENT_DIR}

if [ -e ${CLIENT_DIR} ]; then
    echo "directory already exsits : ${CLIENT_DIR}"
    exit 1
else
    cp -Rip ${DIR_BASE}/../overlays/base/gcp ${CLIENT_DIR}
    ## SED NAMESPACE
    find ${CLIENT_DIR} -type f | xargs sed -i "" "s/__NAMESPACE__/${NAMESPACE}/g"
    ## SED PROJECT ID
    find ${CLIENT_DIR} -type f | xargs sed -i "" "s/__SAMPLE_PROJECT_ID__/\"${DIARKIS_PROJECT_ID}\"/g"
    ## SED AWS ACCOUNT
    find ${CLIENT_DIR} -type f | xargs sed -i "" "s/__GCP_ACCOUNT_ID__/${GCP_ACCOUNT_ID}/g"
fi

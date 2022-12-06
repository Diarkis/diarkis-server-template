#!/usr/bin/env bash
set -ex
DIARKIS_CLI_VERSION=0.0.1
DIARKIC_CLI_URL=https://diarkis.io/cli/v
file=
case $(uname -s) in
    Linux*)     file=diarkis-cli_"$DIARKIS_CLI_VERSION"_linux_amd64.tar.gz;;
    Darwin*)    file=diarkis-cli_"$DIARKIS_CLI_VERSION"_macOS_arm64.tar.gz;;
    *)          exit 1
esac

curl -LO "$DIARKIC_CLI_URL""$DIARKIS_CLI_VERSION"/$file
tar -xzvf $file
rm $file
mv "$(basename $file .tar.gz)" diarkis-cli

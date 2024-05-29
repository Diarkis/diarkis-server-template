#!/bin/bash
if [ $(uname) == 'Darwin' ]; then
    ./puffer-mac . . github.com/Diarkis/diarkis-server-template/puffer/go
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    ./puffer-linux . . github.com/Diarkis/diarkis-server-template/puffer/go
fi
go fmt ./go/...

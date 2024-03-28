#!/bin/bash
if [ $(uname) == 'Darwin' ]; then
    ./puffer-mac . . {0}
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    ./puffer-linux . . {0}
fi

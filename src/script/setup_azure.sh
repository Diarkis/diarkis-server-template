#!/bin/bash

read -r -p "Enter Your Azure Container Registry domain:(sample.azurecr.io) " acr

echo
echo "Your Azure Container Registry domain is: $acr"

if [ $(uname) == 'Darwin' ]; then
    find k8s -type f | xargs sed -i '' -e s/'__ACR_DOMAIN__'/"$acr"/g
    sed  -i '' -e s/'__ACR_DOMAIN__'/"$acr"/g Makefile
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    find k8s -type f  | xargs sed -i -e s/'__ACR_DOMAIN__'/"$acr"/g
    sed -i -e s/'__ACR_DOMAIN__'/"$acr"/g Makefile
else
    echo "Unsupported OS"
    uname -a
    exit 1
fi

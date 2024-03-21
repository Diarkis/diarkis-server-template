#!/bin/bash

read -r -p "Enter Your Azure Container Registry domain:(sample.azurecr.io) " acr

echo
echo "Your Azure Container Registry domain is: $acr"

find k8s -type f  | xargs sed -i '' -e s/'__ACR_DOMAIN__'/"$acr"/g
sed  -i '' -e s/'__ACR_DOMAIN__'/"$acr"/g Makefile

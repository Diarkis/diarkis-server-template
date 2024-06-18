#!/bin/bash

read -r -p "Enter Your GCP ProjectID: " projectid

echo
echo "Your GCP Project ID is: $projectid"

if [ $(uname) == 'Darwin' ]; then
    find k8s -type f | xargs sed -i '' -e s/'__GCP_PROJECT_ID__'/"$projectid"/g
    sed -i '' -e s/'__GCP_PROJECT_ID__'/"$projectid"/g Makefile
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    find k8s -type f  | xargs sed -i -e s/'__GCP_PROJECT_ID__'/"$projectid"/g
    sed -i -e s/'__GCP_PROJECT_ID__'/"$projectid"/g Makefile
else
    echo "Unsupported OS"
    uname -a
    exit 1
fi

#!/bin/bash

read -r -p "Enter Your GCP ProjectID: " projectid

echo
echo "Your GCP Project ID is: $projectid"

find k8s -type f  | xargs sed -i '' -e s/'__GCP_PROJECT_ID__'/"$projectid"/g
sed  -i '' -e s/'__GCP_PROJECT_ID__'/"$projectid"/g Makefile

#!/bin/bash

read -r -p "Enter Your AWS Account Number: (This can be get by commnd \"aws sts get-caller-identity | jq .Account -r\"): " aws_account_num


echo
echo "Your AWS Account Number is: $aws_account_num"

find k8s -type f  | xargs sed -i '' -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g
sed  -i '' -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g Makefile

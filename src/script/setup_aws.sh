#!/bin/bash

read -r -p "Enter Your AWS Account Number: (This can be get by commnd \"aws sts get-caller-identity | jq .Account -r\"): " aws_account_num


echo
echo "Your AWS Account Number is: $aws_account_num"

if [ $(uname) == 'Darwin' ]; then
    find k8s -type f | xargs sed -i '' -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g
    sed -i '' -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g Makefile
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
    find k8s -type f | xargs sed -i -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g
    sed -i -e s/'__AWS_ACCOUNT_NUM__'/"$aws_account_num"/g Makefile
else
    echo "Unsupported OS"
    uname -a
    exit 1
fi

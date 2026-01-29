#!/usr/bin/env sh

VALUE_NAME=
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 {dev|prd}"
    exit 1
fi
case "$1" in
    dev)
        VALUE_NAME=dev-values.yaml
        ;;
    prd)
        VALUE_NAME=prd-values.yaml
        ;;
    *)
        echo "Invalid Argument: $1"
        echo "Usage: $0 {dev|prd}"
        exit 1
        ;;
esac

helm install prometheus prometheus-community/prometheus -n prometheus  -f "$VALUE_NAME"

#!/usr/bin/env sh
#
# Install KEDA, the Kubernetes Event-driven Autoscaler.
# https://keda.sh/docs/

helm repo add kedacore https://kedacore.github.io/charts
helm repo update

helm upgrade --install keda kedacore/keda --namespace keda --create-namespace

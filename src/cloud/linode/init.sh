#!/bin/bash

set -e  # Exit on error

# Install Terraform
if ! command -v terraform &> /dev/null
then
    echo "Terraform not found. Installing..."
    curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
    sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
    sudo apt-get update && sudo apt-get install terraform
else
    echo "Terraform is already installed."
fi

# Install kubectl
if ! command -v kubectl &> /dev/null
then
    echo "kubectl not found. Installing..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
else
    echo "kubectl is already installed."
fi

# Install kustomize
if ! command -v kustomize &> /dev/null
then
    echo "kustomize not found. Installing..."
    curl -LO "https://github.com/kubernetes-sigs/kustomize/releases/latest/download/kustomize_linux_amd64.tar.gz"
    tar -xvzf kustomize_linux_amd64.tar.gz
    sudo mv kustomize /usr/local/bin/
    rm kustomize_linux_amd64.tar.gz
else
    echo "kustomize is already installed."
fi

# Install Helm
if ! command -v helm &> /dev/null
then
    echo "Helm not found. Installing..."
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
else
    echo "Helm is already installed."
fi

# Add Helm repositories
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts || echo "prometheus-community repo already added."
helm repo add grafana https://grafana.github.io/helm-charts || echo "grafana repo already added."

# Update Helm repositories
helm repo update

echo "Init script completed. All required commands are installed, and Helm repositories are added."

#!/bin/bash

set -e

usage() {
  echo "Usage: TF_VAR_token=<token> $0 [workspace] [ROOT_DIR]"
  echo "  workspace (optional): The Terraform workspace to use. (default: dev-diarkis-asia)"
  echo "  TF_VAR_token: The Linode Personal Access Token."
  echo
}

if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
  usage
  exit 0
fi

# Set workspace_name as the first argument or default to dev-diarkis-asia
workspace_name=${1:-dev-diarkis-asia}

# Set ROOT_DIR as the first argument or default to the root directory of the project
ROOT_DIR=$(
  cd $(dirname $0)
  cd ../../
  pwd
)

# Check if TF_VAR_token is set.
if [ -z "$TF_VAR_token" ]; then
  echo "Error: Token is required. Set TF_VAR_token before running."
  exit 1
fi

if [ -z "$LINODE_S3_ACCESS_KEY" ] || [ -z "$LINODE_S3_SECRET_KEY" ]; then
  echo "Error: LINODE_S3_ACCESS_KEY and LINODE_S3_SECRET_KEY are required. Set them before running."
  exit 1
fi

export KUBECONFIG=$ROOT_DIR/terraform/linode/kubeconfig

# Validate all required commands are installed
check_command() {
  if ! command -v "$1" &>/dev/null; then
    echo "Error: $1 is not installed. Please run 'init.sh' to install the required dependencies." >&2
    exit 1
  fi
}

echo "Validating required commands..."
check_command terraform
check_command kubectl
check_command kustomize
check_command helm
check_command curl
echo "All required commands are installed."

# Create the Kubernetes cluster using Terraform
echo "Initializing and applying Terraform configuration..."

# Create a new workspace if not already exists
terraform -chdir=$ROOT_DIR/terraform/linode init
(
  set +e
  terraform -chdir=$ROOT_DIR/terraform/linode workspace select "$workspace_name"
  if [ $? -ne 0 ]; then
    echo "Creating a new workspace: $workspace_name"
    terraform -chdir=$ROOT_DIR/terraform/linode workspace new "$workspace_name"
    terraform -chdir=$ROOT_DIR/terraform/linode workspace select "$workspace_name"
  fi
)
terraform -chdir=$ROOT_DIR/terraform/linode apply -var-file="workspaces/$workspace_name/terraform.tfvars" -auto-approve

# Get results from Terraform output
FIREWALL_ID=$(terraform -chdir=$ROOT_DIR/terraform/linode output firewall_id | tr -d '"')
echo "Firewall has been created with ID: $FIREWALL_ID"

# Extract kubeconfig and set for kubectl
terraform -chdir=$ROOT_DIR/terraform/linode output -raw kubeconfig | base64 -d >$KUBECONFIG
echo "KUBECONFIG is set to $KUBECONFIG"

(
  set +e
  until kubectl get nodes --no-headers 2>&1 | grep -qv "No resources found"; do
    echo "Waiting for nodes to be registered..."
    sleep 10
  done
)

echo "Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready node --all --timeout=600s

# Apply Kubernetes resources using Kustomize
echo "Applying Kustomize manifests..."
kustomize build $ROOT_DIR/k8s/linode/overlays/dev0 | kubectl apply -f -

echo "Waiting for all pods to be ready in the dev0 namespace..."
(
  set +e
  until kubectl get pods --no-headers -n dev0 2>&1 | grep -qv "No resources found"; do
    echo "Waiting for pods to be registered..."
    sleep 1
  done
)
kubectl wait --for=condition=Ready pods --all --namespace=dev0 --timeout=600s

# Apply Firewall configuration using Linode Firewall Operator
curl -s https://raw.githubusercontent.com/gangyi89/deploy-linode-operator/main/deploy-linode-fw-operator.sh | bash -s -- dev0
kubectl get secret linode-api-key -n dev0 || kubectl create secret generic linode-api-key --from-literal=api-key="$TF_VAR_token" -n dev0
sed s/"<MY_FIREWALL_ID>"/$FIREWALL_ID/ $ROOT_DIR/k8s/linode/cluster-firewall.yaml | kubectl apply -n dev0 -f -

# Install Prometheus using Helm
echo "Installing Prometheus using Helm..."
kubectl get namespace monitoring || kubectl create namespace monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --values=$ROOT_DIR/k8s/helm/kube-prometheus-stack/values.yaml --debug

echo "Waiting for Prometheus pods to be ready..."
kubectl wait --for=condition=Ready pods --all --namespace=monitoring --timeout=600s

# Apply ServiceMonitor configuration
echo "Applying ServiceMonitor configuration..."
kubectl apply -f $ROOT_DIR/k8s/servicemonitor.yaml -n monitoring

echo "You can access Grafana by running:"
echo "  kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring"
echo "Use admin/prom-operator as credentials."

# Install Loki using Helm
echo "Installing Loki using Helm..."
kubectl get namespace logging || kubectl create namespace logging
# kubectl get secret linode-s3-creds -n logging || kubectl create secret generic linode-s3-creds --from-literal=AWS_ACCESS_KEY_ID="$LINODE_S3_ACCESS_KEY" --from-literal=AWS_ACCESS_KEY_SECRET="$LINODE_S3_SECRET_KEY" -n logging
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki -n logging --values=$ROOT_DIR/k8s/helm/loki/values.yaml --set loki.storage_config.aws.access_key_id="$LINODE_S3_ACCESS_KEY" --set loki.storage_config.aws.secret_access_key="$LINODE_S3_SECRET_KEY" --set loki.storage.s3.accessKeyId="$LINODE_S3_ACCESS_KEY" --set loki.storage.s3.secretAccessKey="$LINODE_S3_SECRET_KEY" --debug
helm upgrade --install promtail grafana/promtail -n logging --set loki.serviceName=loki-gateway --set loki.servicePort=80 --debug

sleep 1

echo "Waiting for Loki pods to be ready..."
kubectl wait --for=condition=Ready pods --all --namespace=logging --timeout=600s

echo "Setup completed."

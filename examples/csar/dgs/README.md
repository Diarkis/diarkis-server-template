# Overview

This sample explains the following two procedures necessary for developing DGS:

- How to develop DGS locally
- How to deploy DGS to cloud environments

## How to Develop DGS

### 1. Build

#### To build on Linux or macOS

```sh
./run-mage.sh build:local
```

or

```sh
make build-local
```

#### To build on Windows

```
.\run-mage.bat build:local
```

### 2. Build DGS

Develop and build DGS on the client side.
For details, please refer to the client-side DGS development documentation.

### 3. Run mars, HTTP, UDP servers

#### Run servers on Linux or macOS

```sh
./run-mage.sh server mars
./run-mage.sh server http
./run-mage.sh server udp
```

or

```sh
./remote_bin/mars configs/mars/main.json
./remote_bin/http
./remote_bin/udp
```

#### Run servers on Windows

```
.\run-mage.bat server mars
.\run-mage.bat server http
.\run-mage.bat server udp
```

### 4. Run DGS

Execute the DGS built in step 2.
For details, please refer to the client-side DGS development documentation.

---

## How to Deploy DGS to Cloud

This section explains how to deploy DGS developed in Unity to cloud environments.

### 1. Setup Kubernetes Cluster

First, create a Kubernetes cluster.
Please set up according to your environment:

- AWS: https://help.diarkis.io/en/diarkis-server/setup-cloud/aws
- GCP: Under preparation

### 2. Create Container Registry for DGS

Create a container registry for DGS.

- **AWS ECR**

```
aws sts get-caller-identity # Please verify the target is correct
aws ecr create-repository --repository-name dgs
```

- **GCP Artifact Registry**
  If the diarkis repository is already created, no additional work is required.

### 3. Create DGS Binary

Build the DGS binary. For details, please refer to the client-side DGS development documentation.
The executable file name should be `linux-server`. The server template provided on this page assumes executing `linux-server.x86_64` and provides DGS configuration accordingly.

### 4. Copy DGS Binary to `dgs_bin` Directory

Copy the binary built earlier to the `dgs_bin` directory.

Please copy all files necessary to execute `linux-server.x86_64` under this directory.

Example file structure for Unity DGS:

```tree
dgs_bin
├── linux-server_Data
├── linux-server.x86_64 ... executable file
├── log.json ... log configuration file (optional: required for local execution)
├── mesh.json ... mesh configuration file (optional: required for local execution)
├── PLACE_DGS_BINARIES_HERE.md ... build instructions (can be deleted)
├── UnityPlayer.so
└── Other files required for execution (lib, etc.)
```

### 5. Build and Push DGS Image

Execute the following commands to build and push the DGS image.
(mars, http, udp will also be processed simultaneously.)

```sh
# Build DGS container
make build-container-aws
# Push DGS container
make push-container-aws
```

### 6. Install KEDA

DGS scaling uses [KEDA (Kubernetes Event-Driven Autoscaling)](https://keda.sh/). It can be easily installed using helm. Execute the following command while connected to the kubernetes cluster you want to apply:

```sh
# Install helm (if not installed)
#   https://helm.sh/docs/intro/install/
# For macOS using Homebrew, execute:
brew install helm

# Install KEDA
./cloud/aws/install-keda.sh

# Verify installation completion
kubectl get ns

# If the following result is output, installation is complete
NAME                STATUS   AGE
amazon-cloudwatch   Active   116m
default             Active   129m
keda                Active   8s
kube-node-lease     Active   129m
kube-public         Active   129m
kube-system         Active   129m
```

### 7. Deploy DGS

```sh
# Deploy to dev0 namespace
kustomize build k8s/aws/overlays/dev0 | kubectl apply -f -
```

### 8. Connectivity Check

Run the DGS client to access the DGS server and verify connectivity.
For details, please refer to the client-side DGS development documentation.

## DGS Scaling Configuration

When configuring DGS scaling, you need to consider the following:

- CPU resources used by DGS (can be specified in milli units)
- Maximum number of DGS pods (replica count)
  - Set this after calculating the required number based on expected maximum users or game sessions
- Maximum number of nodes in the public nodepool based on the above
- Consider that DGS and other resources consume several hundred milli CPU
  - For example, with 1 node of 4 cores, if DGS uses 1 CPU (1000m), up to 3 pods can be used

The values to set based on these considerations are:

> (Required maximum instance count) = (Maximum DGS pods) / (Available pods per instance)

| Node Core | DGS CPU limits | DGS pods / node | DGS max pod | required max nodes |
| :-------- | :------------- | :-------------- | :---------- | :----------------- |
| 4         | 1000           | 3               | 300         | 300/3 = 100        |
| 4         | 1000           | 3               | 1000        | 1000/3 = 334       |
| 4         | 800            | 4               | 1000        | 1000/4 = 250       |

You also need to consider the following:

- Number of standby DGS pods
  - Increasing this allows response to sudden increases in game sessions
  - However, be careful as this also leads to increased costs

### CPU Resource Configuration

The CPU resource count for DGS can be changed by editing the following file:
**k8s/aws/overlays/dev0/dgs/deploy.yaml**

- `spec.template.spec.containers[].resources.limits.cpu` : Maximum CPU resource value
  - The pod cannot use more CPU resources than this
- `spec.template.spec.containers[].resources.requests.cpu` : Minimum CPU resource value
  - Minimum CPU resources to reserve

At pod startup, `requests` CPU resources are reserved, but up to `limits` can be used depending on load. However, if the node doesn't have sufficient resources, `limits` may not be fully available.

Therefore, especially for DGS, it's recommended to set `requests = limits`.

The following shows an example configuration file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dgs
  labels:
    component: dgs
spec:
  template:
    spec:
      containers:
        - name: dgs
          resources:
            limits:
              # Maximum CPU resource value. Can be specified in 1 milli units
              # Specifying 1 equals 1000m (= 1 core)
              # Specifying 800m equals 0.8 cores
              cpu: 1000m
            requests:
              # Minimum CPU resource value
              # Minimum CPU resources to reserve
              cpu: 1000m
```

To apply file changes, execute kubectl apply.

```sh
# Apply configuration to dev0 namespace
# Note that other files may also be changed
kustomize build k8s/aws/overlays/dev0 | kubectl apply -f -
```

To temporarily change without editing files, execute kubectl edit dgs.

```sh
# Edit dgs deployment resource in dev0 namespace (switches to editor)
kubectl edit deploy dgs -n dev0
```

```yaml
# Edit the cpu: item and save
: resources:
    limits:
      cpu: 800m
    requests:
      cpu: 800m
:
```

### Pod Count Configuration

The DGS pod count can be changed by editing the following file:
**k8s/aws/overlays/dev0/dgs/scaledobject.yaml**

- `spec.minReplicaCount` : Minimum pod count. Will not scale in below this value
- `spec.maxReplicaCount` : Maximum pod count. Will not scale out above this value
- `spec.advanced.scalingModifiers.formula` : By adjusting this value, you can control the number of standby DGS pods.
  - For example, to always have 5 units available, specify `takenCnt + 5`.
  - `takenCnt` represents the number already in use, and by adding to this, you can always keep available DGS servers running.

The following shows an example configuration file. For details, refer to the URLs in each comment.

```yaml
# This is the KEDA ScaledObject spec:
# https://keda.sh/docs/2.17/reference/scaledobject-spec/
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: dgs
spec:
  #  ────────────────────
  # Tell KEDA which workload to scale-up/down
  #  ────────────────────
  scaleTargetRef:
    name: dgs

  #  ────────────────────
  #   Optional tuning
  #  ────────────────────
  pollingInterval: 30 # seconds between polls (default: 30)
  cooldownPeriod: 30 # seconds to wait after a scale event

  minReplicaCount: 1
  maxReplicaCount: 10 # max number of replicas

  #  ────────────────────
  #   Trigger: call to the API and parse the JSON
  #  ────────────────────
  triggers:
    # Metrics API is one of the scalers which can parse the JSON response.
    # https://keda.sh/docs/2.17/scalers/metrics-api/
    - type: metrics-api
      name: takenCnt
      metadata:
        url: "http://http.dev0.svc.cluster.local/state/noderole/DGS"
        format: "json"
        valueLocation: "takenCnt"
        # targetValue: The target value used to determine whether to scale out or in based on the metric value.
        #   - desiredReplicas = metricValue / targetValue
        # However, since scalingModifiers is set, the actual replica count is calculated on scalingModifiers.
        targetValue: "1"
        # activationTargetValue: The activation threshold for starting/stopping the scaler itself.
        #   - If the metric value is less than activationTargetValue, the scaler is stopped.
        #   - If the metric value is greater than activationTargetValue, the scaler is started.
        activationTargetValue: "1"

  advanced:
    # scalingModifiers is the formula to calculate the total number of replicas.
    # https://keda.sh/docs/2.17/concepts/scaling-deployments/#scaling-modifiers
    scalingModifiers:
      # The formula should set the total number of replicas.
      # For example, if you need 5 online replicas(allocatable replicas),
      # you should set "takenCnt + 5".
      formula: "takenCnt + 1"

      # With AverageValue, the total replicas count is calculated as [formula / target].
      #   - totalReplicas = (takenCnt + N) / target = (takenCnt + 5) / 1
      #   - onlineReplicas = totalReplicas - takenCnt = N
      metricType: "AverageValue"
      target: "1"
```

To apply file changes, execute kubectl apply.

```sh
# Apply scaleObject configuration to dev0 namespace
kubectl apply -f k8s/aws/overlays/dev0/dgs/scaledobject.yaml -n dev0

# The following command also applies changes, but note that other file changes will also be applied
kustomize build k8s/aws/overlays/dev0 | kubectl apply -f -
```

To temporarily change without editing files, execute kubectl edit dgs.

```sh
# Edit dgs ScaledObject resource in dev0 namespace (switches to editor)
kubectl edit scaledobject dgs -n dev0
```

```yaml
# Edit items and save
: maxReplicaCount: 300
  minReplicaCount: 10
: scalingModifiers:
    formula: takenCnt + 10
:
```

### Node Count Configuration

Configure the node count so that the maximum pod count can be secured with sufficient nodes.
For configuration methods, please refer to the cloud service documentation.

## Other

### Cleanup

To delete the environment, execute the following command:

```sh
eksctl delete cluster --region=ap-northeast-1 --name=diarkis  --disable-nodegroup-eviction
```

---

_First created on 2025-07-14. Updated on 2025-11-27._

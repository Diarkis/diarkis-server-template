Place the DGS binaries here for docker build.

The DGS container is built by the following command:

```
# for dev0 namespace
make build-container-aws
```

Then, add the container to the container registry:

```
# for AWS ECR
make push-container-aws

# for GCP Artifact Registry
make push-container-gcp
```

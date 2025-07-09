Place the DGS binaries here for docker build.

The DGS container is built by the following command:

```
# for dev0 namespace
make build-container-with-dgs-aws-dev0
```

Then, add the container to the ECR repository:

```
# for dev0 namespace
make push-container-aws-dev0
```

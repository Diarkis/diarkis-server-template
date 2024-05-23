# Overview

This directory includes configuration files used by docker.
These config files are mainly for 2 purposes.

- Json files that are used by docker compose command to run Diarkis locally.
- Dockerfile's to create images used by both Docker compose and Kubernetes

If you want to try your images that are built for public clouds, change `diarkis-local` in `docker-compose.yml` to the url for your docker image registry.

# Run Diarkis locally with Docker Compose

```
# build must be linux arch
make build-local
make run-docker
```

**IMPORTANT**
7000/7200/TCP and 6779/7100/UDP must be open on the host side.

# Stop Diarkis

```
make stop-docker
```

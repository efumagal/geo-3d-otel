# geo-3d-otel

[![golangci-lint](https://github.com/efumagal/geo-3d-otel/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/efumagal/geo-3d-otel/actions/workflows/golangci-lint.yml)

## Introduction

A simple HTTP client to serve APIs to calculate distance between 3D points.  
Implemented using [Fiber](https://gofiber.io) and instrumented using [OpenTelemetry](https://opentelemetry.io).  

## CI/CD

### GitHub action
When pushing to main the [Publish Docker image to GHCR](.github/workflows/ghcr-build-push.yml) is triggered:  
- Build the Docker image and push it to ghcr.io registry
- Scan the docker image using [Snyk](https://snyk.io)
- Update [kustomization.yaml](kustomize/kustomization.yaml) with the newly generated Docker image

### K8s manifest Kustomize
For this example I'm using [Kustomize](https://kustomize.io) to build the K8s manifest, this is different than what I'm currently doing on real projects where I use [Helm Charts](https://helm.sh/docs/topics/charts/).
In order to keep everything in one repo the K8S manifests are kept in [kustomize](kustomize/) folder, unless using a monorepo they should probably be decoupled as the app code should be agnostic to the way is deployed. 
Similarly the files used by [Flux](https://fluxcd.io) are stored in [fluxCD](fluxCD/).  

In this case there are 4 manifests:
- [Deployment](kustomize/deployment.yaml) 
Standard deployment with liveness/readiness, using a secret for the Honeycomb token. The secrets are added using [Sealed secrets](https://github.com/bitnami-labs/sealed-secrets)
- [Service](kustomize/service.yaml)
- [HPA](kustomize/hpa.yaml)
In order to scale out when needed, Average CPU > 80%, min 2, max 6 pods
- [Namespace](kustomize/namespace.yaml)

### Flux
Flux is deployed on a local K8s cluster, and changes in the [kustomize](kustomize/) folder in this repo will trigger the reconciliation.  
On a real project there may be different variants for each env (deployment, staging, production).

## Run Locally

```shell
go run main.go
```

## Load test

Pre-requiste [K6](https://k6.io)

[Load test K6 script](k6-load/load_distance.js)

```shell
k6 run load_distance.js
```

## TO DOs

- For a real app consider structuring the Go code using DDD Hexagonal pattern
- Run Docker build and push only on related code changes
- Add unit tests and run them on PRs
- Generate OpenAPI specs

## Notes

- 
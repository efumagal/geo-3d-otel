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
Standard deployment with liveness/readiness, using a secret for the Honeycomb token. The secrets are added using [Sealed secrets](https://github.com/bitnami-labs/sealed-secrets).
- [Service](kustomize/service.yaml)
- [HPA](kustomize/hpa.yaml)
In order to scale out when needed, Average CPU > 80%, min 2, max 6 pods
- [Namespace](kustomize/namespace.yaml)

### Flux
Flux is deployed on a local K8s cluster, and changes in the [kustomize](kustomize/) folder in this repo will trigger the reconciliation.  
On a real project there may be different variants for each env (deployment, staging, production).  

An example of a merge can be seen in [CICD Video Recording](cicd.md).

## Observability

The code is instrumented with [OpenTelemetry](https://opentelemetry.io) and in this particular example the traces and metrics are exported to [Honeycomb.io](https://www.honeycomb.io).  
[Otel Fiber](https://github.com/gofiber/contrib/tree/main/otelfiber) is used for HTTP calls, plus there are some custom spans created to check specific operations.  

Heatmap of HTTP calls during a load test:  
![Heatmap](https://github.com/efumagal/geo-3d-otel/assets/77152760/69378de0-5dd6-41f1-8713-c3ab1c8212b9)

Single trace:
![Trace](https://github.com/efumagal/geo-3d-otel/assets/77152760/e8e33da7-26d0-4d93-971d-0a77d2fdbfd5)

## Run Locally

```shell
go run main.go
```

The webserver runs on port 8080.

### APIs

- **GET /health** Used for liveness/readiness probes
- **GET /hello** Returns a hello World HTML page
- **GET /distance** Calculates distance between two randomly generated 3D points
- **POST /action/cpu-load** Used to generate some CPU load for testing HPA 

[Geo3D.postman_collection](postman_collection/Geo3D.postman_collection) can be imported in [Postman](https://www.postman.com) to try the endpoints implemented.

## Testing

The Go code should have unit tests, this was not done for time constraints but it is a crucial part in guranteeing the correctness of the application.

### Load test

Pre-requiste [K6](https://k6.io)

The script [Load test K6 script](k6-load/load_distance.js) can be used to generate some load and test performance and instrumentation.

```shell
k6 run load_distance.js
```

## TO DOs

- For this example the K8s cluster is a local one (provided by Docker Desktop) and FluxCD was deployed manually. With more time I could write a Terraform script to set up AWS EKS and bootstrap FluxCD.
- For a real app consider structuring the Go code using DDD/Hexagonal pattern and separate the source code in a separate folder.
- Run Docker build and push only on related code changes
- Add unit tests and run them on PRs
- In case of a real app use REST and generate OpenAPI specs

## Notes

- Before deploying on production there will be probably be some functional/system automated tests

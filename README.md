# geo-3d-otel

[![golangci-lint](https://github.com/efumagal/geo-3d-otel/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/efumagal/geo-3d-otel/actions/workflows/golangci-lint.yml)

## Introduction

A simple HTTP client to serve APIs to calculate distance between 3D points.  
Implemented using [Fiber](https://gofiber.io) and instrumented using [OpenTelemetry](https://opentelemetry.io).  

When pushing to main the [Publish Docker image to GHCR](.github/workflows/ghcr-build-push.yml) is triggered:  
- Build the Docker image and push it to ghcr.io registry
- Scan the docker image using [Snyk](https://snyk.io)
- Update [kustomization.yaml](kustomize/kustomization.yaml) with the newly generated Docker image

In order to keep everything in one repo the K8S manifests are kept in (kustomize)[kustomize/] folder, unless using a monorepo they should probably be decoupled as the app code should be agnostic to the way is deployed.

## Run

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
- Add unit tests and run them on PRs
- Generate OpenAPI specs

## Notes

- 
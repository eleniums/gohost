# gohost

[![Go Report Card](https://goreportcard.com/badge/github.com/eleniums/gohost)](https://goreportcard.com/report/github.com/eleniums/gohost) [![GoDoc](https://godoc.org/github.com/eleniums/gohost?status.svg)](https://godoc.org/github.com/eleniums/gohost)

A tool for hosting a Go service with gRPC and HTTP endpoints.

## Installation

```
go get github.com/eleniums/gohost
```

## Prerequisites

- Requires Go 1.6 or later
- Uses [dep](https://github.com/golang/dep) for dependencies
  - `go get -u github.com/golang/dep/cmd/dep`
- Uses [grpc-go](https://github.com/grpc/grpc-go) for gRPC endpoints
- Uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) for REST endpoints

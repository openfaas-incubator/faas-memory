# faas-memory

faas-memory implements the OpenFaaS provider API and uses in-memory objects to store state.

## Purpose

This provider was created for integration testing and as a first-class example of an OpenFaaS provider for others to follow. It implements the [faas-provider](https://github.com/openfaas/faas-provider) interface.

## Getting started

In one terminal, build and start the provider:

```sh
export GOPATH=$HOME/go
go install github.com/yannip1234/faas-memory
cd faas-memory
go build
```

In another use the CLI with it:

```sh
$ export OPENFAAS_URL=127.0.0.1:8083

$ faas-cli list

Function                      	Invocations    	Replicas

$ faas-cli store deploy figlet

WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.

Deployed. 200 OK.
URL: http://127.0.0.1:8083/function/figlet

$ faas-cli list

Function                      	Invocations    	Replicas
figlet                        	0              	1
```


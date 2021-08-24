---
sidebar_position: 2
---

# Vanilla

You may deploy the Operator in a cluster without Helm.

## Requirements

Make sure to have Prometheus deployed in your target cluster if you want to scrape logs and metrics from the
`/metrics` endpoint. Deploy cert-manager to have webhooks enabled.

## Installation

You can use the official [Docker image](https://hub.docker.com/r/bedag/kubernetes-dbaas) "`bedag/kubernetes-dbaas`" or
build your own:

```
docker build -t yourrepo/imagename . && docker push yourrepo/imagename
make deploy IMG=yourrepo/imagename
```

## Additional information
For more information about the operator-sdk and the enclosed Makefile, consult: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

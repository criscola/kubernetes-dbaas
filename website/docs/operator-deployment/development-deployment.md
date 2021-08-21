---
sidebar_position: 3
---

# Local development

You can try out the Operator on your local development machine and boot the Operator as a normal system process. 

## Requirements

- Install Go 1.16 https://golang.org/doc/install
- Install kubectl v1.21+ https://kubernetes.io/docs/tasks/tools/install-kubectl/
- Install [kube-prometheus-stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) v17.1.3 
- Install minikube v1.21+ https://minikube.sigs.k8s.io/docs/start/
- Install the operator-sdk v1.6+ and its prerequisites: https://sdk.operatorframework.io/docs/installation/
- Have the Operator configured with your endpoint list and DatabaseClasses.

## Installation

1. Install the CRDs

```
make install
```

2. Install an example DatabaseClass

```
kubectl apply -f testdata/dbclass.yaml
```

3. Run the Operator as a local process

```
make run ARGS="--load-config=config/manager/controller_manager_config.yaml --enable-webhooks=false --leaderElection.leaderElect=false --debug=true"
```


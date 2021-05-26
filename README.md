# Kubernetes DbaaS
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

![Bedag](https://www.bedag.ch/wGlobal/wGlobal/layout/images/logo.svg)

*A Kubernetes Database as a Service (DBaaS) Operator for non-Kubernetes managed database management systems.*

## Abstract

This project aims at creating a Kubernetes Operator able to trigger a stored procedure in an external DBMS which in turn provisions a new database instance.
Users are able to create new database instances by writing the API Object configuration using Kubernetes Custom Resources.
The Operator watches for new API Objects and tells the target DBMS to trigger a certain stored procedure based on the custom resource specs.

## Motivation

There are many cases where a company can't or doesn't want to host their precious data in cloud or distributed environments and simply desire a way to bridge the gap between their K8s clusters and DBMS solutions. Imagine an organization composed by developers and system administrators, the former want their database provisioned ASAP whereas the latter want to have as much control as possible on the procedure needed to provision databases while still automating repetitive tasks. If this sounds interesting, keep reading.

## Main technologies

- Go 1.16 or newer
- operator-sdk v1.7.2 or newer 
- Kubernetes v1.21.0 or newer
- Helm v3

## Features

### Create database 

![k8s_dbaas_bedag_create](docs/resources/k8s_dbaas_bedag_create.png)

### Delete database

![k8s_dbaas_bedag_delete](docs/resources/k8s_dbaas_bedag_delete.png)

## Manuals

Set up the Operator using the Sysadmin guide. After that, end-users can use the end-user guide to learn how to provision a database through the Operator. 

- [System administrator guide](docs/sysadmin_guide.md)
- [End-user guide](docs/enduser_guide.md)

## Supported DBMS

- SQLServer
- PostgreSQL
- MySQL/MariaDB

### Additional notes

Encrypted DBMS connections are not supported.

## Quickstart with Helm
Other deployment options are shown in the [System administrator guide]().
### Helm
The Operator provides an official Helm chart.
#### Requirements
When metrics are enabled, the `/metrics` endpoint is protected by [authentication](https://github.com/brancz/kube-rbac-proxy) and scraped by Prometheus through a Service Monitor resource.
If you don't want to publish a `/metrics` endpoint, you may skip the following dependencies. 

- Install [kube-prometheus-stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) v16.0.1
```
helm install prometheus-operator prometheus-community/kube-prometheus-stack --create-namespace --namespace=prometheus
```
- Install [cert-manager](https://artifacthub.io/packages/helm/cert-manager/cert-manager) v1.3.0
```
helm install \                                                                                                                                                                                                                                                                                                                                                                                                                       ±[A1●●][develop]
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.3.1 \
  --set installCRDs=true
  
```
#### Deployment
1. Install the operator
```
helm install kubernetes-dbaas charts/kubernetes-dbaas --create-namespace --namespace=kubernetes-dbaas-system
```

## Known issues

## Code reference

To be done (godoc present on the code).

## Contribute

Please read the [contributing guidelines](docs/contributing.md). 

## Credits

## License

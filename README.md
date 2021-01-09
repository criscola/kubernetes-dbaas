# Kubernetes DbaaS
*A Kubernetes Database as a Service (DBaaS) Operator for non-Kubernetes managed database management systems.*

## Abstract

This project aims at creating a Kubernetes Operator able to trigger a stored procedure in an external DBMS which in turn provisions a new database instance.
Users are able to create new database instances by writing the API Object configuration using Kubernetes Custom Resources.
The Operator watches for new API Objects and tells the target DBMS to trigger a certain stored procedure based on the custom resource specs.

## Motivation

There are many cases where a company can't or doesn't want to host their precious data in cloud or distributed environments and simply desire a way to bridge the gap between their K8s clusters and DBMS solutions. Imagine an organization composed by developers and system administrators, the former want their database provisioned ASAP whereas the latter want to have as much control as possible on the procedure needed to provision databases while still automating repetitive tasks. If this sounds interesting, keep reading.

## Main technologies

- Go 1.15 or newer
- operator-sdk v1.2.x 
- Kubernetes v1.19.0 or newer
- Helm v3

## Features

### Create database 

![k8s_dbaas_bedag_create](docs/resources/k8s_dbaas_bedag_create.png)

### Delete database

![k8s_dbaas_bedag_delete](docs/resources/k8s_dbaas_bedag_delete.png)

### To-do

- Implement additional DBMS drivers (see supported DBMS)
- Test the controller with [KUTTL](https://github.com/kudobuilder/kuttl)
- Tests refactoring
- Extend the Helm chart for a larger number of use cases
- Support db connections encryption
- Maybe convert the current config.yaml to ConfigMap and Secrets

## Manuals

Please setup the Operator using the Sysadmin guide. After that, End-users or testers can use the End-user guide to learn how to provision a database through the Operator. 

Those who wish to contribute to the code should read the contributor guide.

- [System administrator guide](docs/sysadmin_guide.md)
- [End-user guide](docs/enduser_guide.md)
- [Contributor guide](docs/contributor_guide.md)

## Supported DBMS

- SQLServer
- PostgreSQL (to be implemented)
- MariaDB (to be implemented)

### Additional notes

The operator doesn't support encrypted DBMS connections yet.

## Quickstart

To try out the Operator on your local development machine, follow these steps:

1. Install Go 1.15+ https://golang.org/doc/install
2. Install kubectl v1.19+ https://kubernetes.io/docs/tasks/tools/install-kubectl/
3. Install minikube v1.16+ https://minikube.sigs.k8s.io/docs/start/
4. Install the operator-sdk and its prerequisites: https://sdk.operatorframework.io/docs/installation/
5. Configure the Operator by following the [System administrator guide](docs/sysadmin_guide.md)
6. `chmod +x start.sh`
7. `./start.sh`
8. Create and delete a custom resource by following the [End-user guide](docs/enduser_guide.md)

You can also use the supplied Dockerfile to compile your own Docker image. 

For more information about the operator-sdk and the enclosed Makefile, consult: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

## Code reference

To be done (godoc present on the code).

## Tests

## Contribute

## Credits

## License


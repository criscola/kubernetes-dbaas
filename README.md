# Kubernetes DbaaS
*A Kubernetes Database as a Service (DBaaS) Operator for non-Kubernetes managed database management systems.*

## Abstract

This project aims at creating a Kubernetes Operator able to trigger a stored procedure in an external DBMS which in turn provisions a new database instance.
Users are able to create new database instances by writing the API Object configuration using Kubernetes Custom Resources.
The Operator watches for new API Objects and tells the target DBMS to trigger a certain stored procedure based on the custom resource specs.

## Motivation

There are many cases where a company can't or doesn't want to host their precious data in cloud or distributed environments and simply desire a way to bridge the gap between their K8s clusters and DBMS solutions. Imagine an organization composed by developers and system administrators, the former want their database provisioned ASAP while the latter want to have as much control as possible on the procedure needed to provision databases while still automating repetitive tasks. If this sounds interesting, keep reading.

## Main technologies

- Go 1.15
- operator-sdk v1.2.0 

## Features

### Create database 

![k8s_dbaas_bedag_create](https://raw.githubusercontent.com/bedag/kubernetes-dbaas/develop/docs/resources/k8s_dbaas_bedag_create.png)

### Delete database

![k8s_dbaas_bedag_delete](https://raw.githubusercontent.com/bedag/kubernetes-dbaas/develop/docs/resources/k8s_dbaas_bedag_delete.png)

### To-do

- Implement additional DBMS drivers (see supported DBMS)
- Test the controller with [KUTTL](https://github.com/kudobuilder/kuttl)
- Helm chart with appropriate RBAC and monitoring
- Support db connections encryption

## Manuals

Please setup the Operator using the Sysadmin guide. After that, End-users or testers can use the End-user guide to learn how to provision a database through the Operator. 

Those who wish to contribute to the code should read the contributor guide.

- Sysadmin/DevOps guide
- End-user guide
- Contributor guide

## Supported DBMS

- SQLServer
- PostgreSQL (to be implemented)
- MariaDB (to be implemented)

### Additional notes

The operator doesn't support encrypted DBMS connections yet.

## Code reference

To be done (godoc present on the code)

## Tests

## Contribute

## Credits

## License


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

- Go 1.15 or newer
- operator-sdk v1.5.x 
- Kubernetes v1.19.0 or newer
- Helm v3

## Features

### Create database 

![k8s_dbaas_bedag_create](docs/resources/k8s_dbaas_bedag_create.png)

### Delete database

![k8s_dbaas_bedag_delete](docs/resources/k8s_dbaas_bedag_delete.png)

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
6. Install the CRDs
   
```
make install
```

7. Install an example DatabaseClass

```
kubectl apply -f testdata/dbclass.yaml
```

8. Run the Operator in local development mode 
   
```
make run ARGS="--load-config=config/manager/controller_manager_config.yaml --enable-webhooks=false --leaderElection.leaderElect=false --debug=true"
```

9. Create and delete a Database resource by following the [End-user guide](docs/enduser_guide.md)

For more information about the operator-sdk and the enclosed Makefile, consult: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

### Helm deployment
Make sure to have certmanager deployed in your target cluster.
```
helm install charts/kubernetes-dbaas --generate-name --create-namespace --namespace=kubernetes-dbaas-system
```

### Other deployment options
Make sure to have certmanager deployed in your local cluster.

You may deploy the Operator in a local cluster by running the following:

```
docker build -t yourrepo/imagename . && docker push yourrepo/imagename
make deploy IMG=yourrepo/imagename
```


## CLI arguments
|                                          	    | Description                                                                                                                          	|
|---------------------------------------------- |--------------------------------------------------------------------------------------------------------------------------------------	|
| `--debug`                                  	| Enables debug mode for development purposes                                                                                          	|
| `--enable-webhooks`                        	| Enables webhooks servers (default true)                                                                                               	|
| `--health.healthProbeBindAddress <string>` 	| The address the probe endpoint binds to (default ":8081")                                                                            	|
| `-h`, `--help`                               	| help for kubedbaas                                                                                                                   	|
| `--leaderElection.leaderElect`             	| Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager  (default true) 	|
| `--leaderElection.resourceName <string>`   	| The resource name to lock during election cycles (default "bfa62c96.dbaas.bedag.ch")                                                 	|
| `--load-config <string>`                   	| Location of the Operator's config file                                                                                               	|
| `--metrics.bindAddress <string>`           	| The address the metric endpoint binds to (default "127.0.0.1:8080")                                                                  	|
| `--webhook.port <int>`                       	| The port the webhook server binds to (default 9443)                                                                                  	|
| `--log-level <int>`                       	| The verbosity of the logger from 0 to 3 (default 1)                                                                                  	|

The order of precedence is `flags > config file > defaults`. Environment variables are not read.

## Troubleshooting
You can troubleshoot problems in two ways:
1. Look at the events of the resource with `kubectl describe database my-database-resource `
2. Consult the logs of the manager pod.

To avoid leaking possibly sensitive information, events do not contain the full error, only a message along with some
pertinent values if present.

You can also control the verbosity of the logger by setting the `--log-level <int>` flag to a value from 0 (only strictly
necessary logs) to 3 (very verbose). By default, this value is set to 1.

## Known problems

## Code reference

To be done (godoc present on the code).

## Tests

## Contribute

## Credits

## License


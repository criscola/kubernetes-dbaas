# System administrator guide

## Prerequisites

In order to use the Operator, your Database Management Systems (DBMS) must:

- Be supported by the Operator. See [supported DBMS](../README.md#supported-dbms).
- Be available and ready to accept connections from the Operator.
- Contain a `create`, `delete` and `rotate` stored procedure per DBMS as specified in the [Database administrator guide](dba_guide.md).

## Overview
![database resource apply overview](resources/database_resource_apply.svg)

## Operator configuration
The Operator is configured by 2 pieces of configuration, an `OperatorConfig` resource and `DatabaseClass` resources.
The order of precedence is the following:
1. It uses the configuration file specified through the `--load-config` CLI flag
2. It looks for a file named `config.yaml` in the same path as Operator binary
3. It looks for a file named `config.yaml` in `/etc/kubernetes-dbaas`

If you are using the [Helm deployment](#Helm deployment) option, you will specify this configuration in `.Values.operatorConfig` 
(example present in the chart), and it will automatically be referenced in the Operator's Pod through a ConfigMap.

If you are not using the Helm deployment option, it is sufficient to supply the plain configuration using one of the 3
possibilities highlighted before. 
See [config.example.yaml](../config.yaml.example) for an example.

### ComponentConfig
The first part contains configuration intended for the `kube-controller-manager` component.
You can find out more in its [godocs](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/config/v1alpha1#ControllerManagerConfigurationSpec)
and [design proposal](https://github.com/kubernetes-sigs/controller-runtime/blob/master/designs/component-config.md).

```yaml
health:
  healthProbeBindAddress: :8081
metrics:
  bindAddress: 127.0.0.1:8080
webhook:
  port: 9443
leaderElection:
  leaderElect: true
  resourceName: bfa62c96.dbaas.bedag.ch
  resourceNamespace: default
```
### Rate-limiting
As a protection against rogue users (or ugly bugs) the Operator embeds a rate-limiter for its DB drivers. You can specify
how many requests per second are allowed at maximum for an endpoint by setting `rps` to your desired number. If set to 0, 
the rate-limiter is disabled.

```yaml
rps: 1
```
### DBMS configuration
Endpoints should be configured thought the `dbms` key. As you can see, the Operator accepts an array formed by two 
keys, `databaseClassName` and `endpoints`.
- `databaseClassName` is a string specifying the name of the DatabaseClass resource associated with the accompanying 
  `endpoints` attribute.
- `endpoints` is an array containing the list of DBMS endpoints. It accepts two keys: `name` and `dsn`.
  - `name` is a convenient human-readable name associated with the endpoint. It must be unique in the overall list of endpoints and identifiable by end-users, so
    that they can refer to it when they want to create a new Database instance. Endpoint names must be properly documented inside your organization.
  - `dsn` is the [data source name](https://en.wikipedia.org/wiki/Data_source_name) used to connect to the DBMS endpoint.
    This project uses the `xo/dburl` package to parse DSN of different database drivers, you can find out more in [its documentation](https://github.com/xo/dburl).

```yaml
dbms:
  - databaseClassName: "databaseclass-sample-sqlserver"
    endpoints:
      - name: "us-sqlserver-test"
        dsn: "sqlserver://sa:Password&1@localhost:1433/master"
  - databaseClassName: "databaseclass-sample-psql"
    endpoints:
      - name: "us-postgres-test"
        dsn: "postgres://postgres:Password&1@localhost:5432/postgres"
  - databaseClassName: "databaseclass-sample-mariadb"
    endpoints:
      - name: "us-mariadb-test"
        dsn: "mariadb://root:Password&1@localhost:3306/mysql"
```

### Full example
Here's the full example:
```yaml
  health:
    healthProbeBindAddress: :8081
  metrics:
    bindAddress: 127.0.0.1:8080
  webhook:
    port: 9443
  leaderElection:
    leaderElect: true
    resourceName: bfa62c96.dbaas.bedag.ch
  rps: 1
  dbms:
    - databaseClassName: "databaseclass-sample-sqlserver"
      endpoints:
        - name: "us-sqlserver-test"
          dsn: "sqlserver://sa:Password&1@192.168.58.1:1433"
    - databaseClassName: "databaseclass-sample-psql"
      endpoints:
        - name: "us-postgres-test"
          dsn: "postgres://postgres:Password&1@192.168.58.1:5432/postgres"
    - databaseClassName: "databaseclass-sample-mariadb"
      endpoints:
        - name: "us-mariadb-test"
          dsn: "mariadb://root:Password&1@192.168.58.1:3306/mysql"
```

## DatabaseClass
DatabaseClass is the resource describing database operations. 
- `driver` expects a string declaring the driver to be used to execute database operations. It can be either `postgres`, `sqlserver`, `mysql` or `mariadb`.
- `operations` accepts 3 keys: `create`, `delete` and `rotate`. Each operation expects the same keys.
  - `name` expects a string specifying the name of the stored procedure as it is in the relative DBMS endpoint. The Operator will call it when the
    relative operation is triggered.
  - `inputs` expects an arbitrary map of values. Each key is the name of the parameter as specified in the stored procedure, while the value is
    the value supplied to it. See [Templating](#Templating) to learn more.
- `secretFormat` expects an arbitrary map of values. Each key is the name of the key as specified in the Secret resource created during the `create` operation,
  while the value is the value returned by the `create` stored procedure. You can find the values from the `create` operation by using the `.Result` top-level key.
      
```yaml
apiVersion: databaseclass.dbaas.bedag.ch/v1
kind: DatabaseClass
metadata:
  name: databaseclass-sample-psql
spec:
  driver: "postgres"
  operations:
    create:
      name: "sp_create_db_rowset_eav"
      inputs:
        k8sName: "{{ .Metadata.name }}"
    delete:
      name: "sp_delete"
      inputs:
        k8sName: "{{ .Metadata.name }}"
    rotate:
      name: "sp_rotate"
      inputs:
        k8sName: "{{ .Metadata.name }}"
  secretFormat:
    username: "{{ .Result.username }}"
    password: "{{ .Result.password }}"
    port: "{{ .Result.port }}"
    dbName: "{{ .Result.dbName }}"
    server: "{{ .Result.fqdn }}"
    dsn: "psql://{{ .Result.username }}:{{ .Result.password }}@{{ .Result.fqdn }}:{{ .Result.port }}/{{ .Result.dbName }}"
```

### Templating
DatabaseClasses support [Go templates](https://golang.org/pkg/text/template/) for operation inputs. Users can supply an 
arbitrary number of key-value pairs which will be mapped to the relative key as specified in the DatabaseClass responsible for their 
Database instance. See [End-user guide](enduser_guide.md).

The first key starts with a dot and has the first letter capitalized. There are two sources of values:

- `.Metadata` maps values from the `metadata` field of the Database resource.
- `.Params` maps values from the `spec.params` field of the Database resource.

For example, if an end-user has specified `spec.params.department: devops`, you could map it to the create stored procedure 
like that: `spec.operations.create.inputs.department: "{{ .Parameters.department }}"` and it will
be rendered ultimately as following: `department: 'devops'`.

Of course, you can hard-code your own values on a per-DatabaseClass basis, or if you're using the [Helm deployment](#Helm deployment) option, 
render values through the Helm chart.

> **Note:** If a key was specified, but a value was not found during rendering, the resource will generate an error.
Every `.Params.<key>` and `.Metadata.<key>` specified in the Operator configuration must be defined.
To define optional parameters, explicitly ask end-users to provide an empty string as value.

### MySQL/MariaDB
MySQL/MariaDB do not support supplying input parameters by name, only by position. Thus, in this case, the order of 
parameters matter. 

To work around this, DatabaseClasses specifying a MySQL/MariaDB driver must be adapted with a specific configuration.
Keys of `spec.operations.<operation>.inputs` must be integers describing their position relative to the other parameters, 
e.g. `"0"` becomes the first parameter to be supplied and `"1"` becomes the second.

For example:
```yaml
[..]
  driver: "mariadb"
  operations: 
    create:
      name: "sp_create_db_rowset_eav"
      inputs:
        "0": "{{ .Metadata.name }}"
        "1": "{{ .Params.department }}"
[..]
```

## Credential rotation
To order a DBMS to regenerate the credentials for a Database resource, apply the following annotation to the Database resource:

```yaml
dbaas.bedag.ch/rotate: ""
```

The Operator will remove the annotation when the operation has completed.

## CLI arguments
|                                          	    | Description                                                                                                                          	             	       |
|---------------------------------------------- |------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--debug`                                  	| Enables debug mode for development purposes. If set, logging output will be pretty-printed for the command line and `--log-level` will default to `1`        |
| `--disable-webhooks`                        	| Disables webhooks servers (default false)                                                                                               	                   |
| `--health.healthProbeBindAddress <string>` 	| The address the probe endpoint binds to (default ":8081")                                                                                                    |
| `-h`, `--help`                               	| Help for kubedbaas                                                                                                                                           |
| `--leaderElection.leaderElect`                | Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager  (default true)                                |
| `--leaderElection.resourceName <string>`   	| The resource name to lock during election cycles (default "bfa62c96.dbaas.bedag.ch")                                                                         |
| `--load-config <string>`                   	| Location of the Operator's config file                                                                                                                       |
| `--metrics.bindAddress <string>`           	| The address the metric endpoint binds to (default "127.0.0.1:8080")                                                                  	                       |
| `--webhook.port <int>`                       	| The port the webhook server binds to (default 9443)                                                                                  	             	       |
| `--log-level <int>`                       	| The verbosity of the logging output. Can be one out of: `0` info, `1` debug, `2` trace. If debug mode is on, defaults to `1` (default 0)                     |                                                                       	|
| `--disable-stacktrace`                       	| Disable stacktrace printing in logger errors (default false)                                                                                  	           |
| `--rps <int>`                                 | The number of operation executed per second per endpoint. If set to 0, operations won't be rate-limited (default 0)                                                                                  	           |

The order of precedence is `flags > config file > defaults`. Environment variables are not read.

## Installation
### Helm deployment
The Operator provides an official Helm chart.
#### Requirements
When metrics are enabled, the `/metrics` endpoint is protected by authentication and scraped by Prometheus through a Service Monitor resource.

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

### Vanilla deployment
You may deploy the Operator in a cluster without Helm.

Make sure to have cert-manager and Prometheus deployed in your target cluster if you want to have the `/metrics` endpoint enabled.
```
docker build -t yourrepo/imagename . && docker push yourrepo/imagename
make deploy IMG=yourrepo/imagename
```

### Local development deployment
To try out the Operator on your local development machine, follow these steps:

#### Requirements
- Install Go 1.16+ https://golang.org/doc/install
- Install kubectl v1.21+ https://kubernetes.io/docs/tasks/tools/install-kubectl/
- Install minikube v1.20+ https://minikube.sigs.k8s.io/docs/start/
- Install the operator-sdk v1.6+ and its prerequisites: https://sdk.operatorframework.io/docs/installation/
- Have the Operator configured with your endpoint list and DatabaseClasses.

#### Deployment
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

For more information about the operator-sdk and the enclosed Makefile, consult: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

## Samples
You can find a few samples of Database and DatabaseClass resources in the [testdata folder](../testdata).

## Troubleshooting
You can troubleshoot problems in two ways:
1. Look at the events of the resource with `kubectl describe database my-database-resource `
2. Consult the logs of the manager pod.

To avoid leaking possibly sensitive information, events do not contain the full error description, they contain a message along with some helpful values, if present.
Thus, you should refer to the logging output if you need to obtain deeper information about the problem.

You can control the verbosity of the logger by setting the `--log-level <int>` flag.

- `0`: Info level
- `1`: Debug Level
- `2`: Trace level

Errors are always logged.

Sampling is enabled in production mode for every log entry with same level and message. The first 100 entries in one second
are logged, after that only one entry is logged every 100 entries until the next second.

Stacktraces are attached to error logs in both production and development mode. You can disable this behavior by passing the `--disable-stacktrace=true` flag to the Operator binary.

### Monitoring tips
You could implement alerts on malfunctioning Database resources by watching Database events. You can find a complete list 
of Reasons and Messages [here](../pkg/typeutil/constants.go). Alternatively, if you find this too granular, you can
simply watch the `.status.conditions[*].status.type: Ready` field and check whether `.status.conditions[*].status.status`
equals `"False"`, if it does, the Database instance is in an error state, and that could generate an alert.

# System administrator guide
## Installation
### Helm
The Operator provides an official Helm chart.
#### Requirements
When metrics are enabled, the `/metrics` endpoint is protected by authentication and scraped by Prometheus through a Service Monitor resource.

- Install [kube-prometheus-stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) v16.0.1
```
helm install prometheus-operator prometheus-community/kube-prometheus-stack --create-namespace --namespace=kubernetes-dbaas-system
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
helm install charts/kubernetes-dbaas --generate-name --create-namespace --namespace=kubernetes-dbaas-system
```

### Vanilla deployment
To try out the Operator on your local machine, follow these steps:

#### Requirements
- Install Go 1.16+ https://golang.org/doc/install
- Install kubectl v1.21+ https://kubernetes.io/docs/tasks/tools/install-kubectl/
- Install minikube v1.20+ https://minikube.sigs.k8s.io/docs/start/
- Install the operator-sdk and its prerequisites: https://sdk.operatorframework.io/docs/installation/
- Configure the Operator by following the [System administrator guide](docs/sysadmin_guide.md)

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

### Other deployment options
Make sure to have certmanager and Prometheus deployed in your local cluster:

You may deploy the Operator in a local cluster without Helm, by running the following:

```
docker build -t yourrepo/imagename . && docker push yourrepo/imagename
make deploy IMG=yourrepo/imagename
```


## Usage
### Prerequisites

In order to work with the Operator, your DBMS must:

- Be supported by the Operator. See [supported DBMS](../README.md#supported-dbms).
- Be available and ready to accept connections from the Operator.
- Contain an **idempotent** stored procedure for database creation.
  - If the same ID of a **preexisting** DB is provided to the create stored procedure, it should always return the same output values associated with the DB instance.
- Contain an **idempotent** stored procedure for database deletion.
  - If the ID of a non-existing database is provided, the delete stored procedure should return an error.

Stored procedures inputs and outputs are treated as strings.

### Steps
1. Modify the DBMS configuration to suit your needs, the provided example contains the most minimal configuration required to run the Operator, you can find it under the file `config.example.yaml` or in the Helm chart under `operatorConfig`.
    1. Add a DBMS entry for each driver you want to support, e.g. one for sqlserver and another for postgresql. Don't modify the top-level attribute `dbms`.
    2. Add your endpoints under `endpoint`. End-users will use the endpoint `name` to associate their DB with a specific DBMS endpoint.
    3. Specify a `dsn` for your endpoint. See DSN docs for more information (note: currently DSN must be URL encoded by the user).
    4. Under `operations`, describe the `create` and `delete` stored procedures.
        1. For the `create` stored procedure:
            1. `name` is the name of the stored procedure as is in your DBMS.
            2. `inputs` contains input parameters for the stored procedure.
            3. `outputs` contains output parameters for the stored procedure. Currently, outputs do not support templating and are fixed.
                1. There are 5 output parameters: `username`, `password`, `fqdn`, `port`and `dbName`. When the stored procedure returns successfully, the Operator creates a Secret containing each of those output parameters, plus a `dsn` field constructed from those values. The key is fixed, while the value corresponds to the name of the parameters as it is written in the stored procedure.
        2. For the `delete` operation, the same format as the `create` operation applies.

## Templating

The Operator configuration supports [Go templates](https://golang.org/pkg/text/template/) for operation inputs (both create and delete). Users can supply an arbitrary number of keys and values (see [End-user guide](enduser_guide.md)).

Example of templated configuration:

```yaml
dbms:
  - driver: sqlserver
    endpoints:
      - name: us-sqlserver-test
        dsn: sqlserver://sa:Password&1@localhost:1433
    operations:
      create:
        name: sp_create
        inputs:
          k8sName: "{{ .Metadata.namespace }}/{{ .Metadata.name }}"
          paramName: "{{ .Parameters.myCustomUserParam }}"
          env: "dev"
        outputs:
          password: password
          username: username
          dbName: dbName
          fqdn: fqdn
          port: port
      delete:
        name: sp_delete
        inputs:
          k8sName: "{{ .Metadata.namespace }}/{{ .Metadata.name }}"
```

As you can see, the first key starts with a dot and has the first letter capitalized. There are two sources of values:

- `.Metadata`: maps values from the `metadata` field of the Database resource.
- `.Parameters`: maps values from the `paramas` field of the Database resource.

### Notes

- If a key is specified but not mapped by the user or Kubernetes, the resource will generate an error. Every `.Parameters.<key>` and `.Metadata.<key>` specified in the Operator configuration must be defined.
- If `metadata.namespace` is not set in the resource, the Operator will replace it with the value `default` during the rendering.
- If `metadata.name` is not set in the resource, the Operator will replace it with a 16-characters long, random alphanumeric string.


## DSN docs

- SQL Server: https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn (currently only URL format `driver://username:password@host/instance?param1=value&param2=value` is supported).

## User manual
The user manual is available

## CLI arguments
|                                          	    | Description                                                                                                                          	             	       |
|---------------------------------------------- |------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--debug`                                  	| Enables debug mode for development purposes. If set, `--log-level` defaults to `1`                                                                           |
| `--enable-webhooks`                        	| Enables webhooks servers (default true)                                                                                               	                   |
| `--health.healthProbeBindAddress <string>` 	| The address the probe endpoint binds to (default ":8081")                                                                                                    |
| `-h`, `--help`                               	| help for kubedbaas                                                                                                                                           |
| `--leaderElection.leaderElect`         | Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager  (default true)                                |
| `--leaderElection.resourceName <string>`   	| The resource name to lock during election cycles (default "bfa62c96.dbaas.bedag.ch")                                                                         |
| `--load-config <string>`                   	| Location of the Operator's config file                                                                                                                       |
| `--metrics.bindAddress <string>`           	| The address the metric endpoint binds to (default "127.0.0.1:8080")                                                                  	                       |
| `--webhook.port <int>`                       	| The port the webhook server binds to (default 9443)                                                                                  	             	       |
| `--log-level <int>`                       	| The verbosity of the logging output. Can be one out of: `0` info, `1` debug, `2` trace. If debug mode is on, defaults to `1` (default 0)                     |                                                                       	|
| `--disable-stacktrace`                       	| Disable stacktrace printing in logger errors (default false)                                                                                  	           |

The order of precedence is `flags > config file > defaults`. Environment variables are not read.

## Troubleshooting
You can troubleshoot problems in two ways:
1. Look at the events of the resource with `kubectl describe database my-database-resource `
2. Consult the logs of the manager pod.

To avoid leaking possibly sensitive information, events do not contain the full error, only a message along with some
pertinent values if present.

You can control the verbosity of the logger by setting the `--log-level <int>` flag.

- `0`: Info level
- `1`: Debug Level
- `2`: Trace level

Errors are always logged.

Sampling is enabled in production mode for every log entry with same level and message. The first 100 entries in one second
are logged, after that only one entry is logged every 100 entries until the next second.

Stacktraces are attached to error logs in both production and development mode.

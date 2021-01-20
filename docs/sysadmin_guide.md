# System administrator guide

## Prerequisites

In order to work with the Operator, your DBMS must:

- Be supported by the Operator. See [supported DBMS](../README.md#supported-dbms).
- Be available and ready to accept connections from the Operator.
- Contain an **idempotent** stored procedure for database creation.
  - If the same ID of a **preexisting** DB is provided to the create stored procedure, it should always return the same output values associated with the DB instance.
- Contain an **idempotent** stored procedure for database deletion.
  - If the ID of a non-existing database is provided, the delete stored procedure should return an error.

Stored procedures inputs and outputs are treated as strings.

## Steps
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

- `.Metadata`: maps values from the `metadata` field of the KubernetesDbaas resource.
- `.Parameters`: maps values from the `paramas` field of the KubernetesDbaas resource.

### Notes

- If a key is specified but not mapped by the user or Kubernetes, the resource will generate an error. Every `.Parameters.<key>` and `.Metadata.<key>` specified in the Operator configuration must be defined.
- If `metadata.namespace` is not set in the resource, the Operator will replace it with the value `default` during the rendering.
- If `metadata.name` is not set in the resource, the Operator will replace it with a 16-characters long, random alphanumeric string.



## DSN docs

- SQL Server: https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn (currently only URL format `driver://username:password@host/instance?param1=value&param2=value` is supported).

## Helm chart deployment

The operator can be deployed using the convenient Helm chart attached. As a minimum, you will have to:
1. Modify Helm's `values.yaml` file inside `kubernetes-dbaas` to suit your needs. Each attribute is commented on the chart itself except for the `operatorConfig` value, which is documented in [Steps](#Steps).
2. Build your own Docker image.
3. Reference the Docker image inside the Helm `values.yaml` file under `image.repository`.
4. Install the Helm chart.

Example of installation:

```
helm install --generate-name kubernetes-dbaas --namespace kubedbaas-system --create-namespace
```

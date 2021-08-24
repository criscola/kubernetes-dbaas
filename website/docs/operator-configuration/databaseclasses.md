---
sidebar_position: 3
---

# DatabaseClass 

## Format

DatabaseClass is the resource describing database operations.
- `driver` expects a string declaring the driver to be used to execute database operations. It can be either `postgres`, `sqlserver`, `mysql` or `mariadb`.
- `operations` accepts 3 keys: `create`, `delete` and `rotate`. Each operation expects the same keys.
    - `name` expects a string specifying the name of the stored procedure as it is in the relative DBMS endpoint. The Operator will call it when the
      relative operation is triggered.
    - `inputs` expects an arbitrary map of values. Each key is the name of the parameter as specified in the stored procedure, while the value is
      the value supplied to it. See [Templating](/docs/operator-configuration/databaseclasses#templating) to learn more.
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

DatabaseClasses are cluster-wide resources and do not belong to any namespace. They can be recalled on the command line
using the shorthand `dbc` instead of supplying the whole name.

## Templating
DatabaseClasses support [Go templates](https://golang.org/pkg/text/template/) for operation inputs. Users can supply an 
arbitrary number of key-value pairs which will be mapped to the relative key as specified in the DatabaseClass 
responsible for their  Database instance. See [Usage](/docs/usage).

The first key starts with a dot and has the first letter capitalized. There are two sources of values:

- `.Metadata` maps values from the `metadata` field of the Database resource.
- `.Params` maps values from the `spec.params` field of the Database resource.

For example, if an end-user has specified `spec.params.department: devops`, you could map it to the create stored procedure
like that: `spec.operations.create.inputs.department: "{{ .Parameters.department }}"` and it will
be rendered ultimately as following: `department: 'devops'`.

Of course, you can hard-code your own values on a per-DatabaseClass basis, or if you're using the [Helm deployment](/docs/operator-deployment/helm) option,
render values through the Helm chart.

:::caution

If a key was specified, but a value was not found during rendering, the resource will generate an error.
Every `.Params.<key>` and `.Metadata.<key>` specified in the Operator configuration must be defined.
To define optional parameters, explicitly ask end-users to provide an empty string as value.

:::

## Caveats
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

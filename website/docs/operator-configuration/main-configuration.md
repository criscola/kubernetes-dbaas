---
sidebar_position: 2
---

# Main configuration

## Operator configuration

The Operator is configured by 2 pieces of configuration, an `OperatorConfig` resource and `DatabaseClass` resources.
The order of precedence is the following:
1. It uses the configuration file specified through the `--load-config` CLI flag
2. It looks for a file named `config.yaml` in the same path as Operator binary
3. It looks for a file named `config.yaml` in `/etc/kubernetes-dbaas`

If you are using the [Helm deployment](/docs/operator-deployment/helm) option, you will specify this configuration in `.Values.operatorConfig`
(example present in the chart), and it will automatically be referenced in the Operator's Pod through a ConfigMap.

If you are not using the Helm deployment option, it is sufficient to supply the plain configuration using one of the 3
possibilities highlighted before.
See [config.example.yaml](https://github.com/bedag/kubernetes-dbaas/blob/main/config.yaml.example) for an example.

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

As a protection against rogue users or as a limitation to batch operations, the Operator embeds a rate-limiter for its
DB drivers. You can specify how many requests per second are allowed at maximum for an endpoint by setting `rps` to your
desired number. If set to 0, the rate-limiter is disabled.

```yaml
rps: 1
```

### Keepalive

It is possible to enable a keepalive mechanism that checks every *X* seconds whether there is a connection issue between
the Operator and the DBMS endpoints. The following configuration option lets users specify the interval in seconds
between each retry. If a connection is found to be down, the error is logged using the standard logger, and a connection
retry is performed. If the option is set to `0`, the keepalive is disabled.

```yaml
keepalive: 30
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

#### DSN in Secrets

It is recommended to have the DSN of DBMS in their own Secret when deploying the Operator in production.
Secrets can be referenced using the key `secretKeyRef` for each endpoint. `secretKeyRef.name` references the name of the
Secret and `secretKeyRef.key` references the key containing the DSN respectively.

Secrets must be placed in the same
namespace of the Operator Pod and must be present before booting the Operator. Updates to those Secrets do not trigger
any automatic update to the Operator configuration; if a Secret is updated while the Operator is running, the updated
configuration will be loaded during the next Operator boot.

Example for `secretKeyRef`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: "us-sqlserver-test-secret"
type: Opaque
stringData:
  dsn: "sqlserver://sa:Password&1@localhost:1433/master"
```

```yaml
dbms:
  - databaseClassName: "databaseclass-sample-sqlserver"
    endpoints:
      - name: "us-sqlserver-test"
        secretKeyRef:
          name: "us-sqlserver-test-secret"
          key: "dsn"
```

If both `dsn` and `secretKeyRef` are specified, `dsn` will take the precedence, and the referenced Secret will not be pulled.

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
  keepalive: 30
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

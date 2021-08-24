---
sidebar_position: 1
---

# Helm

The Operator provides an official Helm chart.

## Requirements

Install [kube-prometheus-stack](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack) v17.1.3 
used for scraping logs and metrics:

```shell
helm install prometheus-operator prometheus-community/kube-prometheus-stack --create-namespace --namespace=prometheus
```

Install [cert-manager](https://artifacthub.io/packages/helm/cert-manager/cert-manager) v1.4.0 
used for handling TLS certificates for the webhooks:

```shell
helm install \                                                                                                                
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.4.0 \
  --set installCRDs=true
```

## Installation

Finally, install the Operator's Helm Chart:

```shell
helm install kubernetes-dbaas charts/kubernetes-dbaas --create-namespace --namespace=kubernetes-dbaas-system
```

## Helper templates

This Helm Chart contains useful helper templates to facilitate the deployment of the Operator.

### DatabaseClasses generator

The top-level key `dbc` contains an array of entries describing DatabaseClass resources. Each array entry generate one DatabaseClass resource.

Example:

```yaml
dbc:
  - name: "databaseclass-sample-postgres"
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
      lastRotation: "{{ .Result.lastRotation }}"
      dsn: "sqlserver://{{ .Result.username }}:{{ .Result.password }}@{{ .Result.fqdn }}:{{ .Result.port }}/{{ .Result.dbName }}"
```

This entry will be translate into a DatabaseClass and deployed. Its structure is analogous to a standard DatabaseClass spec, only it does not contain Kubernetes-specific fields, such as `spec` and `metadata`. Moreover, it contains an additional key `dbc[*].name` which is rendered as the `metadata.name` of the resource.

### DBMS Secrets generator

The top-level key `dbmsSecrets` contains an array of entries describing Secrets resources which can be referenced in endpoint configurations inside of the `dbms[*].endpoints.secretKeyRef` keys. Each array entry generates one Secret resource.

Example:

```yaml
dbmsSecrets:
  - name: "us-sqlserver-test-credentials"
    stringData:
      dsn: "sqlserver://sa:Password&1@192.168.49.1:1433/master"
  - name: "us-postgres-test-credentials"
    stringData:
      dsn: "postgres://postgres:Password&1@192.168.49.1:5432/postgres"
  - name: "us-mariadb-test-credentials"
    stringData:
      dsn: "mariadb://root:Password&1@192.168.49.1:3306/mysql"
```

`name` is mapped to `metadata.name` and `stringData` is rendered as YAML into `spec.stringData` of the generated Secret.

## Additional information

Additional documentation can be found directly in the Chart's [README](https://github.com/bedag/kubernetes-dbaas/blob/main/charts/kubernetes-dbaas/README.md).
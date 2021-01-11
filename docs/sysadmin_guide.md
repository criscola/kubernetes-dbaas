# System administrator guide

## Prerequisites

In order to work with the Operator, your DBMS must:

- Be supported by the Operator. See [supported DBMS](../README.md#supported-dbms).
- Be available and ready to accept connections from the Operator.
- Contain an **idempotent** stored procedure for database creation.
  - As **input** parameters, there must be at least one parameter which is the ID associated with the database. At the moment, we provide the K8s resource UID as ID because it is guaranteed to be unique and it's easy to use for debugging purposes if something goes wrong. 
  - As **output** parameter, there must be at least two parameters which are username and password used to connect to the generated DB instance.
  - If the same ID of a **preexisting** DB is provided to the create stored procedure, it should always return the same username and password associated with the DB instance.
- Contain an **idempotent** stored procedure for database deletion.
  - As **input** parameter, there must be at least one parameter which is the ID associated with the database. 	 
  - If the ID of a non-existing database is provided, the delete stored procedure should return an error.

Consequently, stored procedures are responsible for:

- Retaining the binding between DB IDs and their real DB name used for db creation and deletion.
- Retaining the binding between DB IDs and their own username and password. 

Stored procedures inputs and outputs are treated as strings.

Other ad-hoc functionalities inside the stored procedures can be customized as needed if they don't break the prerequisites shown before.

## Steps
1. Modify the DBMS configuration to suit your needs, the provided example contains the minimal configuration required to run the Operator.
    1. Add a dbms entry for each driver you want to support, e.g. one for sqlserver and another for postgresql. Don't modify the top-level attribute `dbms`.
    2. Add your endpoints under `endpoint`. End-users will use the endpoint `name` to associate their DB with a specific DBMS.
    3. Specify a DSN for your endpoint. See DSN docs for more information.
    4. Under `operations`, describe the create and delete stored procedures.
        1. `name` is the name of the stored procedure as is in your DBMS.
        2. `inputs` contains an attribute called `k8sName` which is the ID associated with a DB instance. You can define your own parameter name if you need.
        3. `outputs` contains attributes called `username` and `password`, you can define your own naming if you need.
        4. `delete` contains the attribute `k8sName` as well, to reference which DB instance should be deleted.

## DSN docs

- SQL Server: https://github.com/go-sql-driver/mysql#dsn-data-source-name

## Helm chart deployment

The operator can be deployed using the convenient Helm chart attached. As a minimum, you will have to:
1. Modify Helm's `values.yaml` file inside `kubernetes-dbaas` to suit your needs. Each attribute is commented on the chart
   itself except for the `operatorConfig` value, which is documented in [Steps](#Steps).
2. Build your own Docker image.
3. Reference the Docker image inside the Helm values.yaml file under `image.repository`.
4. Install the Helm chart.

Example of installation:

```
helm install --generate-name kubernetes-dbaas --namespace kubedbaas-system --create-namespace
```

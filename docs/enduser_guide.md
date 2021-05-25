# End-user guide

## Prerequisites

You have the Operator ready and deployed. See [the system administrator guide.](sysadmin_guide.md)

## Steps

1. Create a custom resource configuration, e.g. `my-db.yaml`

Here's an example:

```
apiVersion: dbaas.bedag.ch/v1
kind: Database
metadata:
  name: database-sample  
spec:
  endpoint: us-sqlserver-test
  params:
    myCustomUserParam: "myvalue"
```
- `endpoint` defines which DBMS endpoint is responsible for the database instance. Endpoint names are configured in the Operator configuration and should be properly documented inside your organization.
- `params` defines a key-value map of parameters to be supplied to the Operator. Parameters are configured in the Operator configuration and should be properly documented inside your organization. 
  Extra parameters are ignored.

2. Apply the resource:

This will create a new database instance and a Secret resource with the database credentials in the same namespace as your request. 
Secret are named `<your-db-name>-credentials`.

```
kubectl apply -f my-db.yaml
```

3. Delete the resource:

This will delete the relative database instance. The Secret associated with it will be garbage collected.

```
kubectl delete -f my-db.yaml
```

## Troubleshooting

In case your database instance wasn't created successfully, the Operator will write events to the relative Database resource.

```
kubectl describe my-db
```


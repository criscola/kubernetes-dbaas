---
sidebar_position: 2
---

# Operations

## Supported operations

The Operator supports the following operations:
- Database instance creation
- Database instance deletion
- Database instance credential rotation

Operations are implemented on database management systems using their native technique of creating stored procedures.

Each operation maps to a single stored procedure loaded in advance into the targeted DBMS.

:::tip
It is possible to have more than one
stored procedure for each operation, if required (e.g. one for testing and one for production usage). 
:::

## Design specification

The **inputs, output result and body** of a stored procedure are implemented by 
complying to a minimal set of guidelines, i.e. the "**design specification**". This will allow the Operator to rely on the 
stored procedures for the actual implementation of the various operations, calling them whenever needed. 
System administrators will then be able to map inputs and outputs of the stored procedures so that the Operator will 
know how to handle them.

Stored procedures should be properly documented inside your organization:
- Name of the stored procedure as it is stored in its DBMS (case-sensitive).
- Name of each input parameter (see also [Notes](#Notes)).
- List of the resulting outputs and their description.

Input parameters can be supplied directly by end-users or by the infrastructure as configured by system administrators. 
There are no hard-coded input or output parameters.

Every input/output value is treated as a `string` (i.e. `TEXT`, `varchar`...).

:::caution
The Operator can be configured to supply a unique ID to stored procedures. An ID will be bound to a particular database instance.
Make sure to always specify an input parameter acting as ID. 
:::

Stored procedures that return data, will return a row set with at least two columns, named **precisely**
`key` and `value`. 

Example of output result:

| key      	    | value    	|
|-------------- |----------	|
| username 	    | <string\>	|
| password 	    | <string\>	|
| host     	    | <string\>	|
| port     	    | <string\>	|
| dbName   	    | <string\>	|
| lastRotation  | <string\>	|

Errors should be returned using the built-in mechanism of the DBMS involved, e.g. using exceptions. If an exception is returned,
the Operator will re-execute the operation again using an exponential back-off strategy.

Operations should be designed as [idempotent](https://en.wikipedia.org/wiki/Idempotence) (even though the Operator 
will not call an operation if not needed).

### Create

The create operation generates a database instance available to end-users. If the operation succeeds, it returns the connection
data needed in order to connect to and manipulate the newly provisioned database instance, along with
additional information if desired (e.g. timestamp of the latest credential rotation).

#### Caveat

As a rule, the create operation is called only during database creation, but it can be called twice with the same ID in 
the case of data loss about Database resources on the Kubernetes cluster.
Of course, passwords should be stored salted in hash form, so if the create operation is called again, it should
call [Rotate](#Rotate) internally and return the updated values.  Empty strings will overwrite previously filled strings.
All the key-value pairs returned with the first call must be returned.

### Delete

The delete operation deletes a database instance.
The operation will return nothing if the delete operation succeeded.

### Rotate
The rotate operation rotates the Database credentials and must returns the **same** keys as the `create` procedure,
with updated credentials. Empty strings will overwrite previously filled strings.

## Notes
### MySQL/MariaDB
Unfortunately, MySQL/MariaDB do not support supplying input parameters by name, only by position. Thus, in this case, 
the order of the parameters matter and should be documented carefully in order of appearance in the stored procedure.
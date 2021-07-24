t# Database administrator guide
## Prerequisites
In order to use the Operator, your Database Management Systems (DBMS) must be supported by the Operator. 
Please see [supported DBMS](../README.md#supported-dbms).

## Stored procedures
### Supported operations
The Operator supports the following operations:
- Database instance creation
- Database instance deletion
- Database instance credential rotation

The DBA should create at least one stored procedure for each operation per DBMS. It is possible to have more than one 
stored procedure for each operation, if required (e.g. one for testing and one for production usage).

### Common information
The inputs, outputs parameters and body of a stored procedure is defined by the DBA and should be communicated to system 
administrators in order to let them configure the Operator correctly. 

Stored procedure should be properly documented inside your organization:
- Name of stored procedure as it is stored in its DBMS
- Name of each input parameter (see also [Notes](#Notes))

Input parameters can be supplied directly by end-users or by infrastructure.

Every input/output value is treated as a `string` (i.e. `TEXT`, `varchar`...).

Errors should be returned using the built-in mechanism of the DBMS involved, e.g. using exceptions. If an exception is returned,
the Operator will re-execute the operation again using an exponential back-off strategy.

In order to identify each distinct database instance, it is necessary **to have at least an input parameter 
acting as ID for each operation**. Each store procedure call will provide an ID so that store procedures know which 
Kubernetes resource is bound to its relative database instance. This input value will be supplied by the Operator during
stored procedure calls.

### Create
The `create`Row operation will return a set of rows with at least two columns, `key` and `value`. Example:

| key      	    | value    	|
|-------------- |----------	|
| username 	    | \<string>	|
| password 	    | \<string>	|
| host     	    | \<string>	|
| port     	    | \<string>	|
| dbName   	    | \<string>	|
| lastRotation  | \<string>	|

If the `create` operation is called twice with the same ID, it should return the same values or updated values.
Of course, passwords should be stored salted in hash form, so if the `create` operation is called again, it should
call [Rotate](#Rotate) internally and return the updated values. As a rule, the `create` operation is called only during 
database creation.

### Delete
The operation will return nothing if the delete operation succeeded.

### Rotate
The `rotate` operation rotates the Database credentials and returns the **same** keys as the `create` procedure, 
with updated values if necessary. 

## Notes
### MySQL/MariaDB
MySQL/MariaDB do not support supplying input parameters by name, only by position. Thus, in this case, the order of the
parameters matter and should be documented carefully in order of appearance in the stored procedure.

### Samples
A few primitive samples for the `create` `rotate` (only for PostgreSQL) and `delete` stored procedures are present in 
the [testdata folder](../testdata).
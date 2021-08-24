---
sidebar_position: 3
---

# Testing

## Overview

Testing is achieved by using [Ginkgo](https://github.com/onsi/ginkgo), a Behavior-Driven Development
testing framework.

There are 3 types of tests:

- Unit
- Integration
- End-to-end (e2e)

Unit tests do not interact with any external service, and the data used for testing is composed by simple
stubs which are compared to the output of the tested components.

Integration tests interact with one external service, e.g. a DBMS.

E2e are a form of integration testing, which includes the interaction with the Kubernetes API.
End-to-end tests are created using [envtest](https://sdk.operatorframework.io/docs/building-operators/golang/testing/) 
which allows running a minimal
cluster by providing kube-apiserver, kubectl and etcd binaries. This means also that there
aren’t other controllers for built-in resources, for example it is not possible to test if a Secret is successfully
garbage-collected after deleting a Database resource because there is no controller controller watching for
Secret deletion, which on the other hand would be the case for a standard cluster installation.

## How-to

The whole test suite is executed when commits are pushed to the main branch (see [CI pipeline](/docs/contributing/ci)),
but it is much faster to execute it locally during development. 

### Setup DBMS
First, you need to have the supported DBMS available for testing. You can start a few Docker containers for those. See
the following snippet:

```shell
docker run -p 3306:3306 --name mariadb -e 'MARIADB_ROOT_PASSWORD=Password&1' -d mariadb
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Password&1" -p 1433:1433 --name sqlserver -h sqlserver -d mcr.microsoft.com/mssql/server:2019-latest
docker run --name psql1 -p 5432:5432 -e POSTGRES_PASSWORD="Password&1" -d postgres
```

Simply use an ETL tool such as DBeaver or HeidiSQL to load all the stored procedures contained in the 
`testdata/procedures` folder.

### Test configuration
Next, you need to create a test configuration. You can find an example [here](https://github.com/bedag/kubernetes-dbaas/tree/main/testdata)
along with all the resources and stored procedure you might need. 

Tweak the configuration based on what you want to test
along with the necessary DatabaseClasses and credentials for the DB systems. You can find the sample resources in the
`testdata/resources` folder.

Now, specify the full path of your testing configuration by setting the `TEST_CONFIG_PATH` environment variable.

### Executing the test suite

The first time it is advised to run `make test`, which will install the necessary envtest binaries and start
the test suite, afterwards the Ginkgo CLI can be used directly:

```shell
ginkgo -r -v
```

Alternatively, you can make use of a cluster installation of your choice by setting the `TEST_USE_EXISTING_CLUSTER` environment
variable to `true`, in this case the predefined cluster installation will be used, e.g. your local Minikube instance, and
you will need to load the CRDs via `make install`.

Unit and integration tests can be executed in parallel by using the `-p` option of Ginkgo. 
It is possible to focus on the desired test suite portion by using regex:

To run only integration tests in the `database` package:
```shell
ginkgo -focus=".*[i]*." -v pkg/database
```

To run only unit tests in the `database` package:
```shell
ginkgo -focus=".*[u]*." -v pkg/database
```

To run only e2e tests:
```shell
ginkgo -focus=".*[e2e]*." -v
```

To generate the test coverage report:

```shell
ginkgo -v -r -cover -coverprofile=coverage.out -outputdir=testdata/cover
go tool cover -html=testdata/cover/coverage.out -o testdata/cover/coverage_report.html
```

### Writing tests

It is advised to read [Writing controller tests](https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html) on
the Kubebuilder docs for anything related to writing controller tests while using envtest.

There are no particular caveats, other than commenting each step of a test, especially if longer than a few lines, is pretty
much required. 

One thing to keep in mind when writing tests, is to make use of the
`test.FormatTestName(label TestType, description string, extras ...TestAttribute)` method to format the description of tests,
this way regex can be used to include or exclude test cases from execution. Refer to its godocs to learn more.

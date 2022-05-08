# Archiving notice

This project has been temporarely put on ice (set as read-only). The repo's ownership has been transferred from the `bedag` profile to mine before archiving, so there are some problems such as resource groupversions which would need to be adjusted.

Unfortunately for this project, I got a lot of (positive) things going on in my life that leave me out of time to properly maintain the operator in the short and medium-term. That being said, the code is there and with some degree of effort it should work just fine. 
Someone even removed the stored procedures part and included the actual provisioning inside of the operator (something I wanted to do after delivering the operator to Bedag), but I could not find enough time to review the pull-request... still it shows that with some effort the operator can be very well adapted to the needs of your organization.
I still wish to update and enhance this operator a few years down the line, but for now it will stay archieved. 
Thank you for your interest and feel free to get in touch if you have any questions.

# Kubernetes DBaaS
### *A unique Kubernetes Database-as-a-Service (DBaaS) Operator for declarative, self-service database provisioning in DBMS solutions.*

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/bedag/kubernetes-dbaas.svg)](https://pkg.go.dev/github.com/bedag/kubernetes-dbaas)
[![Go Report Card](https://goreportcard.com/badge/github.com/bedag/kubernetes-dbaas)](https://goreportcard.com/report/github.com/bedag/kubernetes-dbaas)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubernetes-dbaas)](https://artifacthub.io/packages/helm/kubernetes-dbaas/kubernetes-dbaas)
[![Test Suite](https://github.com/bedag/kubernetes-dbaas/actions/workflows/go.yaml/badge.svg)](https://github.com/bedag/kubernetes-dbaas/actions/workflows/go.yaml)

<img alt="Kubernetes DBaaS Logo" src="website/static/img/logo.svg" width="300" height="260">

All documentation can be found on the [Kubernetes DBaaS website](https://bedag.github.io/kubernetes-dbaas).

## Description

The **Kubernetes Database-as-a-Service (DBaaS) Operator** ("the Operator") is a [Kubernetes
Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) used
to provision database instances in database management systems:

- The Operator can be easily configured and installed in a Kubernetes cluster
  using the provided Helm Chart.
- End-users such as software developers are able to create new database
  instances by writing simple Database custom resources.
- Operations on DBMS are implemented using stored procedures called by the
  Operator whenever necessary, allowing you to define your own custom logic.
- Credentials to access provisioned database instances are saved into Kubernetes
  Secrets.
  
Written using [Go](https://go.dev/) and the [Operator SDK](https://github.com/operator-framework/operator-sdk).
  
## Main features

- Level-based logging
- Event logging, metrics
- Credential rotation
- Helm deployment
- Rate-limited requests
- Modern tech-stack
- Flexible configuration

## Supported DBMS

- SQLServer
- PostgreSQL
- MySQL/MariaDB

## Contributing

There is a comprehensive set of guidelines and helpful documentation.

Please read [How to contribute](https://bedag.github.io/kubernetes-dbaas/docs/contributing/how-to-contribute) to get started.

## Credits

Copyright © 2021 [Cristiano Colangelo](https://github.com/criscola) —
Developed for Bedag Informatik AG.

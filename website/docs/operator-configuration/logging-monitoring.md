---
sidebar_position: 5
---

# Logging & troubleshooting

## Overview

By default `--debug` is set to false, in this case the production logger is initialized.
If `--debug` is set to true the development logger is initialized instead.

The development logger has pretty printing, suited for command-line consumption.

The production logger has JSON output and sampling.

Make sure to have a look at [CLI arguments](/docs/operator-configuration/cli-arguments).

## Troubleshooting



You can troubleshoot problems in two ways:
1. Look at the events of the resource with `kubectl describe db my-database-resource `
2. Consult the logs of the manager pod.

To avoid leaking possibly sensitive information, events do not contain the full error description, they contain a message along with some helpful values, if present.
Thus, you should refer to the logging output if you need to obtain deeper information about the problem.

You can control the verbosity of the logger by setting the `--log-level <int>` flag.

- `0`: Info level
- `1`: Debug Level
- `2`: Trace level

Errors are always logged.

### Example of event

The following shows what is written to the Database resource when it is written to the cluster and the relative DB instance
is provisioned successfully on the DBMS endpoint:

```shell
Type    Reason                    Age     From                   Message
----    ------                    ----    ----                   -------
Normal  DatabaseCreateInProgress  1s      database-controller    database instance is being provisioned on dbms endpoint
Normal  DatabaseCreateSuccess     1s      database-controller    database instance provisioned successfully on dbms endpoint
Normal  SecretCreateSuccess       1s      database-controller    secret created successfully: {"secret": "database-sample-789-credentials"}
```



## Sampling

Sampling is enabled in production mode for every log entry with same level and message. The first 100 entries in one second
are logged, after that only one entry is logged every 100 entries until the next second.

## Metrics

Metrics are information related to the operational status of the Operator, e.g. how much memory it
consumes. A ServiceMonitor resource is included, which enables Prometheus to scrape metrics from the
Operator. Metrics are written to the `/metrics` endpoint of the Operator.

Metrics are protected using [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy), a small HTTP proxy that can perform RBAC authorization
against the Kubernetes API. It is deployed alongside the controller, acting as a proxy for inbound requests.

See also the [kubebuilder documentation](https://book.kubebuilder.io/reference/metrics.html) about metrics.

## Additional information

Stacktraces can be enabled by setting the flag `--enable-stacktrace` to `true`. Defaults to `true` in debug mode.
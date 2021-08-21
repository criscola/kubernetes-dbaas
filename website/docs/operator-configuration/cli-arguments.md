---
sidebar_position: 6
---

# CLI arguments

The order of precedence is `flags > config file > defaults`. Environment variables are not read.

|                                                 | Description                                                  |
| ----------------------------------------------- | ------------------------------------------------------------ |
| `-h`, `--help`                                  | Help for kubedbaas                                           |
| `--debug <bool>`                                | Enables debug mode for development purposes. If set, logging output will be pretty-printed for the command line and `--log-level` will default to `1` |
| `--disable-webhooks <bool>`                     | Disables webhooks servers (default `false`)                  |
| `--health.healthProbeBindAddress <string>`      | The address the probe endpoint binds to (default `:8081`)    |
| `--leaderElection.leaderElect <bool>`           | Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager (default `true`) |
| `--leaderElection.resourceName <string>`        | The resource name to lock during election cycles (default `bfa62c96.dbaas.bedag.ch`) |
| `--leaderElection.resourceNamespace <string>`   | The namespace in which to create the leader election lock resource (defaults to the namespace of the Operator Pod) |
| `--load-config <string>`                        | Location of the Operator's config file                       |
| `--metrics.bindAddress <string>`                | The address the metric endpoint binds to (default `127.0.0.1:8080`) |
| `--webhook.port <int>`                          | The port the webhook server binds to (default `9443`)        |
| `--log-level <int>`                             | The verbosity of the logging output. Can be one out of: `0` info, `1` debug, `2` trace. If debug mode is on, defaults to `1` (default 0) |
| `--enable-stacktrace <bool>`                    | Enable stacktrace printing in logger errors, If debug mode is on, defaults to `true` (default `false`) |
| `--rps <int>`                                   | The maximum number of operations executed per second per endpoint. If set to `0`, operations won't be rate-limited (default `0`) |
| `--keepalive <int>`                             | The interval in seconds between connection checks for the endpoints (default `30`) |

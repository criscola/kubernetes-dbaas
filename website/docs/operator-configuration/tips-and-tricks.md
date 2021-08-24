---
sidebar_position: 7
---

# Tips & tricks

### Monitoring
You could implement alerts on malfunctioning Database resources by watching Database events. You can find a complete list
of Reasons and Messages [here](https://github.com/bedag/kubernetes-dbaas/blob/main/pkg/typeutil/constants.go). Alternatively, if you find this too granular, you can
simply watch the `.status.conditions[*].status.type: Ready` field and check whether `.status.conditions[*].status.status`
equals `"False"` for a certain number of time, if it does, that could generate an alert.

### Credential rotation
Credential rotation could be performed periodically by using a simple `CronJob` for each Database resource in the cluster, see the [bitnami/kubectl](https://hub.docker.com/r/bitnami/kubectl) Docker image.

### Restoring resources
If something bad happened and you've lost all your Database resources, you can simply reapply all your Database yaml files. Given that the `create` and `rotate` operations were implemented according to the specification, Database resources will be regenerated without the need of manual intervention.

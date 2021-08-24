---
sidebar_position: 4
---

# Credential rotation

To order a DBMS to regenerate the credentials for a Database resource, apply the following annotation to the Database resource:

```yaml
dbaas.bedag.ch/rotate: "true"
```

The Operator will attempt to rotate the credentials immediately. The Operator will remove the annotation once the
operation has completed successfully.

Credential rotation can be triggered also when a Secret resource generated during a create operation is deleted by the 
user.
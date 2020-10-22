---
description: ""
weight: 3
---

# Server Classes

Server classes are a way to group distinct server resources.
The "qualifiers" key allows the administrator to specify criteria upon which to group these servers.
There are currently three keys: cpu, systemInformation, and labelSelectors.
Each of these keys accepts a list of entries.
The top level keys are a "logical AND", while the lists under each key are a "logical OR".
Qualifiers that are not specified are not evaluated.

An example:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
metadata:
  name: default
spec:
  qualifiers:
    cpu:
      - manufacturer: Intel(R) Corporation
        version: Intel(R) Atom(TM) CPU C3558 @ 2.20GHz
      - manufacturer: Advanced Micro Devices, Inc.
        version: AMD Ryzen 7 2700X Eight-Core Processor
    labelSelectors:
      - "my-server-label": "true"
```

Servers would only be added to the above class if they had _EITHER_ CPU info, _AND_ the label associated with the server resource.

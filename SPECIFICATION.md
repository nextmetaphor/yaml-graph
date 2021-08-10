# Specification

### Class
Each file contains `0..n` definitions of a particular `Class` which must be specified at the top-level as is shown in the example below:

```yaml
Class: "Service"
```

There must be exactly one such statement per file.

### References
If **all** definitions in the file share a number of the same references to other definitions, these can be stated in the (optional) top-level `References` section, as is shown in the example below. Within the (optional) `References` section there can be `0..n` actual references; in the example below there are two.

```yaml
References: [
  { Class: "Provider", ID: "azure", Relationship: "PROVIDED_BY" },
  { Class: "Category", ID: "compute", Relationship: "TYPE_OF" }
]
```

### Definitions
The (optional) top-level `Definitions` section contains the actual definitions for the `Class` specified previously, as shown in the example below.
Within the (optional) `Definitions` section there can be `0..n` actual definitions; in the example below there are two.

```yaml
Definitions:
  app-service:
    Fields:
      Name: "App Service"
      Description: "A fully managed platform for building, deploying and scaling your web apps"
      Link: "https://azure.microsoft.com/en-gb/services/app-service/"

  azure-function:
    Fields:
      Name: "Azure Function"
      Description: "Azure Function"
      Link: "https://azure.microsoft.com/en-gb/services/functions/"
```
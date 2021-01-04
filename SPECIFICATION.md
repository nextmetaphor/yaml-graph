# Specification

## Class
Each file contains `0..n` definitions of a particular `Class` which must be specified at the top-level as is shown in the example below:

```yaml
Class: "Animal"
```

There must be exactly one such statement per file.

## References
If **all** definitions in the file share a number of the same references to other definitions, these can be stated in the optional top-level `References` section, as is shown in the example below:

```yaml
References: [
  Class: "",
  ID: "",
  Relationship: ""
]
```

## Definitions
The (optional) top-level `Definitions` section contains the actual definitions for the `Class` specified previously. 
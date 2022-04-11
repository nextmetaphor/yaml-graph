# `Specification` Schema
```yaml
# MANDATORY [1] field which represents the class of all definitions within this file 
Class: string

# OPTIONAL [0..1] field which, if specified, must be an map of Definition objects, keyed by ID. There can be `0..n` definitions within this map.
Definitions: map[string]Definition

# OPTIONAL [0..1] field which, if specified, must be an array of reference objects. There can be `0..n` references within this array.
References: []Reference
```

# `Definition` Schema
```yaml
# OPTIONAL [0..1] field which represents the fields of the definition
Fields: map[string]string

# OPTIONAL [0..1] field which represents the fields of the definition which are stored in external files
FileFields: map[string]FileDefinition

# OPTIONAL [0..1] field which represents the references
References: []Reference

  # OPTIONAL [0..1] field which represents any subdefinitions that this definition has
SubDefinitions: map[string]Specification
```

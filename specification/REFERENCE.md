# `yaml-graph` Definition Quick Reference
This document provides a definitive reference as to the structure of the definition files used by `yaml-graph`.
Refer to the [detailed documentation](README.md) to understand the purpose of each field and structure, and for
examples. The applicable source code is [here](../src/definition/definition.go).

Each definition file must contain an implementation of the `Specification` schema as detailed below.

## `Specification` Schema
```yaml
# MANDATORY field which represents the class of all definitions within this file.
Class: string

# OPTIONAL field which, if specified, must be a map of Definition objects, keyed by ID. There can be `0..n` definitions
# within this map.
Definitions: map[string]Definition

# OPTIONAL field which, if specified, must be an array of reference objects. There can be `0..n` references within this
# array.
References: []Reference
```

## `Definition` Schema
```yaml
# OPTIONAL field which represents the fields of the definition, keyed by field name.
Fields: map[string]string

# OPTIONAL field which represents the fields of the definition which are stored in external files, keyed by field name.
FileFields: map[string]FileDefinition

# OPTIONAL field which represents the references for this definition.
References: []Reference

# OPTIONAL field which represents any sub-definitions that this definition has, keyed by relationship.
# Note that using this approach to define relationships doesn't currently permit direction to be indicated or
# relationship fields. Refer to the Reference section for further details.
SubDefinitions: map[string]Specification
```

## `FileDefinition` Schema
```yaml
# MANDATORY field which specifies the path to the file.
Path: string

# OPTIONAL field which specifies a prefix to be added to the underlying definition.
Prefix: string

# OPTIONAL field which specifies the encoding to be used. Defaults to no encoding, but "base64" can also be specified.
Encoding: string
```

## `Reference` Schema
```yaml
# MANDATORY field which represents the class of the field to which a reference is being made.
Class: string

# MANDATORY field which represents the ID of the field of the class specified above to which a reference is being made.
ID: string

# MANDATORY field which represents the type of the relationship between the definitions.
Relationship: string

# OPTIONAL field which indicates whether the relationship is directed from the Definition containing this Relationship. Defaults to false.
RelationshipFrom: bool

# OPTIONAL field which indicates whether the relationship is directed to the Definition containing this Relationship. Defaults to false.
RelationshipTo: bool

# OPTIONAL field which specifies any additional data in the form of fields that should be attached to the relationship.
Fields: map[string]string
```
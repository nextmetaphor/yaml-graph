# Specification
Each `yaml` file is used to specify `0..x` definitions of a particular type, referred to as a `Class`.  

The top-level elements of each such `yaml` file are as follows:
```yaml
Class: "class_name"
Definitions: [{definition_object_1}, {definition_object_2}, ...]
References: [{reference_object_1}, {reference_object_2}, ...]
```

With regard to 'standard' graph database terminology:
* A `Definition` is synonymous with a `Node`
* A `Reference` is synonymous with an `Edge`

Each of these top-level elements are discussed in the sections below.

## Top-Level Elements
### `Class` Element
There must be **a single** `Class` element per file, and must be a `string`. All subsequent `Definitions` within the file will be of this `Class`.

For example:
```yaml
Class: "Service"
```
### `References` Element
There can be **at most** one `References` section per file, and it must be an array of reference objects. There can be `0..n` references within this array.

**ALL** references specified in this element will be applied to **ALL** definitions within the file. Therefore, if each definition is required to have different references, then these should be defined within the definition itself, and not the top-level `References` section. The use of the top-level `References` section is useful to avoid duplication for when every definition in this file has the same reference.

Each element in the `References` array is a reference object with the following attributes:
* `Class`: a mandatory `string` attribute which specifies to the class of definition to which the reference refers to to
* `ID`: a mandatory `string` attribute which specifies the ID of the definition of the class specified above
* `Relationship`: a mandatory `string` attributes which specifies the nature of the relationship between the source and target definition
* `RelationshipFrom`: an optional `boolean` attribute which specifies whether the relationship is directed from the source to the target definition
* `RelationshipTo`: an optional `boolean` attribute which specifies whether the relationship is directed to the source from the target definition

Note that bidirectional relationships, when the relationship is directed from the source to the target **AND** from the target to the source are valid and permitted.

For example:
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

The two elements directly below the `Definitions` element, in this case `app-service` and `azure-function`, specify the `ID` for each of the definitions. Each definition must have a unique `ID` for the `Class` that it belongs to.

#### `Fields` Element
This section unsurprisingly contains the fields that make up the content of the definition. Each field is a name/value pair where the key is the field name, and the value is the data associated with the field. Note that currently only `string` values are supported.

For both of the definitions in the example above, there are three fields: `Name`, `Description` and `Link`.

#### `FileFields` Element
The `FileFields` section is an alternative method to populate the `Fields` of a definition, but instead allows the use of a separate file for the value of the field. This is useful if the content is a more complex Markdown document that is preferable to maintain separately from the main definition YAML file, or an image file.




#### `References` Element
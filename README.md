# yaml-graph
golang-based utility to help define, build and build reports on graph representations of data models from simple YAML files.

*Note: this utility is currently still in development and is subject to breaking changes.*  

## Install

### Prerequisites
* Local [Docker](https://www.docker.com/) installation: compilation, test and execution of `yaml-graph` is implemented within a Docker container
* Local [make](https://www.gnu.org/software/make/) installation: a `makefile` is used to co-ordinate the activities listed above

### Building
`yaml-graph` is compiled and packaged into a Docker container, complete with a [neo4j](https://neo4j.com) installation for ease as follows:
```bash
# clone the yaml-graph repository
$ git clone git@github.com:nextmetaphor/yaml-graph.git

# move into the root of the repository
$ cd yaml-graph

# invoke the docker-build target in the makefile
$ make docker-build
``` 

## Usage

### Running
When running a `yaml-graph` Docker container, we need to mount one or more directories on the host machine containing the YAML definitions within the container. In the example below, we mount host directory `$(PWD)/example-definitions/CloudTaxonomy` which is a sample directory provided in this repository for example purposes only. We will usually also mount a directory containing report definitions to allow us to build reports; in the example below this is set to  
```bash
docker run -it -p7474:7474 -p7687:7687 -v $(PWD)/example-definition/CloudTaxonomy:/home/ymlgraph/definition -v $(PWD)/example-report:/home/ymlgraph/report nextmetaphor/yaml-graph
```

From within the Docker container, the definition directory provided will then be mounted at `/home/ymlgraph/definitions`. This directory will need to be specified for the majority of the `yaml-graph` operations, as demonstrated below. 


### Command Line Options
From within the running container, use the `--help` flag to examine the various command line options available.

```bash
                       _                             _
 _   _  __ _ _ __ ___ | |       __ _ _ __ __ _ _ __ | |__
| | | |/ _` | '_ ` _ \| |_____ / _` | '__/ _` | '_ \| '_ \
| |_| | (_| | | | | | | |_____| (_| | | | (_| | |_) | | | |
 \__, |\__,_|_| |_| |_|_|      \__, |_|  \__,_| .__/|_| |_|
 |___/                         |___/          |_|
version:0.3.6

yaml-graph $ yaml-graph --help
Define data in YAML then generate graph representations to model relationships

Usage:
  yaml-graph [command]

Available Commands:
  help        Help about any command
  load        Load definition files into graph representation
  report      Generate report from graph representation
  validate    Validate definition files
  version     Print the version number of yaml-graph

Flags:
  -d, --dbURL string      URL of graph database (default "bolt://localhost:7687")
  -e, --ext string        file extension for definitions (default "yaml")
  -h, --help              help for yaml-graph
  -l, --logLevel int8     log level (0=debug, 1=info, 2=warn, 3=error) (default 2)
  -p, --password string   password (default "password")
  -u, --username string   username for graph database (default "username")

Use "yaml-graph [command] --help" for more information about a command.
```

### Validate Definitions
To validate the YAML definitions, execute the following command:
```bash
yaml-graph $ yaml-graph validate -f definition/definition-format.yml -s definition
successfully validated definitions
```
Note that multiple definition source directories can be supplied. This is useful, for example, if you are referencing a separate taxonomy that you want to reference, but it makes sense to control the definition files separately.
```bash
# note: example only
yaml-graph $ yaml-graph validate -s BASE_DEFINITION_DIRECTORY -s ADDITIONAL_DEFINITION_DIRECTORY
successfully validated definitions
```

### Load Definitions
To load the YAML definitions into a graph representation, execute the following command:
```bash
yaml-graph $ yaml-graph load -s definition
```

### Visualise Graph Representation
Examine the graph database structure at http://localhost:7474/browser/ using the CYPHER of `match (n) return n`

### Generate Report
To generate a report from the loaded graph representation, based on specified classes and fields together with a
[gotemplate](https://golang.org/pkg/text/template/) report definition, execute the following command:
```bash
yaml-graph $ yaml-graph report -f report/fields.yaml -t report/template.gohtml > report/output.html
```
The HTML report is available on the host machine at `$(PWD)/example-report`

## Licence
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This project is licenced under the terms of the [Apache 2.0 License](LICENCE.md) licence.
# yaml-graph
golang-based utility to help define and build graph representations of data models from simple YAML files.

*Note: this utility is currently still in development and is subject to breaking changes. The only representation currently supported uses the [neo4j](https://neo4j.com) graph database; additional future representations planned include [d3](https://d3js.org/) and [pdf](https://www.adobe.com/devnet/pdf/pdf_reference.html).*  

## Install

### Prerequisites
* Local [golang](https://golang.org/) installation: for building purposes
* Local [Docker](https://www.docker.com/) installation: to host a local [neo4j](https://neo4j.com) database
* Running [neo4j]() database container, which should be started with `docker run -p7474:7474 -p7687:7687 --env=NEO4J_AUTH=none neo4j`

### Build Steps
To build `yaml-graph`, follow the steps below.
```bash
# clone the yaml-graph repository
$ git clone git@github.com:nextmetaphor/yaml-graph.git
Cloning into 'yaml-graph'...
remote: Enumerating objects: 149, done.
remote: Counting objects: 100% (149/149), done.
remote: Compressing objects: 100% (88/88), done.
remote: Total 149 (delta 55), reused 122 (delta 31), pack-reused 0
Receiving objects: 100% (149/149), 34.36 KiB | 429.00 KiB/s, done.
Resolving deltas: 100% (55/55), done.

# move into the src directory
$ cd yaml-graph/src/
$ cd yaml-graph/src/
src $

# execute the build.sh script
$ ./build.sh 
### go fmt ###
### go vet ###

### golint ###
### go build ###
``` 

## Testing
To test `yaml-graph`, follow the steps below.
```bash
# move into the src directory of the yaml-graph repository
$ cd yaml-graph/src/

# execute the test.sh script
src $ ./test.sh 

### go test ###
?   	github.com/nextmetaphor/yaml-graph	[no test files]
?   	github.com/nextmetaphor/yaml-graph/cmd	[no test files]
ok  	github.com/nextmetaphor/yaml-graph/definition	0.345s	coverage: 30.3% of statements in ./...
ok  	github.com/nextmetaphor/yaml-graph/graph	0.191s	coverage: 15.1% of statements in ./...

### go tool cover ###
```

## Deployment

### Command Line Options
Use the `--help` flag to examine the various command line options available.

```bash
$ yaml-graph --help
Define data in YAML then generate graph representations to model relationships

Usage:
  yaml-graph [command]

Available Commands:
  help        Help about any command
  parse       Parse definition files into graph representation
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

### Running
The `yaml-graph` repository has a directory which contains a number of sample definitions then can be used for test
purposes. Generate a graph representation of this data set as follows.
```bash
# note: execute the command below from the src directory from which the code was built
$ ./yaml-graph parse -s ../example/CloudTaxonomy
```

## Validation
Examine the graph database structure at http://localhost:7474/browser/ using the CYPHER of `match (n) return n`
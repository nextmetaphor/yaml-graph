# `yaml-graph` Unit Test Instructions
### Build `yaml-graph` Image
```bash
# Execute from the root of the yaml-graph git repository
yaml-graph $ make docker-build
...
[+] Building 0.8s (11/11) FINISHED  
````

### Start the `yaml-graph` Container
```bash
# Execute from the root of the yaml-graph git repository
yaml-graph $ docker run -it -p7474:7474 -p7687:7687 -v $(PWD)/example-definition:/home/ymlgraph/definition -v $(PWD)/example-report:/home/ymlgraph/report nextmetaphor/yaml-graph 
                        _                             _
 _   _  __ _ _ __ ___ | |       __ _ _ __ __ _ _ __ | |__
| | | |/ _` | '_ ` _ \| |_____ / _` | '__/ _` | '_ \| '_ \
| |_| | (_| | | | | | | |_____| (_| | | | (_| | |_) | | | |
 \__, |\__,_|_| |_| |_|_|      \__, |_|  \__,_| .__/|_| |_|
 |___/                         |___/          |_|
version:0.3.18
```

### Validate and Load the Test Definitions 
```bash
# Execute from within the yaml-graph container:
yaml-graph $ yaml-graph validate -f definition/CloudTaxonomy/definition-format.yml -s definition
successfully validated definitions

yaml-graph $ yaml-graph load -s definition
```

### Execute the Unit Tests
```bash
# Execute from outside of the container, in the root of the yaml-graph git repository:
yaml-graph $ cd src
src $ ./test.sh 

### go test ###
?       github.com/nextmetaphor/yaml-graph      [no test files]
ok      github.com/nextmetaphor/yaml-graph/cmd  0.472s  coverage: 7.7% of statements in ./...
?       github.com/nextmetaphor/yaml-graph/cui  [no test files]
ok      github.com/nextmetaphor/yaml-graph/definition   0.599s  coverage: 14.0% of statements in ./...
ok      github.com/nextmetaphor/yaml-graph/graph        0.251s  coverage: 6.8% of statements in ./...
ok      github.com/nextmetaphor/yaml-graph/parser       1.419s  coverage: 38.7% of statements in ./...

### go tool cover ###
```
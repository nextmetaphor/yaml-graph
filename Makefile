# variable for common working directory and build cache arguments
docker_dir_args = -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest

.PHONY: help
help:	## show makefile help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:	## build yaml-graph using a docker build container
	# optionally pass GOOS and GOARCH parameters e.g. make build GOOS=darwin GOARCH=amd64
	docker run --rm $(docker_dir_args) ./build.sh $(GOOS) $(GOARCH)

	# copy the built binary to the docker installation files
	cp src/yaml-graph docker/utils

test:	## test yaml-graph using a docker test container
	docker run --rm $(docker_dir_args) ./test.sh

docker-build: build	## build yaml-graph docker image
	docker build --tag nextmetaphor/yaml-graph:latest docker

docker-run: docker-build
	docker run -it -p7474:7474 -p7687:7687 -v $(PWD)/example-definition/CloudTaxonomy:/home/ymlgraph/definition -v $(PWD)/example-template:/home/ymlgraph/report nextmetaphor/yaml-graph
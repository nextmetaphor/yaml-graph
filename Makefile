# variable for common working directory and build cache arguments
docker_dir_args = -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest

help:	## show makefile help
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

build:	## build yaml-graph using a docker build container
	## optionally pass GOOS and GOARCH parameters e.g. make build GOOS=darwin GOARCH=amd64
	docker run --rm $(docker_dir_args) ./build.sh $(GOOS) $(GOARCH)

test:	## test yaml-graph using a docker test container
	docker run --rm $(docker_dir_args) golang:latest ./test.sh
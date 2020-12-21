build:
	printf "### go fmt ###\n"
	docker run --rm -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest go fmt ./...

	printf "### go vet ###\n"
	docker run --rm -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest go vet ./...

	printf "\n### golint ###\n"
	docker run --rm -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest go get -u golang.org/x/lint/golint; golint ./...

	printf "### go build ###\n"
	docker run --rm -v $(PWD)/src:/usr/src/yaml-graph -v $(PWD)/.buildcache:/go -w /usr/src/yaml-graph golang:latest go build -v
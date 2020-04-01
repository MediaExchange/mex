.PHONY: all docker mex ui
.SILENT: dockerList

HASH=$(shell git rev-parse --short HEAD)

# Dummy target that lets mex run locally.
all: ui mex

# Build the docker container
dockerBuild:
	docker build -t mex:latest -t mex:${HASH} ${PWD}

# List the files in the docker container. Useful for verifying that
# everything is in the expected locations.
dockerList:
	docker create --name="tmp_mex" mex:latest > /dev/null
	docker export tmp_mex | tar t
	docker rm tmp_mex > /dev/null

# Builds the mex executable.
mex:
	go build

# Builds the angular UI.
ui:
	cd ui && npm update && ./node_modules/.bin/ng build --prod=true

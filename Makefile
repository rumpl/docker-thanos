STATIC_FLAGS=CGO_ENABLED=0
GO_BUILD=$(STATIC_FLAGS) go build
OUTPUT=docker-thanos
DIR := ${CURDIR}

plugin:
	$(GO_BUILD) -o $(OUTPUT) cmd/main.go

link:
	ln -sf $(DIR)/docker-thanos ~/.docker/cli-plugins/docker-thanos

.PHONY: plugin link
.DEFAULT: plugin

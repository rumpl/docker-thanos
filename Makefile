LDFLAGS="-s -w"
STATIC_FLAGS=CGO_ENABLED=0
GO_BUILD=$(STATIC_FLAGS) go build -ldflags=$(LDFLAGS)
OUTPUT=thanos

cmd:
	$(GO_BUILD) -o $(OUTPUT) cmd/main.go

.PHONY: cmd
.DEFAULT: cmd

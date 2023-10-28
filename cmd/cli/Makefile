BIN=cli

BUILDFLAGS=CGO_ENABLED=0
LDFLAGS=-ldflags "-X github.com/n4-networks/openusp/pkg/cli.buildtime=`date +%Y-%m-%d.%H:%M:%S`"

build: *.go
	@echo "Building OpenUSP CLI..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

install: *.go
	@echo "Installing OpenUSP CLI..."
	$(BUILDFLAGS) go install $(LDFLAGS)

.PHONY: clean build install
clean:
	go clean -x -i github.com/n4-networks/openusp/cmd/cli/...

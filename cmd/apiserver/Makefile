BIN=apiserver

BUILDFLAGS=CGO_ENABLED=0 GOOS=linux 
LDFLAGS=-ldflags "-X github.com/n4-networks/openusp/pkg/cntlr.buildtime=`date +%Y-%m-%d.%H:%M:%S`"

build: *.go
	@echo "Building OpenUSP ApiServer..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

install: *.go
	@echo "Installing OpenUSP ApiServer..."
	$(BUILDFLAGS) go install $(LDFLAGS)

.PHONY: clean build install
clean:
	go clean -x -i github.com/n4-networks/openusp/cmd/apiserver/...

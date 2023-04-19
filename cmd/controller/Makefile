BIN=mtp

BUILDFLAGS=CGO_ENABLED=0 GOOS=linux 
LDFLAGS=-ldflags "-X github.com/n4-networks/server.buildtime=`date +%Y-%m-%d.%H:%M:%S`"

$(BIN): *.go ../*.go
	@echo "Building N4 MTP..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

.PHONY: clean
clean:
	go clean -n github.com/n4-networks/cmd/server/...
	rm mtp

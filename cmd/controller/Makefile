BIN=controller
ACTIVEMQ=10.5.0.3:61613
DB_ADDR=10.5.0.2:27017
DB_USER=admin
DB_PASSWD=admin
MQTT_ADDR=10.5.0.3:1883
CACHE_ADDR=10.5.0.4:6379
LOGGING=all

BUILDFLAGS=CGO_ENABLED=0 GOOS=linux 
LDFLAGS=-ldflags "-X github.com/n4-networks/openusp/pkg/cntlr.buildtime=`date +%Y-%m-%d.%H:%M:%S`"

build: *.go
	@echo "Building OpenUSP Controller..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

install: *.go
	@echo "Installing OpenUSP Controller..."
	$(BUILDFLAGS) go install $(LDFLAGS)

.PHONY: clean build install
clean:
	go clean -x -i github.com/n4-networks/openusp/cmd/controller/...

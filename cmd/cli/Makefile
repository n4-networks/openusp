
#AGENT=os::N4a_9-02:42:ac:11:00:0b
#AGENT=os::012345-B827EB73AF33
#AGENT=os::N4a_8-02:42:ac:11:00:0d

# Docker
AGENT=os::N4est-02:42:ac:11:00:06

# Raspberry PI: LAN: 
#AGENT=os::012345-B827EB351734

# Local
#ACTIVEMQ=172.17.0.3:61613
REST=http://172.17.0.2:8081

# GCloud
#ACTIVEMQ=35.240.207.101:61613

LOGGING=all
BIN=cli

BUILDFLAGS=CGO_ENABLED=0 GOOS=linux
LDFLAGS=-ldflags "-X github.com/n4-networks/openusp/cmd/cli.buildtime=`date -u +%Y-%m-%dT%H:%M:%SZ`"
RUNFLAGS=AGENT_ID=$(AGENT) REST_SRV_ADDR=$(REST) CLI_LOGGING=$(LOGGING)

$(BIN): ../*.go
	@echo "Building N4 CLI..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

# command line mode
runc: $(BIN)
	$(RUNFLAGS) ./$(BIN)

.PHONY: testall
testall:
	$(RUNFLAGS) go test -v ../test

# The below target to be used to run a specefic test or a set of tests by providing tc="regex"
# example; make test tc=Ip where test cases starting with Ip will be executed alone
# By default go test will have cache enabled which means if the tests were run before and there were no changes on source or test code results will be shown from cache instead of running it again
# To make tests always running atleast once added count value to 1, please see more details using "go help testflag"
test: 
	$(RUNFLAGS) go test -v ../test -run $(tc) -count=1

.PHONY: clean
clean:
	go clean -n github.com/n4-networks/openusp/cmd/cli/cmd/...
	rm $(BIN)

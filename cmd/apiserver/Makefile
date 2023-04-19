# Local
#ACTIVEMQ=172.17.0.3:61613
MONGO=172.17.0.4:27017
MTP=:9001
HTTP=8081

# GCloud
#ACTIVEMQ=35.240.207.101:61613
#MONGO=34.87.135.72:27017
#MTP=34.126.123.116:9001
DBUSER=n4admin
DBPASSWD=n4defaultpass

LOGGING=all
BIN=rest


BUILDFLAGS=CGO_ENABLED=0 GOOS=linux 
LDFLAGS=-ldflags "-X github.com/n4-networks/rest.buildtime=`date +%Y-%m-%d.%H:%M:%S`"


$(BIN): *.go ../*.go
	@echo "Building N4 RESt Server..."
	$(BUILDFLAGS) go build -o $(BIN) $(LDFLAGS)

TEST_FLAGS=HTTP_PORT=$(HTTP) MTP_GRPC_ADDR=$(MTP) DB_ADDR=$(MONGO) DB_USER=${DBUSER} DB_PASSWD=$(DBPASSWD) LOGGING=$(LOGGING)
TEST_FLAGS_TLS=$(TEST_FLAGS) HTTP_TLS=1

.PHONY: testall
testall:
	$(TEST_FLAGS) go test -v ../test

# The below target to be used to run a specefic test or a set of tests by providing tc="regex"
# example; make test tc=Ip where test cases starting with Ip will be executed alone
# By default go test will have cache enabled which means if the tests were run before and there were no changes on source or test code results will be shown from cache instead of running it again
# To make tests always running atleast once added count value to 1, please see more details using "go help testflag"
test: 
	$(TEST_FLAGS) go test -v ../test -run $(tc) -count=1

.PHONY: clean
clean:
	go clean -n github.com/n4-networks/rest/cmd/...
	rm $(BIN)

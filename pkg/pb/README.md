# pb
All proto and generated go files

# How to use?

## Install protoc compiler
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip
### Install unzip
apt update
apt install unzip
### extract and install protoc compiler
unzip protoc-3.13.0-linux-x86_64.zip
cp bin/protoc /usr/local/bin

## Install grpc-go plugin

https://grpc.io/docs/languages/go/quickstart/

$ export GO111MODULE=on  # Enable module mode
go get github.com/golang/protobuf/protoc-gen-go  google.golang.org/grpc/cmd/protoc-gen-go-grpc
make sure gopath/bin is in your bash PATH

## Build USP Message and Record protocol buffer go files
cd usp/usp_msg
./gengo.sh

cd usp/usp_record
./gengo.sh

## Build cwmpGrpc and mtpgrpc grpc go files
cd mtpgrpc
./gengo.sh

cd cwmpGrpc
./gengo.sh



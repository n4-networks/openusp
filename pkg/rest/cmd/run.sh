#!/bin/bash

## External Server addresses
DB_SERVER_ADDR="172.17.0.4:27017"
MTP_SERVER_ADDR=":9001"

## REST configuration
SERVER_BIN=rest
SERVER_PORT=8081
DB_USERNAME=n4admin
DB_PASSWD=n4defaultpass
LOGGING=all


ENV="HTTP_PORT=$SERVER_PORT MTP_GRPC_ADDR=$MTP_SERVER_ADDR DB_ADDR=$DB_SERVER_ADDR DB_USER=$DB_USERNAME DB_PASSWD=$DB_PASSWD LOGGING=$LOGGING"

ENV_TLS=$ENV HTTP_TLS=1

#echo "$ENV ./$SERVER_BIN"
eval "$ENV ./$SERVER_BIN"


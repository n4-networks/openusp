#!/bin/bash
# Copyright 2023 N4-Networks.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


## External Server addresses
DB_SERVER_ADDR="10.5.0.2:27017"
MTP_SERVER_ADDR=":9001"

## REST configuration
SERVER_BIN=apiserver
SERVER_PORT=8081
DB_USERNAME=admin
DB_PASSWD=admin
LOGGING=all


ENV="HTTP_PORT=$SERVER_PORT MTP_GRPC_ADDR=$MTP_SERVER_ADDR DB_ADDR=$DB_SERVER_ADDR DB_USER=$DB_USERNAME DB_PASSWD=$DB_PASSWD LOGGING=$LOGGING"

ENV_TLS=$ENV HTTP_TLS=1

#echo "$ENV ./$SERVER_BIN"
eval "$ENV ./$SERVER_BIN"


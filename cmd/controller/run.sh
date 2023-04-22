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


#ACTIVEMQ=172.21.0.3:61613
ACTIVEMQ=10.5.0.3:61613
MONGO=10.5.0.2:27017
REDIS=10.5.0.4:6379
MQTT=10.5.0.3:1883
DBUSER=admin
DBPASSWD=admin
AGENT=os::012345-525400C8712E
LOGGING=all

AGENT_ID=$AGENT STOMP_ADDR=$ACTIVEMQ MQTT_ADDR=$MQTT DB_ADDR=$MONGO DB_USER=$DBUSER DB_PASSWD=$DBPASSWD CACHE_ADDR=$REDIS ./controller -c


#!/bin/bash

#ACTIVEMQ=172.21.0.3:61613
ACTIVEMQ=10.5.0.3:61613
MONGO=10.5.0.2:27017
REDIS=10.5.0.4:6379
MQTT=10.5.0.3:1883
DBUSER=admin
DBPASSWD=admin
#AGENT=os::N4est-02:42:ac:11:00:06
AGENT=os::012345-525400C8712E
LOGGING=all

AGENT_ID=$AGENT STOMP_ADDR=$ACTIVEMQ MQTT_ADDR=$MQTT DB_ADDR=$MONGO DB_USER=$DBUSER DB_PASSWD=$DBPASSWD CACHE_ADDR=$REDIS ./server -c


#!/bin/bash

#AGENT=os::N4est-02:42:ac:11:00:06
#AGENT=proto::rx_usp_agent_mqtt
AGENT=os::012345-525400C8712E
LOGGING=all
REST=http://localhost:8081
#REST=https://172.17.0.2:8081

AGENT_ID=$AGENT REST_SRV_ADDR=$REST CLI_LOGGING=$LOGGING ./cli


# Message Transfer Protocol Manager
N4 Networks's MTP Manager is an architectural block of USP Controller. This handles all transport related data and message processing. This cloud native application consists of one or many containers managed by Kubernetes deployment framework.

## Introduction
USP (User Service Platform) supports various message transfer protocols and mechanism to create communication paths between Controller and Agent or two Endpoints. These protocols are mainly at application level and use UDP or TCP as transport protocols underneath.

## Protocol Binding Mechanisms
While USP supports message brokers and various other client server protocols to cater to the devices having small form-factors (e.g. IoT devices etc.) N4 Controller does support TR-069 based network devices. The default MTP of TR-069 is HTTP/HTTPS.

* CoAP: Constrained Application Protocol
* STOMP: Simple Text-Oriented Messaging Protocol
* MQTT: MQ Telemetry Transport
* WebSocket: Self explainatory
* HTTP: Hyper Text Transfer Protocol

### CoAP

## Message Encoding and Decoding

## Configuration
Even though it was designed to take configuration files from mtt.yaml in yaml format but few parameters like port and addresses etc. are taken from environment varibles to support ConfigMap kind of deployment

### Environment variables
1. STOMP_ADDR. Ex: STOMP_ADDR=172.17.0.4:61613 where 61613 is the port of activemq server. If not set it would take default as ":61613"
2. DB_ADDR.  Ex: DB_ADDR=172.17.0.3:27017. If not set it would take default as ":27017".
3. AGENT_ID. Ex. AGENT_ID=1234-N4-221133. If not set it and if mtp is running in CLI mode (-c) it would given error

### Run Options (command line opts)
1. MTP is designed to be running a deamon but it does support command line cli. To run in cli mode use -c option (Ex: ./mtp -c).
2. You can provide a config file by giving -f option. Ex. ./mtp -f myconfig.yaml, default is mtp.yaml

### gRPC Support
MTP provides a set of GRPC call to send messages to agent. Please refer to pb/mtpgrpc/mtpgrpc.proto for more details


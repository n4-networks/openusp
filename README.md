[![CI Build Status](https://github.com/n4-networks/openusp/actions/workflows/build.yml/badge.svg)](https://github.com/n4-networks/openusp/actions/workflows/build.yml)

# Summary 
Open source implementation of USP (User Services Platform) controller based on Broadband Forum's [USP](https://usp.technology) Specification.

## Background
So far most of the residential broadband devices e.g. Residential Gateways and Routers are being remotely managed by TR-069 based ACS (Auto configuration server) controllers or by proprietary software.

User Services Platform ([USP](https://usp.technology)) is a standardized framework to remotely manage networked devices (RGs, IoT Devices) through standard protocols as defined by [TR-369a2](https://usp.technology/specification/index.html).

OpenUSP is an Apache Licensed 2.0 repository, written primarily in [golang](https://go.dev). 

# How to use it?
The most simplest way to use is using docker-compose to pull all the required images from docker hub. 

```
git clone git@github.com:/n4-networks/openusp
cd openusp
source scripts/bash/aliases
dc up -d

ubuntu@openusp:~/openusp$ dc up -d
[+] Running 7/7
 ✔ Network openusp               Created                                  0.0s 
 ✔ Container openusp-cache       Started                                  0.0s 
 ✔ Container openusp-db          Healthy                                  0.0s 
 ✔ Container openusp-broker      Started                                  0.0s 
 ✔ Container openusp-controller  Started                                  0.0s 
 ✔ Container openusp-agent       Started                                  0.0s 
 ✔ Container openusp-apiserver   Started                                  0.0s 
ubuntu@openusp:~/openusp$ 
ubuntu@openusp:~/openusp$ cli
OpenUsp-Cli>>  
**************************************************************
                          OpenUsp Cli
**************************************************************
OpenUsp-Cli>> show agent all-ids
Agent Number              : 1           
EndpointId                : os::012345-000000000000
-------------------------------------------------
Agent Number              : 2           
EndpointId                : os::012345-02420A050007
-------------------------------------------------
OpenUsp-Cli>> 
OpenUsp-Cli>> 
OpenUsp-Cli>> 
OpenUsp-Cli>> 
OpenUsp-Cli>> show agent
Object Path                 : Device.LocalAgent.
 UpTime                     : 362         
 SupportedProtocols         : STOMP, CoAP, MQTT, WebSocket
 SoftwareVersion            : 7.0.2       
 EndpointID                 : os::012345-000000000000
 CertificateNumberOfEntries : 0           
 SupportedFingerprintAlgorithms : SHA-1, SHA-224, SHA-256, SHA-384, SHA-512
 ControllerNumberOfEntries  : 1           
 MTPNumberOfEntries         : 1           
 SubscriptionNumberOfEntries : 0           
 RequestNumberOfEntries     : 0           
-------------------------------------------------
OpenUsp-Cli>> 
OpenUsp-Cli>> 

```

For more details visit https://openusp.io



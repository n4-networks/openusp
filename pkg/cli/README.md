# Overview
**C**ommand **L**ine **I**nterface (CLI) N4 is centralized command center to issue various commands to all the distributed subsystems. It interacts with rest of the N4 modules through gRPC protocol and uses database methods to interact with the DB. 

It issues variour USP commands to agents through MTP module.

It also reads datamodel, parameter and instance information of objects from Database and manages the db objects, tables.

# Syntax
A standard syntax of "verb noun" form has been followed for all the CLI commands. 

# Command Summary
```
add instance|wifi <path|ssid> <ssid:securitytype> <radio>
connect to mtp|db <addr:port>
operate command <command> Ex: operate command Device.Reboot()
remove object|collection|instance|wifi 
remove db <object|collection> <objname|collectionname>
set agentid|param|wifi
show agentid|param|datamodel|instance|wifi
update param|datamodel|instance
```

## Command Details
### Add
Add Commands can be used to create new instances of Device2 datamodel objects. Instances of objects having multi-instance capabilities can only be created or removed.
```
add instance <path> <paramname> <paramvalue>
add wifi <ssid>
```
### Connect

### Remove
#### IP
```
remove ip intf <id|name> [ipv4addr|ipv6addr] [id|name]
remove ip intf <id|name>
remove ip intf *
remove ip intf
remove ip intf <id|name> ipv4addr <id|name>
remove ip intf <id|name> ipv4addr
remove ip intf <id|name> ipv4addr *
```




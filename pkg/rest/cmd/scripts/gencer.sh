#!/bin/sh
openssl req -new -nodes -newkey rsa:4096 -days 365 -x509 -keyout ../ssl/server.key -out ../ssl/server.csr -subj "/C=IN/ST=Telengana/L=Hyderabad/O=N4-Networks/OU=IT/CN=n4-networks.com"

#/C=Country
#/ST=State
#/L=Location
#/O=Organization
#/OU=Organization Unit
#/CN=Common Name





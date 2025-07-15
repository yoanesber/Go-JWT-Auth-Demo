#!/usr/bin/bash

# Create directory for certificates if it doesn't exist
mkdir -p cert
cd cert/

# ip=$(ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
# ip=$(hostname -I | awk '{print $1}')
ip=localhost
echo $ip

# Generate a private key for the Certificate Authority (CA) with a key length of 2048 bits
openssl genrsa -out myCA.key 2048

# Create a self-signed certificate for the Certificate Authority (CA) using the private key 'myCA.key'
# The certificate is valid for 730 days (2 years) and uses the Common Name (CN) as the IP address (or hostname)
# The -x509 option indicates that we are creating a self-signed certificate
# This certificate is used by the client to verify the server's certificate
openssl req -x509 -new -key myCA.key -out myCA.cer -days 730 -subj /CN=$ip

# Generate a private key for your server's certificate with a key length of 2048 bits
# This key is used to sign the server certificate
openssl genrsa -out mycert.key 2048

# Create a configuration file for Subject Alternative Names (SAN)
# This file specifies the IP address or domain name that will be associated with the certificate
# The SAN (Subject Alternative Name) is used to specify additional host names for the certificate, e.g., DNS:example.com, DNS:www.example.com, DNS:api.example.com
cat > san.cnf <<EOF
[req]
distinguished_name=req
req_extensions=req_ext
prompt=no

[req_ext]
subjectAltName=DNS:$ip
EOF

# Generate a Certificate Signing Request (CSR) for the server's certificate
# The CSR includes the server's public key and is signed with the private key 'mycert.key'
# The CN (Common Name) is set to the IP address, and the SAN configuration is used to include the IP in the request
# The CSR is signed by the CA to create the server certificate
openssl req -new -key mycert.key -out mycert.csr -subj /CN=$ip -config san.cnf

# Sign the CSR ('mycert.csr') with the CA's private key ('myCA.key') to generate the server certificate ('mycert.cer')
# This command creates a signed certificate for the server, valid for 365 days (1 year), using the CA's certificate ('myCA.cer')
# The SAN extension is also included, allowing the certificate to recognize the IP (or domain name)
openssl x509 -req -in mycert.csr -CA myCA.cer -CAkey myCA.key -CAcreateserial -out mycert.cer -days 365 -extfile san.cnf -extensions req_ext

# Cleanup file konfigurasi sementara
rm san.cnf mycert.csr

cd ../
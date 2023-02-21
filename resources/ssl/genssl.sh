#!/usr/bin/env bash
set -ex

: ${1?'missing key directory'}

key_dir="$1"

chmod 0700 "$key_dir"
cd "$key_dir"

cat >server.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no
[req_distinguished_name]
CN = luojun96.k8s-admission-webhook-server.svc
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = luojun96.k8s-admission-webhook-server.svc
EOF

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Admission Controller Webhook CA"
# Generate the private key for the webhook server
openssl genrsa -out k8s-admission-webhook-server.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with private key of the CA.
openssl req -new -key k8s-admission-webhook-server.key -subj "/CN=luojun96.k8s-admission-webhook-server.svc" -config server.conf \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out k8s-admission-webhook-server.crt -extensions v3_req -extfile server.conf

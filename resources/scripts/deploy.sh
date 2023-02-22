#!/usr/bin/env bash

set -euo pipefail

basedir="$(dirname "$0")"
keydir="$(mktemp -d)"

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
"${basedir}/ssl/genssl.sh" "$keydir"

echo "Creating Kubernetes objects..."
# Create the `admission-webhook` namespace. This cannot be part of the YAML file as we first need to create The TLS secret.

kubectl create namespace admission-webhook

# Create the TLS secret for the generated keys.
kubectl -n admission-webhook create secret tls admission-webhook-server-tls \
    --cert "${keydir}/k8s-admission-webhook-server.crt" \
    --key "${keydir}/k8s-admission-webhook-server.key"

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML template with it.
# Then, create the Kubernetes resources.
ca_pem_b64="$(openssl base64 -A <"${keydir}/ca.crt")"
sed -e 's/${CA_PEM_B64}/'"$ca_pem_b64"'/g' <"${basedir}/deployment.yaml.template" \
    | kubectl create -f -

# Delete the key directory
rm -rf "$keydir"

echo "The admission webhook server has been deployed and configured."

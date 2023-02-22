#!/usr/bin/env bash

set -euo pipefail

basedir="$(dirname "$0")"
keydir="$(mktemp -d)"

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
"${basedir}/ssl/genssl.sh" "$keydir"



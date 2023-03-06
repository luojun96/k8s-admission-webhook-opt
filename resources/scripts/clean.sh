#!/usr/bin/env bash

set -euo pipefail

kubectl delete namespace admission-webhook

kubectl delete mutatingwebhookconfigurations admission-webhook-opt

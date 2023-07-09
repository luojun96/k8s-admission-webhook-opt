# A Simple Admission Controller Webhook

## Introduction

This admission controller webhook is used to configure the security context for pods that are created in non Kubernetes-owned namespaces (`kube-system` and `kube-public`). The detail logic of the webhook is that for every pod that is created (outside of Kubernetes namespaces), it first checks if `runAsNonRoot` is set. If it is not, it is set to a default value of `false`. Furthermore, if `runAsUser` is not set (and `runAsNonRoot` was not initially set), it defaults `runAsUser` to a value of 1234.

## Getting Started

### Build

```bash
make 
```

### Deploy

```bash
./resources/scripts/deploy.sh
```

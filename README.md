# Operator SDK Demo: Pod Notification Restarter

A demo project showcasing a Kubernetes Operator built using the Operator SDK. This operator monitors pod restarts in a Kubernetes cluster and sends notifications to Discord, Slack, or Teams when a pod exceeds a specified restart threshold.

## Introduction

This Operator SDK-based project demonstrates how to:

Monitor Kubernetes pods for restart events.
Trigger notifications using Discord, Slack, or Microsoft Teams webhooks.
Use a custom resource definition (CRD) to configure notification thresholds.

## Quick Start

1. **Clone the repository:**

```bash
git clone https://github.com/NoNickeD/operator-sdk-demo.git

cd operator-sdk-demo
```

2. **Deploy a local Kubernetes cluster:**

```bash
cd ./Infrastructure

task full-deploy-local
```

3. **Confirm the cluster is running:**

```bash
kubectl get nodes

kubectl get pods -A
```

4. **Build and deploy the operator:**

```bash
cd ..

task build-local

make deploy IMG=ttl.sh/$(uuidgen):2h
```

5. **Apply the sample CR and restart deployment:**

```bash
cd ./app
kubectl apply --filename podnotifrestart-sample.yaml

kubectl apply --filename restart-deployment.yaml
```

6. **Verify the custom resource and pods:**

```bash
kubectl get podnotifrestarts

kubectl get pods -A
```

![Discord Notifications](./img/discord_notif.png)

## Creating the Project

The project was created with the following commands (already executed; do not repeat):

```bash
operator-sdk init --domain=vodafone.com --repo=github.com/NoNickeD/operator-sdk-demo

operator-sdk create api --group notifier --version v1alpha1 --kind PodNotifRestart --resource --controller

# Add your code and then:
make generate

make manifests
```

### Configuring Webhooks

Add your webhook URLs to the operator's environment variables in `config/manager/manager.yaml`. Example for Discord:

```bash
- name: DISCORD_WEBHOOK_URL
  value: "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID"
```

## Clean Up

To delete the local Kubernetes cluster:

```bash

cd ../Infrastructure

task delete-local-cluster
```

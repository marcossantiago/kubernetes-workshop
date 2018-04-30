---
showLeftCol: 'false'
hideFirstStyle: 'false'
showFooterText : 'true'
title: Production Grade Kubernetes
---

## Workshop Setup

---

We have provided everything you need to participate in this workshop.

We have provisioned you a Virtual Machine running a Kubernetes cluster.

The following steps will get you access to your cluster.

---

### Connect to your Virtual Machine

Open a Terminal window and login to the VM via SSH

```
$ ssh csuser@[VIRTUAL MACHINE IP]
```

Enter the password provided to you.

---

### Retrieve configuration files

Retrieve the permissions file

```
$ wget https://storage.googleapis.com/kubecon-2018-pgk/ca.pem
```

---

## Configure your Kubernetes cluster

Execute the following commands to setup your cluster

---

## Step 1: Point to the Cluster

```bash
$ kubectl config set-cluster workshop \
  --server=https://35.195.196.5 \
  --certificate-authority=ca.pem
```

---

## Step 2: Authenticate

```bash
$ kubectl config set-credentials workshop-user \
  --username=[CLUSTER USERNAME] \
  --password=[CLUSTER PASSWORD]
```

---

## Step 3: Set Context

```bash
$ kubectl config set-context workshop \
  --cluster=workshop \
  --user=workshop-user \
  --namespace=[CLUSTER USERNAME]
```

---

## Step 4: Connect

```bash
$ kubectl config use-context workshop
```

---

### Verify the Configuration

Check the cluster is up and running correctly by retrieving information about it

```
$ kubectl cluster-info

Kubernetes master is running at https://...
....
```

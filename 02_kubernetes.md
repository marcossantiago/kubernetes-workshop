---
showLeftCol: 'false'
hideFirstStyle: 'false'
showFooterText : 'true'
title: Production Grade Kubernetes
---

# Kubernetes

---

## Background

* Open-source container orchestrator from Google
* Now part of Cloud Native Computing Foundation
* Popular and Active: >32K stars on Github

---

## Features

 - Based on Google's experience running containers (Borg)
 - Bakes in various features
   - Load-balancing, secret management
 - Opinionated
   - Impact on application design

---

# The Kubernetes Cluster

---

<img src="img/30. High-level Kubernetes architecture diagram-01.png" alt="Kubernetes Architecture">

---

## Core Components

 - Nodes
 - Pods
 - Labels & Selectors
 - Services
 - Deployments
 - ReplicaSets
 - Namespaces
 - Jobs

---

## Nodes

 * Worker machine
 * May be a VM or physical machine
 * A Node can host one or multiple Pods
 * Include Container runtime, kubelet and kube-proxy

---

<img src="img/32. Kubernetes Node diagram-01.png" alt="Kubernetes Node">

---

## Pods

 - A group of one or more running containers
 - Always co-located, co-scheduled, and run in a shared context
 - Atomic unit (scaling and deployment)
 - Containers in a Pod share an IP address
 - Containers share resources within the pod
 - Single container pods are common
 - Pods are ephemeral

---

<img src="img/33. Kubernetes Pods diagram-01.png" alt="Kubernetes Pod">

---

## Kubernetes Configuration

 - Use configuration files to manage resources
 - Specified in YAML or JSON
  - YAML tends to be more user-friendly
 - Made up of 3 parts:
  - API, Metadata, Spec
 - Often combined
 - Should be stored in version control

---

## Pod Configuration

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.7.9
    ports:
    - containerPort: 80
```

---

## Labels

 - Key/Value pairs attached to objects
    - e.g: "version: dev", "tier: frontend"
 - Objects include Pods, ReplicaSets, Services
 - Label selectors then used to identify groups
 - Used for load-balancing etc...

---

<img src="img/35. Kubernetes Labels diagram-01.png" alt="Kubernetes Labels">

---

## Pod with Labels

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    app: workshop-demo
    tier: front
spec:
  containers:
  - name: nginx
    image: nginx:1.7.9
    ports:
    - containerPort: 80
```

---

## Selectors

 - Used to query Labels
   - environment = production
   - tier != frontend
 - Also set based comparisons
   - environment in (production, staging)
   - tier notin (frontend, backend)

---

<img src="img/35. Kubernetes Labels diagram-01.png" alt="Kubernetes Selectors">

---

## Services

 - Stable endpoints addressed by name
 - Forward traffic to Pods
 - Pods are selected by Labels
 - Round-robin load-balancing
 - Separates endpoint from implementation

---

<img src="img/37. Kubernetes Services diagram_A-01.png" alt="Kubernetes Services 1">

---

<img src="img/37. Kubernetes Services diagram-01.png" alt="Kubernetes Services 2">

---

<img src="img/36. Kubernetes Selectors diagram-01.png" alt="Kubernetes Services 3">

---

## Service Configuration

```
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  ports:
  - protocol: TCP
    port: 8000 # Listen on inside the cluster
    targetPort: 80 # Forward to port on endpoints
  selector:
    app: nginx
```

---

## Types of Service

* ClusterIP (default)
* NodePort
* LoadBalancer
* ExternalName

---

## ClusterIP

 * Uses internal IP for service
 * No external exposure

---

<img src="img/41. Cluster IP Service diagram-01.png" alt="Kubernetes ClusterIP">

---

## NodePort

 * Service is externally exposed via port on host
 * Same port on every host
 * Port automatically chosen or can be specified

---

<img src="img/43. NodePort service diagram-01.png" alt="Kubernetes NodePort">

---

## LoadBalancer

 * Exposes service externally
 * Implementation dependent on cloud provider

---

<img src="img/Load balancer-01.png" alt="Kubernetes LoadBalancer">

---

## ExternalName

 * For forwarding to resources outside of Kubernetes
   * e.g. Existing database
 * CNAME created

---

## Flat networking space

 - All pods, across all hosts, are in the same network space
   - Can see each other without NAT
 - Simple cross host communication

---

<img src="img/34. Flat Networking Space-01.png" alt="Kubernetes Networking">

---

## Deployments

 * Declarative state of your Pods and ReplicaSets
 * Handles Rollout, Rollback & Updates

---

<img src="img/45. Kubernetes Deployments diagram-01.png" alt="Kubernetes Deployments">

---

## ReplicaSets

 * Explicit definition of Desired State
 * Monitors the status of Pods
   * Too few? Start a new one
   * Too many? Kill some

---

<img src="img/Deployment Replicaset.png" alt="Kubernetes Deployments">

---

## Namespaces

 - Resources can be partitioned into Namespaces
 - Logical groups
 - System resources run in their own Namespace
 - Normally only use one Namespace

---

<img src="img/47. Kubernetes Namespaces diagram-01.png" alt="Kubernetes Namespaces">

---

## Jobs

 - Typically for performing batch processing
 - Spins up short-lived pods
 - Ensures given number run to completion

---

### More & More

 - Volumes
 - Stateful Sets
 - Ingress
 - Annotations
 - Daemon Sets
 - Network Policies
 - Security Context
 - Service Accounts
 - ...

---

## Dashboard

* Simple Web User Interface
* Good *high-level* overview of the cluster
* Can drill down into details
* Useful for debugging

---

<img src="img/dashboard.png" alt="Kubernetes Dashboard">

---

## Get Hands on!

---

### Introducting Kubectl

`kubectl` is the command line interface (CLI) tool for sending commands to a Kubernetes cluster.

We will use this tool to deploy, view, and access an application on our cluster.

---

## Step 1: Kubectl Basics

* The format of a kubectl command is:
```
kubectl [action] [resource]
```
* This performs the specified action  (like `create`, `describe`) on the specified resource (like `node`, `container`).
* Use `--help` after the command to get additional info about possible parameters
```
$ kubectl get nodes --help
```

---

Check that kubectl is configured to talk to your cluster, by running the kubectl version command:
```bash
$ kubectl version
```

You can see both the client and the server versions.

---

To view how to reach the cluster, run the `cluster-info` command:

```
$ kubectl cluster-info

Kubernetes master is running at https://35.205.211.112
GLBCDefaultBackend is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/default-http-backend:http/proxy
Heapster is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/heapster/proxy
KubeDNS is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
KubeRegistry is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/kube-registry:registry/proxy
kubernetes-dashboard is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy
Metrics-server is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/https:metrics-server:/proxy
Grafana is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/monitoring-grafana/proxy
InfluxDB is running at https://35.205.211.112/api/v1/namespaces/kube-system/services/monitoring-influxdb:http/proxy
```

To further debug and diagnose cluster problems, use:
```
kubectl cluster-info dump
```

---

## Step 2: Deploy an Application

Let’s run our first application on Kubernetes with the kubectl run command. The `run` command creates a new deployment for the specified container. This is the simplest way of deploying a container.

```bash
$ kubectl run hello-kubernetes \
--image=gcr.io/google_containers/echoserver:1.4 --port=8080

deployment "hello-kubernetes" created
```

---

This performed a few things:
* Searched for a suitable node.
* Scheduled the application to run on that node.
* Configured the cluster to reschedule the instance on a new node when needed.

---

List your deployments

```bash
$ kubectl get deployments

NAME        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE
hello-kubernetes   1         1         1            1
```

We see that there is 1 deployment running a single instance of your app.

---

Gather information about the status of your objects (pods, deployments, services, etc) using

```
kubectl get <object>
```
and
```
kubectl describe <object>
```

---

## Step 3: Make the App Visible

By default applications are only visible inside the cluster. We can create a proxy to connect to our application.

Start by finding out the pod name:
```
$ kubectl get pod

NAME                               READY     STATUS    RESTARTS
hello-kubernetes-624527933-nth9d   1/1       Running   0
```

---

Create a proxy for the pod

```bash
$ kubectl port-forward <POD NAME> 8080 &
```
We now have a connection between our host and the Kubernetes cluster.

---

## Step 4: Access the App

To see the output of our application, run a curl request to the local port:
```bash
$ curl http://localhost:8080

CLIENT VALUES:
client_address=127.0.0.1
command=GET
real path=/
query=nil
request_version=1.1
request_uri=http://0.0.0.0:8080/

SERVER VALUES:
server_version=nginx: 1.10.0 - lua: 10001

HEADERS RECEIVED:
accept=*/*
host=0.0.0.0:8080
user-agent=curl/7.51.0
BODY:
-no body in request-
```

---

## Step 5: Clean Up

`port-forward` is meant for testing services that are not exposed. To expose the application, use a Service (covered later).

Kill port forward
```
$ kill %2
```

Delete old Deployment
```
$ kubectl delete deployment hello-kubernetes
```

---

## Step 6: Create a new Deployment & Service

```
$ kubectl run hello --image=gcr.io/google_containers/echoserver:1.4 \
   --port=8080 \
   --expose \
   --service-overrides='{ "spec": { "type": "NodePort" } }'

service "hello" created
deployment "hello" created
```

This creates a new Deployment and Service of type:NodePort. A random high port will be allocated to which we can connect.

---

View the Service

```
$ kubectl get service

NAME               CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
hello         10.0.122.112   <nodes>       8080:30659/TCP   10m
```

We can see the port on which it is exposed, but what is the external IP?

---

To find the IP on which to call we need information on the nodes (use the EXTERNAL-IPs from any node):

```
$ kubectl get nodes -o wide

NAME                           STATUS                     AGE       VERSION   EXTERNAL-IP      OS-IMAGE                             KERNEL-VERSION
kubernetes-master              Ready,SchedulingDisabled   17m       v1.7.5    35.187.38.163    Container-Optimized OS from Google   4.4.52+
kubernetes-minion-group-c9bz   Ready                      17m       v1.7.5    35.189.206.159   Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
kubernetes-minion-group-cfzx   Ready                      17m       v1.7.5    35.195.36.237    Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
kubernetes-minion-group-ftw1   Ready                      17m       v1.7.5    35.195.61.242    Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
```

---

Access the external IP with Curl:

```
$ curl 35.189.206.159:30659

CLIENT VALUES:
client_address=10.132.0.3
command=GET
real path=/
query=nil
request_version=1.1
request_uri=http://35.187.76.71:8080/

SERVER VALUES:
server_version=nginx: 1.10.0 - lua: 10001

HEADERS RECEIVED:
accept=*/*
host=35.187.76.71:8080
user-agent=curl/7.52.1
BODY:
-no body in request-
```

---

## Step 7: Clean Up

Delete the Deployment

```
$ kubectl delete deploy hello
```

Delete the Service

```
$ kubectl delete svc hello
```

---

## What have we Learned?

* The history of Container Orchestration.
* Basics of Kubernetes.
* How to deploy a simple application on to our own cluster.

---

[Next up, a real application!](../03_real_application.md)


![kubernetes](/img/1. kubernetes-logo.png)

The contents of this workshop have been borrowed from the internet. 

---

### Contents

* What is Kubernetes ?
* Core Concepts
* Kubernetes Commands
* [Hands On! ***Level 101***](/#/introducting-kubectl)
* [Real application](/#/a-real-application)
* [Next resources](/#/what-next-)

---

# What is Kubernetes?

<blockquote>
 Kubernetes is a portable container platform for managing containerized workloads and services, that facilitates both declarative **configuration** and **automation**. 
</blockquote> <!-- .element: class="fragment" data-fragment-index="2" -->

Note:
Plataform for allowing you to maintain deploying containers into production
Automated Health Checks, Rolling Deployments and restarts,

---

# What is not?

Does not: <!-- .element: class="fragment" data-fragment-index="2" -->
* limit the type of applicattion supported <!-- .element: class="fragment" data-fragment-index="3" -->
* deploy source code and does not build your application <!-- .element: class="fragment" data-fragment-index="4" -->
* provide application-level services: middleware, data-processing frameworks, databases, caches, nor cluster storage systems as built-in services. <!-- .element: class="fragment" data-fragment-index="5" -->
* dictate logging, monitoring, or alerting solutions. <!-- .element: class="fragment" data-fragment-index="6" -->

---

# What can it do?

* Deployment <!-- .element: class="fragment" data-fragment-index="2" -->
* Auto Scaling <!-- .element: class="fragment" data-fragment-index="2" -->
* Scheduling <!-- .element: class="fragment" data-fragment-index="2" -->
* Monitoring <!-- .element: class="fragment" data-fragment-index="2" -->
* Load Balancing <!-- .element: class="fragment" data-fragment-index="2" -->
* Secret management <!-- .element: class="fragment" data-fragment-index="2" -->
* It orchestrates computing, networking, and storage infrastructure <!-- .element: class="fragment" data-fragment-index="3" -->

---

# What is Orchestration?

<blockquote>
"The planning or coordination of the elements of a situation to produce a desired
effect."
</blockquote><!-- .element: class="fragment" data-fragment-index="2" -->

Oxford English Dictionary<!-- .element: class="fragment" data-fragment-index="2" -->

---

## The Elements

 - Containers
 - Hosts
 - Networking

---

## The Desired Effect

 - Running application
 - Automatically scaling
 - Fault tolerance
   - e.g. Failover, node re-balancing, health checks
 - Efficient use of resources
 - Little manual intervention

---

# One word

*Desired Stage Manager* <!-- .element: class="fragment" data-fragment-index="2" -->

---

# Container Orchestrators

---

## Common Components

* Container Runtime
* Resource Manager
* Scheduler
* Service Discovery
* Advanced Networking

---

## Many Options

* Kubernetes
* ECS 
* Mesos, DC/OS
* Docker Swarm
* Plus others
  * Nomad
  * Fleet from CoreOS (no more)
  * PaaSs...

---

# Kubernetes

---

## Background

* Open-source container orchestrator from Google
* Now part of Cloud Native Computing Foundation
* Popular and Active: >36K stars on Github

---

## Features

* Based on Google's experience running containers
* Bakes in various features
  * Load-balancing, secret management...
* Impact on application design

---

## Components

* Nodes
* Pods
* Labels & Selectors
* Services
* ReplicaSets
* Deployments
* Jobs
* Namespaces

---

<img src="img/30. High-level Kubernetes architecture diagram-01.png">

---

## Nodes

<img src="img/32. Kubernetes Node diagram-01.png">

---

 * Worker machine
 * May be a VM or physical machine
 * A Node can host one or multiple Pods
 * Include Container runtime, kubelet and kube-proxy

Note:
 kube-proxy is responsible for implementing a form of virtual IP for Services of type other than ExternalName

---

## Pods

<img src="img/33. Kubernetes Pods diagram-01.png">

---

* Groups of containers deployed and scheduled together
* Atomic unit (scaling and deployment)
* Containers in a pod share IP address
* Single container pods are common
* Pods are ephemeral

---

## Flat networking space

<img src="img/34. Flat Networking Space-01.png">

---

* All pods, across all hosts, are in the same network space
  * Can see each other without NAT
* Simple cross host communication

---

## Labels

<img src="img/35. Kubernetes Labels diagram-01.png">

---

* Key/Value pairs attached to objects
   * e.g: "version: dev", "tier: frontend"
* Objects include Pods, ReplicaSets, Services
* Label selectors then used to identify groups
* Used for load-balancing etc

---

## Selectors

<img src="img/36. Kubernetes Selectors diagram-01.png">

---

* Used to query labels
  * environment = production
  * tier != frontend
* Also set based comparisons
  * environment in (production, staging)
  * tier notin (frontend, backend)

---

## Services

* Stable endpoints addressed by name
* Forward traffic to Pods
* Pods are selected by Labels
* Round-robin load-balancing
* Separates endpoint from implementation

---

<img src="img/37. Kubernetes Services diagram_A-01.png">

---

<img src="img/37. Kubernetes Services diagram-01.png">

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

<img src="img/41. Cluster IP Service diagram-01.png">

---

## NodePort

 * Service is externally exposed via port on host
 * Same port on every host
 * Port automatically chosen or can be specified

---

<img src="img/43. NodePort service diagram-01.png">

---

## LoadBalancer

 * Exposes service externally
 * Implementation dependent on cloud provider

---

<img src="img/41. Cluster IP Service diagram-01.png">

---

## ExternalName

 * For forwarding to resources outside of Kubernetes
   * e.g. Existing database
 * CNAME created

---

## Deployments

* Deployments start ReplicaSets
* Rollout/Rollback & Updates

---

<img src="img/45. Kubernetes Deployments diagram-01.png" alt="Deployments">

---

## ReplicaSets

* Monitor the status of Pods
  * define number of pods to run
  * start/stop pods as needed

---

## Jobs

* Typically for performing batch processing
* Spins up short-lived pods
* Ensures given number run to completion

---

## Namespaces

* Resources can be partitioned into Namespaces
* Logical groups
* System resources run in their own Namespace
* Normally only use one Namespace

---

<img src="img/47. Kubernetes Namespaces diagram-01.png" alt="Namespaces">

---

### More & More

* Volumes
* Stateful Sets
* Ingress
* Annotations
* Daemon Sets
* Horizontal Pod Autoscaling
* Network Policies
* Resource Quotas
* Secrets
* Security Context
* Service Accounts
* ...

---

## Dashboard

* Simple Web User Interface
* Good *high-level* overview of the cluster
* Can drill down into details
* Useful for debugging

---

<img src="img/dashboard.png">

---

# Kubernetes Configuration

---

## Configuring a Cluster

* Use configuration files to manage resources
* Specified in YAML or JSON
 * YAML tends to be more user-friendly
* Can be combined
* Should be stored in version control

---

## Pod Example

```
apiVersion: v1
   kind: Pod
   metadata:
     name: hello-node
     labels:
       app: hello-node
   spec:
     containers:
       - name: hello-node
         image: hello-node:v1
         ports:
           - containerPort: 8080
```

---

## Service Example

```
apiVersion: v1
kind: Service
metadata:
  name: hello-service
spec:
  type: NodePort
  selector:
    app: hello-node
  ports:
    - name: http
      nodePort: 36000
      targetPort: 8080
      port: 80
      protocol: TCP
```

---

## Hands on!
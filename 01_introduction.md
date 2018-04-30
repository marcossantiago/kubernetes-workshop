---
showLeftCol: 'false'
hideFirstStyle: 'false'
showFooterText : 'true'
title: Production Grade Kubernetes
---

# A Brief History of
# Container Orchestration

---

# In the Beginning...

---

<img src="img/ops_problem.jpeg" alt="Ops Problem" width="150%">

---

<img src="img/devops-wall.png" alt="Wall of Confusion" width="150%">

---

## DevOps

 * Acknowledges Dev and Ops are on the same team
 * Attempts to tear down the Wall
 * **Teams** become responsible for running services

---

<img src="img/ci_manual_deployment.jpg" alt="Manual Deployment" width="200%">

---

<img src="img/continous_delivery.png" alt="Continuous Delivery">

---

# Continuous Delivery

* Build, Test and ship code all with the push of a button
* Development becomes more efficient, reliable and continuous
* Enable developers to safely push new features at a rapid pace
* Removes risk from the build process

---

<img src="img/microservices_mem.jpg" alt="Microservices" width="150%">

---

<img src="img/13. Abstract diagram to represent a MicroServices architecture-01.png" alt="Microservices">

---

# Microservices

* Architectural ideas from lightweight Service Oriented Architecture
 * services are small - fine-grained to perform a single function
 * talk over uniform APIs
* Organisational approaches like DevOps and Agile
 * embrace automation of deployment and testing
 * easing the burden on management and operations
* New technology like Containers and Programmable Infrastructure

---

<img src="img/containers_meme.jpeg" alt="Containers Everywhere" width="120%">

---

<img src="img/containers.png" alt="Containers" width="60%">

---

# Containers

- Enables separation of applications from infrastructure
- Allows developers to package up an application and ship it all out as one
- **Developers** can focus on writing code
- **Operations** get flexibility and a (potential) reduction in the number of systems needed

---

<img src="img/cloud_native.png" alt="Cloud Native" width="50%">

---

## Cloud Native

* Architect applications for the Cloud - **First**
* Delivery - Minimise cycle-time, Automate deployment
* Performance - Responsiveness, Concurrency, Efficiency
* Resilience - Fault-tolerance, Self-healing
* Elasticity - Automatic scaling
* Diagnosability - Logs, Traces & Metrics
* **Assume something will go wrong**

---

# Mission Accomplished?

---

## Co-ordination Challenges

 * DevOps everything
 * Code flowing from Dev to Production
 * A monolith split into dozens and dozens and dozens of pieces..
 * Containers by the boatload
 * **How do we manage all of these?**

---

# Orchestration

---

<blockquote>
"The planning or coordination of the elements of a situation to produce a desired
effect, especially surreptitiously"
</blockquote>

<small>Oxford English Dictionary</small>

---

<blockquote>
"The planning or coordination of **the elements of a situation** to produce a desired
effect, especially surreptitiously"
</blockquote>

---

## The Elements

 - Containers
 - Hosts
 - Networking

---

<blockquote>
"The planning or coordination of the elements of a situation to **produce a desired
effect**, especially surreptitiously"
</blockquote>

---

## The Desired Effect

 - Running application
 - Automatically scaling
 - Fault-tolerance
   - e.g. Failover, node re-balancing, health checks
 - Efficient use of resources
 - Little manual intervention

---

<blockquote>
"The planning or coordination of the elements of a situation to produce a desired
effect, **especially surreptitiously**"
</blockquote>

---

## Surreptitiously

 - Should happen in the background
 - Complexity is hidden
 - User doesn't need to know the details

---

# Container Orchestrators

---

## Common Components

 - Container Runtime
  - e.g. Docker
 - Resource Manager
  - Manages CPU, Memory and Disk
 - Workload Scheduler
  - Makes sure containers run in the desired number on hosts
 - Service Discovery
  - Allow containers to be found by other containers or services
 - Networking
  - e.g. flannel, Weave, Calico, Pipeworks...

---

## Many Options

 - Kubernetes
 - Mesos, DC/OS
 - Docker Swarm
 - Plus others
   - Nomad
   - Fleet from CoreOS (no more)
   - PaaSs...

---

[Tell me more about Kubernetes](../02_kubernetes.md)

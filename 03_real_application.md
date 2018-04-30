---
showLeftCol: 'false'
hideFirstStyle: 'false'
showFooterText : 'true'
title: Production Grade Kubernetes
---

# A Real Application

---

## In this section we will

* Work with a non-trivial application
* Create a Deployment
* Deploy the application on your cluster
* Scale the application
* Create a Service
* Expose the application on your cluster

---

## Our Demo Application

This application is composed of multiple pieces:

- One main back-end service.
- A front-end (UI) service.
- A data layer.

We will deploy these pieces one at a time on our cluster.

---

<img src="img/2. Demo App-01.png">

---

## Demo application

Clone the demo application's repository to your VM

```bash
$ git clone https://github.com/ContainerSolutions/ws-production-grade-kubernetes.git .
```

---

## Recap of Resource Hierarchy

* **Deployment**: Manages ReplicaSets and defines how updates to Pods should be rolled out

* **ReplicaSet**: Ensures that a specified number of Pods are running at any given time

* **Pod**: A group of one or more containers deployed and scheduled together

---

<img src="img/45. Kubernetes Deployments diagram-01.png" alt="Deployments">

---

### Deployment Configuration

The "./real-app" directory of the Git repo contains all of the yaml configurations we will need:

```
# ./real-app/deployment-back-end.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: pgk-back-deployment
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: real-app
        tier: backend
    spec:
      containers:
        - name: production-grade-be
          image: icrosby/k8s-real-demo:latest
          ports:
            - containerPort: 8080
```

---

### Deploy Back-end to the cluster

```
$ kubectl apply -f real-app/deployment-back-end.yaml
```

---

### View Resource Details

Use the "kubectl get" and "kubectl describe" to view details of the deployed resources:

```
$ kubectl get deployments
$ kubectl get replicasets
$ kubectl get pods
```

```
$ kubectl describe pods <POD-NAME>
```

---

### Interact with a Pod remotely

* Pods get a private IP address by default.
* Cannot be reached from outside the cluster.
* Use `kubectl port-forward` to map a local port to a port inside a pod.


---

```
$ kubectl port-forward <POD-NAME> 8080 &
```

```
$ curl 0.0.0.0:8080

Hello from Container Solutions.
I'm running version 1.0 on ...
```

---

### Clean Up

`port-forward` is meant for testing services that are not exposed. To expose the application, use a Service (covered later).

Kill port forward
```
$ kill %2
```

---

### Scaling Deployments

* Deployments manage ReplicaSets.
* Each deployment is mapped to one active ReplicaSet.
* Use `kubectl get replicasets` to view the current set of replicas.
* `kubectl get deployments` will give us the same info (plus more)

```
$ kubectl get rs

NAME                   DESIRED   CURRENT   READY     AGE
...                    1         1         1         16s
```

---

### Scale up/down the Deployment

```
$ kubectl scale deployments pgk-back-deployment --replicas=2

deployment "pgk-back-deployment" scaled
```

---

### Check the status of the Deployment

Notice the new Pod(s)
```
$ kubectl get pods
```

Look at the `Events` at the bottom

```
$ kubectl describe deployment pgk-back-deployment
```

---

### Fault Tolerance

What happens if we kill one of the Pods?

```
$ kubectl get pods
$ kubectl delete pod <POD-NAME>
```

---

## Debugging

---

### View the logs of a Pod

Use `kubectl logs` to view the logs for the `<POD-NAME>` Pod:

```
$ kubectl logs <POD-NAME>
```

> Use the -f flag and observe what happens.

---

### Run an interactive shell inside a Pod

Execute a shell in a Pod, like in Docker:

```
$ kubectl exec -ti <POD-NAME> /bin/sh
```

---

## Accessing our Application

---

## Reminder

* Pods are ephemeral (no fixed IP)
* Port-forwarding strictly a debugging tool
* Need to be able to scale

---

### Services
* Stable endpoints for Pods.
* Based on Labels and Selectors.

---

### Labels & Selectors
* Label: key/value pair attached to objects (e.g. Pods)
* Selector: Identify and group a set of objects.

---

<img src="img/labels2.svg" height="600">

---

### Service Types

* ClusterIP (Default): Exposes the service on a cluster-internal IP.

* NodePort: Expose the service on a specific port on each node.

* LoadBalancer: Use a loadbalancer from a Cloud Provider. Creates `NodePort` and `ClusterIP`.

* ExternalName: Connect an external service (CNAME) to the cluster.

---

### Service Configuration

Look in "./real-app" folder of the Git repo for the Service configuration.

```
# real-app/service-back-end.yaml
apiVersion: v1
kind: Service
metadata:
  name: pgk-back-service
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: 8080
  selector:
    tier: backend

```

---

### Create the Service

```
$ kubectl apply -f ./real-app/service-back-end.yaml
```

---

### Query the Service

Find the NodePort (via Service) and IP (via Node)

```
$ curl [IP]:[NODE_PORT]
```

---

### Test Load Balancing

Make several calls to the service and notice the different responses.

---

### Updating Deployments

(`RollingUpdate`)

* RollingUpdate is the default strategy.
* Updates Pods one (or a few) at a time.

---

### Common Workflow

* Update the application, and create a new version.
* Build the new image and tag it with the new version, i.e. v2.
* Update the Deployment with the new image

---

### Try It Out

First check the current version running

```
$ curl [EXTERNAL_IP]:[NodePort]

Hello from Container Solutions.
I'm running version 1.0 on ...
```

Next, update the image:

```
$ kubectl set image \
  deployment/pgk-back-deployment k8s-real-demo=icrosby/k8s-real-demo:v2
```

---

## Monitor the Deployment

Check status via

```
kubectl rollout status deployment pgk-back-deployment
```

Now verify the new version

```
$ curl [EXTERNAL_IP]:[NodePort]
```

---

## Now it's your turn

---

## Step 1: Build your own Image

Build your own image and push to Docker Hub.

```
docker build -t [DOCKERHUB_USER]/k8s-real-demo:v1.0.0 .
```

```
docker push [DOCKERHUB_USER]/k8s-real-demo:v1.0.0
```

---

## Step 2: Create a Deployment

* Create a Deployment configuration for your Image.
  * Use the same Image name as above.
* Deploy on the cluster.

---

## Step 3: Scale

* Scale the Deployment to 3 instances.
* Verify the scaling was successful and all instances are getting requests.

---

## Step 4: Update

* Modify the `Dockerfile` to return a different Version.
* Build the Image and tag as `v2`.
* Update the Deployment to use the new tag.
* Verify the new version by making an HTTP request.
* View the logs of the application.

---

## Step 5: Deploy the Front-end

In the `./real-app/` folder you will find configuration files for the front-end (Deployment and Service).

* ./real-app/deployment-front-end.yaml
* ./real-app/service-front-end.yaml

Using these configuration files deploy and expose the application on to the cluster.

```
$ kubectl apply -f ./real-app/deployment-front-end.yaml
```

```
$ kubectl apply -f ./real-app/service-front-end.yaml

```

---

## Step 6: Accessing the Front-end

Find the port on which the front end is exposed (via the Service) and access this in your browser.

```
$ kubectl get svc pgk-front-service
```

---

## Bonus Exercise

---

### Storage

While we would like to ideally run stateless applications on Kubernetes, we will eventually run into the challenge of requiring state within our cluster.

Add some state to our application via the front end and refresh the browser a few times.

---

## What is happening to our entries?

The current state is local only to a particular instance. So the view you get depends on which back end is hit. This is obviously not a production ready setup.

We want the state to be shared across all instances.

---

## Problem #2

Kill one of the pods and see how this affects the front end view.

```
$ kubectl delete pod <POD-NAME>
```

---

## No more state

The state currently shares the same lifetime as the pod. This not an acceptable setup.

How can we solve this?

---

## CockroachDB

An open source 'Cloud Native' SQL database.

We will deploy CockroachDB to maintain the state for our demo application.

---

### Deploying CockroachDB

```
$ kubectl apply -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cockroachdb-statefulset.yaml

service "cockroachdb-public" created
service "cockroachdb" created
poddisruptionbudget "cockroachdb-budget" unchanged
statefulset "cockroachdb" created
```

---

### What did this do?

* Pull the Kubernetes configuration file from Github
* Created two Services
* Created a `poddisruptionbudget`
* Created a `StatefulSet`
  * Three Pods
  * Three PersistentVolumes
  * Three PersistentVolumeClaims

---

Don't worry, we will cover StatefulSets and PersistentVolumes later on

---

### Connecting to the Database

Finally we need to configure our back end to use the database instead of local storage.

Add the following environment variable to the back end.

```
        env:
        - name: SQL_DATASTORE_CONNECTION
          value: "host=cockroachdb-public port=26257 user=root dbname=gorm sslmode=disable"
```

Restart the back end (kill all pods) for this to take effect (we will look at a better way to solve this later on).

---

## What have we Learned?

* How to deploy a 'real world' application on Kubernetes
* Deal with Deployment and Services
* Connecting Services with labels and selectors
* Scale up/down
* Update a Deployment (rolling update)

---

[Next up, heading to Production...](../04_productionize.md)


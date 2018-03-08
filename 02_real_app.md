## Beyond The Basics

---

### In this section we will:

* Work with a non trivial application
* Create a Deployment
* Deploy the application on your cluster
* Scale the application
* Create a Service
* Expose the application on your cluster

---

## 'Real' demo application

We will work the following demo application: `https://github.com/idcrosby/k8s-example`
Clone the repo to your VM:

```
$ git clone https://github.com/idcrosby/k8s-example.git
```

This application is composed of multiple pieces. One main backend service, a front end (UI) service, and a data layer. We will deploy these pieces one at a time on our cluster.

---

**Architecture Diagram**

---

## Recap of Resource Hierarchy

---

### Pod
A Pod is a group of one or more containers deployed and scheduled together.

### ReplicaSet
A ReplicaSet ensures that a specified number of Pods are running at any given time.

### Deployment
A Deployment manages ReplicaSets and defines how updates to Pods should be rolled out.

---

<!-- .slide: data-background="img/deployments-to-containers.png" data-background-size="70%"-->

---

### Creating a Deployment

./resources/deployment.yaml (in the resources folder)

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: k8s-real-demo
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: k8s-real-demo
    spec:
      containers:
      - name: k8s-real-demo
        image: icrosby/k8s-real-demo
        ports:
        - containerPort: 8080
```

---

### Deploy to K8s

```
$ kubectl apply -f resources/deployment.yaml
```

---

### View Resource Details

Use the `kubectl get` and `kubectl describe` to view details for the `k8s-real-demo` resources:

```
$ kubectl get deployments
$ kubectl get replicasets
$ kubectl get pods
```

```
$ kubectl describe pods <pod-name>
...
```

---

### Interact with a Pod remotely

* Pods get a private IP address by default.
* Cannot be reached from outside the cluster.
* Use `kubectl port-forward` to map a local port to a port inside the `k8s-real-demo` pod.


---

### Use two terminals

* Terminal 1

```
$ kubectl port-forward <pod-name> 8080:8080
```

* Terminal 2

```
$ curl 0.0.0.0:8080
Hello from Container Solutions.
I'm running version 1.0 on k8s-real-demo-648d67845-hh8bn
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
k8s-real-demo-364036756   1         1         1         16s
```

---

### Scale up/down the Deployment

```
$ kubectl scale deployments k8s-real-demo --replicas=2
deployment "k8s-real-demo" scaled
```

---

### Check the status of the Deployment

Notice the new pod(s)
```
$ kubectl get pods
```

Look at the `Events` at the bottom

```
$ kubectl describe deployment k8s-real-demo
```

---

### Fault tolerance

What happens if we kill one of the pods?

```
$ kubectl delete po <pod-name>
```

---

## Debugging

### View the logs of a Pod

Use `kubectl logs` to view the logs for the `<pod-name>` Pod:

```
$ kubectl logs <pod-name>
```

> Use the -f flag and observe what happens.

---

### Run an interactive shell inside a Pod

Execute a shell in a Pod, like in Docker:

```
$ kubectl exec -ti <pod-name> /bin/sh
```

---

## How to access our application?

* Port-forwarding strictly a debugging tool
* Pods are ephemal (no fixed IP)
* Need to be able to scale 

---

## Services

---

### Recap of services
* Stable endpoints for Pods.
* Based on labels and selectors.

---

### Labels & Selectors
* Label: key/value pair attached to objects (e.g. Pods)
* Selector: Identify and group a set of objects. 

---

<img src="img/labels2.svg" height="600">

---

### Service types

* `ClusterIP` (Default) Exposes the service on a cluster-internal IP.

* `NodePort` Expose the service on a specific port on each node.

* `LoadBalancer` Use a loadbalancer from a Cloud Provider. Creates `NodePort` and `ClusterIP`.

* `ExternalName` Connect an external service (CNAME) to the cluster.

---

### Create a Service

Explore the ./resources/service.yaml service configuration file:

```
apiVersion: v1
kind: Service
metadata:
  name: k8s-real-demo
  labels:
    app: k8s-real-demo
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: 8080
  selector:
    app: k8s-real-demo
```

---

Create the ./resources/service.yaml service using kubectl:

```
$ kubectl apply -f ./resources/service.yaml
service "k8s-real-demo" created
```

---

### Query the Service

Find the NodePort (via Service) and IP (via Node)

```
$ curl [IP]:[NODE_PORT]
```

---

### Load Balancing

Make several calls to the service and notice the different responses.

---

### Explore the k8s-real-demo Service

```bash
$ kubectl get services k8s-real-demo
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
k8s-real-demo   10.0.0.142   <nodes>       8080:30080/TCP   1m
```

```bash
$ kubectl describe services k8s-real-demo
```

Notice the `Endpoints:` entry

---

## Labels

---

### Using Labels

Use `kubectl get pods` with a label query, e.g. for troubleshooting.

```
$ kubectl get pods -l "app=k8s-real-demo"
```

Use `kubectl label` to add labels to a pod.

```
$ kubectl label pod [POD_NAME] 'secure=disabled'
```

```
$ kubectl get pods -l "app=k8s-real-demo"
$ kubectl get pods -l "secure=disabled"
```

---

### Using Labels

We can also modify existing labels

```
$ kubectl label pod [POD_NAME] "app=new-label" --overwrite
$ kubectl describe pod [POD_NAME]
```

---

### Endpoints

View the endpoints of the `k8s-real-demo` service:

(Note the difference from the last call to `describe`. What has happened?)

```
kubectl describe services k8s-real-demo
```


Revert the label to the orginal setting.

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
I'm running version 1.0 on k8s-real-demo-648d67845-jml8j
```

Next, update the image:

```
$ kubectl set image \
  deployment/k8s-real-demo k8s-real-demo=icrosby/k8s-real-demo:v2
```

---

Check status via 

```
kubectl rollout status deployment k8s-real-demo
```

Now verify the new version

```
$ curl [EXTERNAL_IP]:[NodePort]
```

---

### Exercise - Putting It Together

* Build your own Image.
  * In the cloned repo (`cd k8s-example/`)
  * `docker build -t localhost:5000/[YOUR USER]/k8s-real-demo:v1.0.0 .`
  * `docker push localhost:5000/[YOUR USER]/k8s-real-demo:v1.0.0`
* Create a deployment config for your image 
* Deploy on the cluster

---

### Exercise (cont.)

* Scale the deployment to 3 instances
* Verify the scaling was successful and all instances are getting requests.
* Modify the `CONFIG/DOCKERFILE` to return different Version
* Build the image and tag as `v2`
* Update the deployment to use the new tag
* Verify the new version by making an HTTP request
* View the logs of the application

---

### Deploy the Front End

In the /resources folder you will find configuration files for the Front End (Deployment and Service). 

* ./resources/front-end-deploy.yaml
* ./resources/front-end-svc.yaml

Using these configuration files deploy and expose the application on to the cluster.

```
$ kubectl apply -f ./resources/front-end-deploy.yaml

$ kubectl apply -f ./resources/front-end-svc.yaml

```

---

### Accessing the Front End

Find the port on which the front end is exposed (via the service)
And access this via the browser

```
$ kubectl get svc front-end
```

---

### Bonus (if time permits)

Storage!

While we would like to ideally run stateless applications on Kubernetes, we will eventually run into the challenge of requiring state within our cluster.

We will deploy CockroachDB to maintain the state for our demo application.

CockroachDB (link) is an open source 'cloud native' SQL database.

---

## CockroachDB

A Cloud Native SQL Database.

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

* Pull the Kubernetes configuration file from the web (github)
* Created two services (`kubectl get services`)
* Created a `poddisruptionbudget`
* Created a `StatefulSet`
  * Creates 3 Pods
  * Creates 3 PersistentVolumes
  * Creates 3 PersistentVolumeClaims

`We will cover StatefulSets and PersistentVolumes later on`

---

<Verify storage?>

---

### Summary

What have we learned?
* How to deploy a 'real world' application on Kubernetes
* Deal with Deployment and Services
* Connecting Services with labels and selectors
* Scale up/down
* Update a Deployment (rolling update)

---

[Next up, heading to Production...](../03_productionize.md)


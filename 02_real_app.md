## Beyond The Basics

---

### In this section you will:

* Take a non trivial application to work with
* Create a Deployment configuration
* Deploy the application on your cluster
* Scale the application
* Create a Service configuration
* Expose the application on your cluster

---

## 'Real' demo application

We will work the following demo application: `https://github.com/idcrosby/k8s-example` feel free to clone this repo locally, or fork it to your own account.

This application is composed of multiple pieces. One main backend service, a front end (UI) service, and a data layer. We will deploy these pieces one at a time on our cluster.

**Architecture Diagram**

---

## Recap of resource hierarchy

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

./resources/deployment.yaml

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
$ kubectl create -f resources/deployment.yaml
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
$ kubectl port-forward k8s-real-demo 8080:8080
```

* Terminal 2

```
$ curl 0.0.0.0:8080
Hello World!
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
$ kubectl delete po <POD_NAME>
```

---

### Updating Deployments 

(`RollingUpdate`)

* RollingUpdate is the default strategy.
* Updates Pods one (or a few) at a time.

---

### Common workflow

* Update the application, and create a new version.
* Build the new image and tag it with the new version, i.e. v2.
* Update the Deployment with the new image

---

Let's try this

First let's check the current version running (use the same IP and Node Port from before)

```
$ curl [EXTERNAL_IP]:[NodePort]
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

That was fast :)

Now verify the new version

```
$ curl [EXTERNAL_IP]:[NodePort]
```

---

### Exercise - ???

* Remove existing deals pod
* Tag the image as `v1`
* Create a deployment config for the `v1` Deals image 
* Deploy on the cluster
* Scale the deployment to 3 instances
* Verify the scaling was successful
* Modify the `./resources/deals-app/deals.json` file to return different deals
* Build the image and tag as `v2`
* Update the deployment to use the new tag

---

### Try it yourself
* Modify the `./resources/nginx.conf` to return a custom message.

* Create a custom Nginx docker image (see ./resources/Dockerfile-nginx)
* Build and push the container.
* Create a Pod manifest using the image.
* Query the application using `curl` or a browser.
* Access the pod on port 80 using port-forward.
* View the logs of the nginx container.

---

### Debugging

### View the logs of a Pod

Use `kubectl logs` to view the logs for the `<PODNAME>` Pod:

```
$ kubectl logs <PODNAME>
```

> Use the -f flag and observe what happens.

---

### Run an interactive shell inside a Pod

Execute a shell in a Pod, like in Docker:

```
$ kubectl exec -ti <PODNAME> /bin/sh
```

---

### Exercise - ??? Update application to version 3 (tag v3)

Modify the `./resources/deals-app/deals.json` file to return custom deals

* Build and push the docker image
* Create a pod configuration file for the new image
* Create the pod on your cluster
* Access the application via `curl` or a browser
* Check the logs of the application

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

LABEL AND SELECTOR DIAGRAM(S)

---

### Service types

* `ClusterIP` (Default) Exposes the service on a cluster-internal IP.

* `NodePort` Expose the service on a specific port on each node.

* `LoadBalancer` Use a loadbalancer from a Cloud Provider. Creates `NodePort` and `ClusterIP`.

* `ExternalName` Connect an external service (CNAME) to the cluster.

---

### Create a Service

Explore the XXXX service configuration file:

```
```

---

Create the XXXX service using kubectl:

```
$ kubectl create -f XXXX.yaml
service "XXXX" created
```

---

### Query the Service

Use the IP of any of your nodes.

```
$ curl -i [cluster-node-ip]:30080
```

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

---

### Using labels

Use `kubectl get pods` with a label query, e.g. for troubleshooting.

```
kubectl get pods -l "app=k8s-real-demo"
```

Use `kubectl label` to add labels.

```
kubectl label pods k8s-real-demo 'secure=disabled'
```
... and to modify
```
kubectl label pods k8s-real-demo "app=new-label" --overwrite
kubectl describe pod k8s-real-demo
```

---

### Endpoints

View the endpoints of the `k8s-real-demo` service:

(Note the difference from the last call to `describe`)

```
kubectl describe services k8s-real-demo
```

---

### Exercise - Labels

* Add the `k8s-real-demo` label to our previously deployed pod (`hello-794f7449f5-bm5k4`).
* Call the service several times and see what happens.
* Remove the label.

---

### Exercise - Deploy the Front End

In the /resources folder you will find configuration files for the Front End (Deployment and Service). Using these configuration files deploy the service on to your cluster, and access the application through the browser.

---

### Bonus (if time permits)

Storage!

While we would like to ideally run stateless applications on Kubernetes, we will eventually run into the challenge of requiring state within our cluster.

We will deploy CockroachDB to maintain the state for our demo application.

CockroachDB (link) is an open source 'cloud native' SQL database.

---

### Deploying CockRoachDB

```
$ kubectl apply -f https://raw.githubusercontent.com/cockroachdb/cockroach/master/cloud/kubernetes/cockroachdb-statefulset.yaml
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

---

...

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


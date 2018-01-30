## A more realistic application

---

In this section you will:

* 
* Create a docker container image.
* 
* 
* Interact with Pods remotely using kubectl.

---

### What is a Pod?

* Collection of
  * Application container(s)
  * Storage
  * Network
* Unit of deployment
* Unit of scaling

---


---

## Deployments

---

### Creating and Managing Deployments
In this section we will
* Combine what we learned about Pods and Services
* Create a deployment manifest
* Scale our Deployment / ReplicaSet
* Update our application (Rolling Update)

---

### ReplicaSet
A ReplicaSet ensures that a specified number of Pods are running at any given time.

### Deployment
A Deployment manages ReplicaSets and defines how updates to Pods should be rolled out.

---

### Creating a Deployment

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: hello-node
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: hello-node
    spec:
      containers:
      - name: hello-node
        image: nginx:1.12
        ports:
        - containerPort: 8080
```

---

### Deploy to K8s

```
$ kubectl create -f resources/deployment-v1.yaml
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
hello-node-364036756   1         1         1         16s
```

---

### Scale up/down the Deployment

```
$ kubectl scale deployments hello-node --replicas=2
deployment "hello-node" scaled
```

---

### Check the status of the Deployment

Notice the new pod(s)
```
$ kubectl get pods
```

Look at the `Events` at the bottom

```
$ kubectl describe deployment hello-node
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
  deployment/hello-node hello-node=icrosby/hello-node:v2
```

---

Check status via 

```
kubectl rollout status deployment hello-node
```

That was fast :)

Now verify the new version

```
$ curl [EXTERNAL_IP]:[NodePort]
```

---

### Cleanup

```
$ kubectl delete svc hello-node
$ kubectl delete -f resources/deployment-v1.yaml
```
* If the number of Pods is large, this may take a while to complete.
* To leave the Pods running instead,  
use `--cascade=false`.
* If you try to delete the Pods before deleting the Deployment, the ReplicaSet will just replace them.

---

### Exercise - Deploy Deals microservice

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

### View Pod details

Use the `kubectl get` and `kubectl describe` to view details for the `hello-node` Pod:

```
$ kubectl get pods
```

```
$ kubectl describe pods <pod-name>
```

---

### Interact with a Pod remotely

* Pods get a private IP address by default.
* Cannot be reached from outside the cluster.
* Use `kubectl port-forward` to map a local port to a port inside the `hello-node` pod.


---

### Use two terminals

* Terminal 1

```
$ kubectl port-forward hello-node 8080:8080
```

* Terminal 2

```
$ curl 0.0.0.0:8080
Hello World!
```

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

### Exercise - Deploy the Deals microservice

Modify the `./resources/deals-app/deals.json` file to return custom deals

* Build and push the docker image
* Create a pod configuration file for the new image
* Create the pod on your cluster
* Access the application via `curl` or a browser
* Check the logs of the application

---

## Services

---

### Introduction to services
* Stable endpoints for Pods.
* Based on labels and selectors.

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

### Explore the hello-node Service

```bash
$ kubectl get services hello-node
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
hello-node   10.0.0.142   <nodes>       8080:30080/TCP   1m
```

```bash
$ kubectl describe services hello-node
```

---

### Using labels

Use `kubectl get pods` with a label query, e.g. for troubleshooting.

```
kubectl get pods -l "app=hello-node"
```

Use `kubectl label` to add labels.

```
kubectl label pods hello-node 'secure=disabled'
```
... and to modify
```
kubectl label pods hello-node "app=goodbye-node" --overwrite
kubectl describe pod hello-node
```

---

View the endpoints of the `hello-node` service:

(Note the difference from the last call to `describe`)

```
kubectl describe services hello-node
```

---

### Exercise - Expose the Deals microservice

* Create a service for the deals pod.
* Expose the service via nodePort
* Access the service using `curl` or a browser.

---

### Cleanup

```
kubectl delete po --all
```

We will leave the service for now

---

[Next up Deployments...](../04_deployments.md)



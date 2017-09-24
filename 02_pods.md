## Creating and managing pods

In this lab you will:

* Create a simple Hello World node.js application.
* Create a docker container image.
* Write a Pod configuration file.
* Create and inspect Pods.
* Interact with Pods remotely using kubectl.

----

### What is a Pod?

* Collection of
  * Application container(s)
  * Storage
  * Network
* Unit of deployment
* Unit of scaling

----

### Deploy application to Kubernetes
```
kubectl run hello-node --image=nginx:1.12 --port=80

deployment "hello-node" created
```

----

### Check Deployment and Pod

```
kubectl get deployment
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
hello-node   1         1         1            1           49s

kubectl get pod
NAME                          READY     STATUS    RESTARTS   AGE
hello-node-2399519400-02z6l   1/1       Running   0          54s
```

----

### Check metadata about the cluster, events and kubectl configuration

```
kubectl cluster-info
kubectl get events
kubectl config view
```

----

### Creating a Pod manifest

Explore the `hello-node` pod configuration file:

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
      image: nginx:1.12
      ports:
        - containerPort: 80
```

----

### Create the Pod using kubectl:

```bash
kubectl delete deployment hello-node
kubectl create -f configs/pod.yaml
```

----

### View Pod details

Use the `kubectl get` and `kubectl describe` to view details for the `hello-node` Pod:

```
kubectl get pods
```

```
kubectl describe pods <pod-name>
```

----

### Interact with a Pod remotely

* Pods get a private IP address by default.
* Cannot be reached from outside the cluster.
* Use `kubectl port-forward` to map a local port to a port inside the `hello-node` pod.



----

Use two terminals.

* Terminal 1

```
kubectl port-forward hello-node 8080:8080
```

* Terminal 2

```
curl 0.0.0.0:8080
Hello World!
```

----

### Do it yourself
* Create an `nginx.conf` which returns a  
`200 "Hello Kiwi"`.
* Create a custom Nginx image.
* Build the container on Minikube.
* Create a Pod manifest using the image.
* Query the application using `curl` or a browser.
* Access the pod on port 80 using port-forward.
* View the logs of the nginx container.

----

### Debugging

### View the logs of a Pod

Use `kubectl logs` to view the logs for the `<PODNAME>` Pod:

```
kubectl logs <PODNAME>
```

> Use the -f flag and observe what happens.

----

### Run an interactive shell inside a Pod

Execute a shell in a Pod, like in Docker:

```
kubectl exec -ti <PODNAME> /bin/sh
```

----

[Next up Services...](../04_services.md)


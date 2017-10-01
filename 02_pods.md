## Pods

---

## Creating and managing pods

In this section you will:

* Create a simple Hello World node.js application.
* Create a docker container image.
* Write a Pod configuration file.
* Create and inspect Pods.
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

### Deploy application to Kubernetes
```
$ kubectl run hello-node --image=nginx:1.12 --port=80

deployment "hello-node" created
```

---

### Check Deployment and Pod

```
$ kubectl get deployment
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
hello-node   1         1         1            1           49s

$ kubectl get pod
NAME                          READY     STATUS    RESTARTS   AGE
hello-node-2399519400-02z6l   1/1       Running   0          54s
```

---

#### Check metadata about the cluster, events and kubectl configuration

```
kubectl cluster-info

kubectl get events

kubectl config view
```

---

### Creating a Pod manifest

Explore the `hello-node` pod configuration file:
(./resources/hello-node-pod.yaml)

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

---

### Create the Pod using kubectl:

First clean up original

```bash
$ kubectl delete deployment hello-node
```

```bash
$ kubectl create -f resources/hello-node-pod.yaml
```

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

[Next up Services...](../03_services.md)


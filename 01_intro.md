## Kubernetes Basics

---

### Step 1 kubectl basics

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
```bash
$ kubectl cluster-info
Kubernetes master is running at https://35.189.206.159
```

To further debug and diagnose cluster problems, use `kubectl cluster-info dump`


---

### Step 2 deploy a simple application 

Let’s run our first app on Kubernetes with the kubectl run command. The `run` command creates a new deployment for the specified container. This is the simpliest way of deploying a container.

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

### List your deployments

```bash
$ kubectl get deployments
NAME        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
hello-kubernetes   1         1         1            1           31s
```

We see that there is 1 deployment running a single instance of your app. 

---

### Inspect your application

With 
```
kubectl get <obejct>
```
and 
```
kubectl describe <object>
```
you can gather information about the status of your objects like pods, deployments, services, etc.

---

### Step 3 View our app

By default applications are only visible inside the cluster. We can create a proxy to connect to our application.  
Find out the pod name:
```
$ kubectl get pod
NAME                               READY     STATUS    RESTARTS   AGE
hello-kubernetes-624527933-nth9d   1/1       Running   0          2m
```
Create the proxy:
```bash
$ kubectl port-forward hello-kubernetes-624527933-nth9d 8080 
```
We now have a connection between our host and the Kubernetes cluster.

---

### Accessing the application

To see the output of our application, run a curl request in a new terminal window:
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

### Expose service while creating the deployment

`kubectl port-forward` is meant for testing services that are not exposed. To expose the application, use a service (covered later).

Delete old deployment
```
$ kubectl delete deployment hello-kubernetes
```

---

Create a new **Deployment** and a **Service**

```
$ kubectl run hello --image=gcr.io/google_containers/echoserver:1.4 \
   --port=8080 \
   --expose \
   --service-overrides='{ "spec": { "type": "NodePort" } }'
service "hello" created
deployment "hello" created
```

This creates a new deployment and a service of type:NodePort. A random high port will be allocated to which we can connect.

---

View the **Service**:

```
$ kubectl get service
NAME               CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
hello   			10.0.122.112   <nodes>       8080:30659/TCP   10m
```

We can see the port on which it is exposed, but what is the external IP?

---

To find the IP on which to call we need information on the nodes:

```
$ kubectl get nodes -o wide
NAME                           STATUS                     AGE       VERSION   EXTERNAL-IP      OS-IMAGE                             KERNEL-VERSION
kubernetes-master              Ready,SchedulingDisabled   17m       v1.7.5    35.187.38.163    Container-Optimized OS from Google   4.4.52+
kubernetes-minion-group-c9bz   Ready                      17m       v1.7.5    35.189.206.159   Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
kubernetes-minion-group-cfzx   Ready                      17m       v1.7.5    35.195.36.237    Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
kubernetes-minion-group-ftw1   Ready                      17m       v1.7.5    35.195.61.242    Debian GNU/Linux 7 (wheezy)          3.16.0-4-amd64
```

---

Access the external IP with curl:

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

### Cleanup

```
$ kubectl delete deployment,service hello
deployment "hello" deleted
service "hello" deleted
```

---

[Next up Pods...](../02_pods.md)


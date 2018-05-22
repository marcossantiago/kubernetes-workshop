
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

```bash
$ kubectl cluster-info

Kubernetes master is running at https://api-example-cluster-k8s-l-1dt7vk-944985396.us-east-2.elb.amazonaws.com
KubeDNS is running at https://api-example-cluster-k8s-l-1dt7vk-944985396.us-east-2.elb.amazonaws.com/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
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

Create a port-forward for the pod

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
$ jobs 
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
hello   			10.0.122.112   <nodes>       8080:30659/TCP   10m
```

We can see the port on which it is exposed, but what is the external IP?

---

To find the IP on which to call we need information on the nodes (use the EXTERNAL-IPs from any node):

```
$ $ kubectl get nodes -o wide
NAME                                           STATUS    ROLES     AGE       VERSION   EXTERNAL-IP      OS-IMAGE                      KERNEL-VERSION
ip-172-20-119-162.us-east-2.compute.internal   Ready     master    23m       v1.9.3    52.15.56.60      Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-125-29.us-east-2.compute.internal    Ready     node      22m       v1.9.3    18.220.252.62    Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-50-208.us-east-2.compute.internal    Ready     node      22m       v1.9.3    18.220.253.174   Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-55-67.us-east-2.compute.internal     Ready     node      22m       v1.9.3    52.15.126.218    Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-62-48.us-east-2.compute.internal     Ready     master    23m       v1.9.3    13.58.187.3      Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-73-52.us-east-2.compute.internal     Ready     node      22m       v1.9.3    52.15.67.117     Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-75-166.us-east-2.compute.internal    Ready     node      22m       v1.9.3    18.219.33.72     Debian GNU/Linux 8 (jessie)   4.4.115-k8s
ip-172-20-93-40.us-east-2.compute.internal     Ready     master    23m       v1.9.3    18.218.224.237   Debian GNU/Linux 8 (jessie)   4.4.115-k8s
```

---

Access the external IP with Curl:

```
$ curl 3 52.15.126.218:30659

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

* Basics of Kubernetes.
* How to deploy a simple application on to our own cluster.

---

[Next up, a real application!](/#/a-real-application)

or

[What next?](/#/what-next-)
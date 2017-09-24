### Step 1 kubectl basics

* The format of a kubectl command is: 
```
kubectl [action] [resource]
```
* This performs the specified action  (like `create`, `describe`) on the specified resource (like `node`, `container`). 
* Use `--help` after the command to get additional info about possible parameters
```
kubectl get nodes --help
```

----

Check that kubectl is configured to talk to your cluster, by running the kubectl version command:
```bash
kubectl version
```

You can see both the client and the server versions.

----

To view the nodes in the cluster, run the `kubectl get nodes` command:
```bash	
kubectl get nodes
NAME       STATUS    AGE       VERSION
minikube   Ready     7m        v1.6.0
```

Here we see the available nodes, just one in our case. Kubernetes will choose where to deploy our application based on the available Node resources.

----

### Step 2 deploy a simple application 

Let’s run our first app on Kubernetes with the kubectl run command. The `run` command creates a new deployment for the specified container. This is the simpliest way of deploying a container.

```bash
kubectl run hello-kubernetes \
--image=gcr.io/google_containers/echoserver:1.4 --port=8080

deployment "hello-kubernetes" created
```

----

This performed a few things:
* Searched for a suitable node.
* Scheduled the application to run on that node.
* Configured the cluster to reschedule the instance on a new node when needed.

----

### List your deployments

```bash
kubectl get deployments
NAME        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
hello-kubernetes   1         1         1            1           31s
```

We see that there is 1 deployment running a single instance of your app. 

----

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

----

### Step 3 View our app

By default applications are only visible inside the cluster. We can create a proxy to connect to our application.  
Find out the pod name:
```
kubectl get pod
```
Create the proxy:
```bash
kubectl port-forward hello-kubernetes-3015430129-g95j6 8080 
```
We now have a connection between our host and the Kubernetes cluster.

----

### Accessing the application

To see the output of our application, run a curl request in a new terminal window:
```bash
curl http://localhost:8080
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

----

### Expose service while creating the deployment

`kubectl port-forward` is meant for testing services that are not exposed. To expose the application, use a service.

Delete old deployment
```
kubectl delete deployment hello-kubernetes
```

----

Create a new deployment and a service
```
kubectl run hello-kubernetes \
--image=gcr.io/google_containers/echoserver:1.4 \
--port=8080 --expose \
--service-overrides='{ "spec": { "type": "NodePort" } }'

service "hello-kubernetes" created
deployment "hello-kubernetes" created
```
This creates a new deployment and a service of type:NodePort. A random high port will be allocated to which we can connect.

----

View the service:
```
kubectl get service
kubectl get svc
NAME             CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
hello-kubernetes   10.0.0.233   <nodes>       8080:31075/TCP   24s
kubernetes       10.0.0.1     <none>        443/TCP          28m
```

Access the application with curl

(use the IP of one of your nodes)

```
curl 0.0.0.0:31075
```

----

### Cleanup

```
kubectl delete deployment,service hello-kubernetes
deployment "hello-kubernetes" deleted
service "hello-kubernetes" deleted
```

----

[Next up Pods...](../03_pods.md)


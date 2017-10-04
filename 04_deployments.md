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

[Next up Probes](../05_probes.md)

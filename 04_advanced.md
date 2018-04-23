# Advanced Features

---

We now have a 'real life' application deployed and running. We have added the necessary pieces to consider this 'production ready'.

---

## In this section we will cover

- Volumes / Storage
- Ingress
- Autoscaling
- Advanced Deployments

---

### Storage

 - A Pod is made up of one or more containers and data volumes that can be mounted inside the containers. 
 - In this section you will learn how to: 

* define a deployment backed by a emptyDir
* define a deployment backed by a emptyDir(memory backed storage)
* define a deployment backed by a persistent volume and persistent volume claim 
* define a deployment backed by a persistent volume and persistent volume claim using a StorageClass

---

### Volumes

Volumes are means to save data, as well as share it between containers. Any volumes in a pod are accessible by all containers running inside that same pod. The data is persisted across container restarts. 

There are a large number of implementations to back the storage. From local options such as `hostPath`, cloud specific options such as `awsElasticBlockStore`, distributed storage such as `ceph` or `nfs`, and many many more.

---

### emptyDir

* In this exercise we will demonstrate the use of an emptyDir as a volume.

---

 * The volume is of type `emptyDir`. 
 * The kubelet will create an empty directory on the node when the Pod is scheduled. 
 * Once the Pod is destroyed, the kubelet will delete the directory.

---

Find the `empty-dir.yaml` file in the `./kuberetes/` folder

```
apiVersion: v1
kind: Pod
metadata:
  name: busybox
spec:
  containers:
  - name: busy
    image: busybox
    volumeMounts:
    - name: test
      mountPath: /busy
    command:  
      - sleep
      - "3600"
  - name: box
    image: busybox
    volumeMounts:
    - name: test
      mountPath: /box
    command:
      - sleep
      - "3600"
  volumes:
  - name: test
    emptyDir: {} 
``` 

---

Apply the configuration file:

```
kubectl apply -f ./kubernetes/empty-dir.yaml
```

This create a single pod with two containers and one volume.

---

Once the pod is deployed, `exec` into one of the containers:

```
$ kubectl exec -ti busybox -c box /bin/sh
#
```

Now create a new file in the `/box` directory:

```
$ touch /box/foobar
$ ls /box
```

---

Exit from the ccontainer and exec into the second:

```
#exit
$ kubectl exec -ti busybox -c busy
#
```

List the contents of the `/busy` directory

```
#ls -l /busy
total 0
-rw-r--r--    1 root     root             0 Nov 19 16:26 foobar
```

Note that this works in both directions. But once the pod dies, the directory is removed.

---

### Persistent Volumes

While handling volumes with a persistent volume definition and abstracting the storage provider using a claim is powerful, an administrator of the cluster still needs to create those volumes in the first place.

---

## Dynamic Provisioning
Alternatively a `StorageClass` can be configured to dynamically provision Persistent Volumes (PV).  User requests a Persistent Volume Claim (PVC), and the provisioner defined in the StorageClass will dynamically create the volume.

---

Here is an example of a StorageClass on AWS:

```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: slow
provisioner: kubernetes.io/aws-ebs
parameters:
  type: io1
  zones: us-east-1d, us-east-1c
  iopsPerGB: "10"
```
---

GKE comes with a default StorageClass that will dynamically provision persitent disks on demand. We can see this by running:

```
$ kubectl get storageclass
...
$ kubectl describe storageclass standard
...
```

---

First we create a persistent volume claim including the above storage class. (This file can be found at `./kubernetes/pvc.yaml`). Be sure to replace the <MY_SC_CLAIM>.

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: <MY_SC_CLAIM>
  annotations:
    volume.beta.kubernetes.io/storage-class: "standard"
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

---

Let's create the claim and then verify that a persistent volume is created automatically. It should be bound to the claim requesting storage.

```
$ kubectl create -f pvc.yaml
$ kubectl get pv
$ kubectl get pvc
```

---

Finally, if we delete the persistent volume claim, we can see the volume gets released and is automatically deleted

```
$ kubectl delete pvc <MY_SC_NAME>
$ kubectl get pv
```

---

### Do it yourself

* Our sample app currently is writing logs to `/var/log/app.log`
* Create a Persistent Volume Claim to persist logs.
* Update the exiting deployment to use the PVC.
* Verify the logs are persisted (even after the application dies)

---

### Cheat

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: logclaim
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 8Gi
```

---

```
...
    volumeMounts:
    - name: logs
      mountPath: /var/log
  volumes:
    - name: logs
      persistentVolumeClaim:
        claimName: logclaim
```

---

### Ingress

Typically, services and pods have IPs only routable by the cluster network. All traffic that ends up at an edge router is either dropped or forwarded elsewhere. Conceptually, this might look like:
```
    internet
        |
  ------------
  [ Services ]
```

---

An Ingress is a collection of rules that allow inbound connections to reach the cluster services.
```
    internet
        |
   [ Ingress ]
   --|-----|--
   [ Services ]
```

---

It can be configured to:
* Give services externally-reachable urls
* Loadbalance traffic
* Terminate SSL
* Offer name based virtual hosting

An Ingress controller is responsible for fulfilling the Ingress, usually with a loadbalancer, though it may also configure your edge router or additional frontends to help handle the traffic in an HA manner.

---

### Ingress controller

In order for the Ingress resource to work, the cluster must have an `Ingress Controller` running.

An `Ingress Controller` is a daemon, deployed as a Kubernetes Pod, that watches the ApiServer's /ingresses endpoint for updates to the Ingress resource. Its job is to satisfy requests for ingress.

---

### Ingress Workflow

* Poll until apiserver reports a new Ingress.
* Write the LB config file based on a go text/template.
* Reload LB config.

---

### Example
Ingress resource
```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: frontend-ingress
spec:
  rules:
    - host: frontend.example.com
      http:
        paths:
          -
            backend:
              serviceName: front-end
              servicePort: 80
            path: /
```
*POSTing this to the API server will have no effect if you have not configured an Ingress controller.*

---

### Ingress Controllers

There are multiple Ingress controllers we can use, including nginx, HAProxy or Traefik. However there are some subtle differences between them.

Our cluster has Google's L7 ingress controller already installed.

---

### Deploy some application

First we need to deploy some application to publish. To keep this simple we will use the echoheaders app that just returns information about the http request as output.
```
kubectl run echoheaders --image=gcr.io/google_containers/echoserver:1.4 \
  --replicas=1 --port=8080
```
Now we expose the same application in two different services (so we can create different Ingress rules).
```
kubectl expose deployment echoheaders --type=NodePort --port=80 \
  --target-port=8080 --name=echoheaders-x
kubectl expose deployment echoheaders --type=NodePort --port=80 \
--target-port=8080 --name=echoheaders-y
```

---

### Add ingress rules (in `./kubernetes/echo-ingress.yaml`)

```
  rules:
    - host: foo.bar.com
      http:
        paths:
          - path: /foo
            backend:
              serviceName: echoheaders-x
              servicePort: 80           
    - host: bar.baz.com
      http:
        paths:
          - path: /bar
            backend:
              serviceName: echoheaders-y
              servicePort: 80
          - path: /foo
            backend:
              serviceName: echoheaders-x
              servicePort: 80
```

---


```
kubectl create -f ./kubernetes/echo-ingress.yaml
```

Note this may take a few minutes spin up.

Wait until you see an IP under the `ADDRESS` via `kubectl get ingress`

---

### Accessing the application

To access the applications via a browser you need either to edit your `/etc/hosts` file with the domains `foo.bar.com` and `bar.baz.com` pointing to the IP of your k8s cluster. Or use a browser plugin to manipulate the host header.

Here we'll use `curl`.

```
curl -H "Host: foo.bar.com" http://<INGRESS_ADDRESS>/bar
curl -H "Host: bar.baz.com" http://<INGRESS_ADDRESS>/bar
curl -H "Host: bar.baz.com" http://<INGRESS_ADDRESS>/foo
```

---

### Do it yourself

* Expose our application via Ingress.
* Write an ingress manifest to expose the front-end service on port 80 listening on goto-chicago.example.com/demo
* Access the application via `curl` or a browser on port 80.

---

## Auto Scaling

As we have seen in the previous section scaling our applications is very simple. `kubectl scale deployment/k8s-real-demo replicas=5`.  However, ideally this would not be a manual action.

Kubernetes supports this via the `HorizontalPodAutoscaler` resource.

---

## HorizontalPodAutoscaler (HPA)

* Periodically fetches metrics (default 30 seconds)
* Compares to user specified target value.
* Adjusts the number of replicas (pods) of a Deployment if needed.
* CPU usage is built in.
* Fetched from Heapster.
* Can also read Prometheus (for custom metrics)

---

<img src="img/5. Auto sclaing process-01.png">

---

### Adding an HPA

We can add an HPA to our existing Deployment

./kubernetes/hpa.yaml

```
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: k8s-real-demo-hpa
spec:
  scaleTargetRef:
    kind: Deployment
    name: k8s-real-demo
  minReplicas: 1
  maxReplicas: 5
  targetCPUUtilizationPercentage: 50
```

---

We can create this resource on the cluster just as we have done with the others:

```
$ kubectl apply -f ./kubernetes/hpa.yaml
```

And view the resource with `kubectl get` and `kubectl describe`

---

### Testing the HPA

In order to test that the HPA works as expected we will need our application to consume > 80% CPU. Luckily our application has a very CPU intensive endpoint (`/mineBitcoin`).

---

Make a request (or several) to the endpoint `/mineBitcoin` (note this endpoint can take a query parameter `seconds`).

See if your Deployment scales successfully.

---

### Advanced Deployments

Kubernetes offers a variety of ways to release an application.

The correct option depends on your specific requirements and use case.

In this section we will try out several options and look at some possible use cases.

---

## Deployment Strategies

- **Recreate**: Terminate the old version then release a new one
- **Ramped**: Release a new version via a rolling update
- **Blue/Green**: Release a new version alongside the old version then switch traffic
- **Canary**: Release a new version to a subset of users, then proceed to a full rollout
- **A/B testing**: release a new version to a subset of users in a precise way (HTTP headers, cookie, weight, etc.).

---

## Recreate deployment

First terminate (all instances of) the old version and then release the new one.

Add the below `strategy` section to your existing deployment.yaml and remove the `env` section

```yaml
...
spec:
  replicas: 3
  strategy:
    type: Recreate
...

```

```
$ kubectl apply -f ./resources/deployment.yaml
```

---

Verify the deployment:

```
$ curl $EXTERNAL_IP:[NodePort]
Hello from Container Solutions.
I'm running version 2.1 on k8s-real-demo-5449767b94-hnk78
```

In the output we can see both the Pod ID as well as the version of the application.

---

To see the deployment in action, open a new terminal and run the following command in another terminal:

```
$ watch -n1 kubectl get po
```
or
```
$ kubectl get po -w
```

Then deploy the version 2 of the application:

```
$ kubectl set image deploy/k8s-real-demo k8s-real-demo=icrosby/k8s-real-demo:v2
```

---

Now test the second deployment progress.

(N.B. Since we are not removing the service, the NodePort will not change)

```
$ export SERVICE_URL=$EXTERNAL_IP:[NodePort]
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

---

## Blue/Green Deployment
Release a new version alongside the old version then switch traffic


Deploy the first application

```
$ kubectl apply -f ./resources/deployment.yaml
```

---

Test if the deployment was successful

```
$ curl $EXTERNAL_IP:[NodePort]
Hello from Container Solutions.
I'm running version 1.0 on k8s-real-demo-5449767b94-hnk78
```

---

To see the deployment in action, open a new terminal and run the following command:

```
$ kubectl get po -w
```

Then create a second deploy with version 2 of the application:

```
$ kubectl apply -f ./resources/deployment-v2.yaml
```

---

Side by side, 3 pods are running with version 2 but the service still send traffic to the first deployment.

Try manually test one of the new pods by port-forwarding it to your local environment.

---

Once your are ready, you can switch the traffic to the new version by patching the service to send traffic
to all pods with label app: k8s-real-demo-v2:

```
$ kubectl patch service my-app -p \
'{"spec":{"selector":{"app: k8s-real-demo-v2}}}'
```

Alternatively you can use

```
$ kubectl edit service k8s-real-demo
```

---

Test if the second deployment was successful:

```
$ export SERVICE_URL=$EXTERNAL_IP:[NodePort]
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

In case you need to rollback to the previous version:

```bash
$ kubectl patch service my-app -p \
'{"spec":{"selector":{"version":"v1.0.0"}}}'
```

---

If everything is working as expected, you can then delete the v1.0.0 deployment:

```
$ kubectl delete deploy my-app-v1
```

---

### Clean Up

```
$ kubectl deployment k8s-real-demo
```

---

[Next up Setting up an HA cluster](../05_cluster.md)

## Creating and Managing Services

In this section you will create a `hello-node` service and "expose" the `hello-node` Pod. 

You will learn how to:

* Create a service.
* Use labels and selectors to expose Pods.

----

### Introduction to services
* Stable endpoints for Pods.
* Based on labels and selectors.

----

### Service types

* `ClusterIP` (Default) Exposes the service on a cluster-internal IP.

* `NodePort` Expose the service on a specific port on each node.

* `LoadBalancer` Use a loadbalancer from a Cloud Provider. Creates `NodePort` and `ClusterIP`.

* `ExternalName` Connect an external service (CNAME) to the cluster.

----

### Create a Service

Explore the hello-node service configuration file:
(./resources/hello-node-service.yaml)

```
apiVersion: v1
kind: Service
metadata:
  name: hello-node
spec:
  type: NodePort
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 80
    nodePort: 30080
  selector:
    app: hello-node
```

*Setting nodePort is optional.

----

Create the hello-node service using kubectl:

```
$ kubectl create -f resources/hello-node-service.yaml
service "hello-node" created
```

----

### Query the Service

Use the IP of any of your nodes.

```
$ curl -i [cluster-node-ip]:30080
```

----

### Explore the hello-node Service

```bash
$ kubectl get services hello-node
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
hello-node   10.0.0.142   <nodes>       8080:30080/TCP   1m
```

```bash
$ kubectl describe services hello-node
```

----

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

----

View the endpoints of the `hello-node` service:

(Note the difference from the last call to `describe`)

```
kubectl describe services hello-node
```

----

### Do it yourself

* Create a service for the nginx pods (Created earlier).
* Expose port 80 to a static nodePort 31000.
* Access the service using `curl` or a browser.

----

### Cleanup

```
kubectl delete po --all
kubectl delete svc --all
```

----

[Next up Deployments...](../04_deployments.md)

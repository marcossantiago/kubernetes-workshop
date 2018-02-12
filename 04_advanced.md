## Advanced

- Autoscaling
- Advanced Deployments
- Storage

---

### Auto Scaling

---

Check for existing material?

- Create application which can consume CPU (have example created for Max/Michael)
- deploy HPA
- Cause CPU usage -> scale up
- Remove, scale down
- Combine with monitoring??

---

---

### Advanced Deployments

Kubernetes offers a variety of ways to release an application. 

The correct option depends on your specific requirements and use case.

In this section we will try out each option and look at some possible use cases.

---

## Deployment Strategies

- **recreate**: Terminate the old version then release a new one
- **ramped**: Release a new version via a rolling update
- **blue/green**: Release a new version alongside the old version then switch traffic
- **canary**: Release a new version to a subset of users, then proceed to a full rollout
- **a/b testing**: release a new version to a subset of users in a precise way (HTTP headers, cookie, weight, etc.).

---

You can find the code for the sample deploy application under ./resources/deploy-app

Thanks to my colleague Etienne Tremel for the application.

---

## Recreate deployment
First terminate (all instances of) the old version and then release the new one.

Let's take a look at the config file and then deploy the first application

```yaml
...
spec:
  replicas: 3
  strategy:
    type: Recreate
...

```

```
$ kubectl apply -f ./resources/recreate-app-v1.yaml
```

---

Retrieve the IP of one of the nodes and store this as an envrionment variable:
```
$ export EXTERNAL_IP=$(kubectl get nodes \
-o jsonpath='{.items[1].status.addresses[?(@.type=="ExternalIP")].address}')
```

Next we need the NodePort on which our service is exposed:
```
$ kubectl describe svc my-app                                                                                            ✓  10354  15:01:43
Name:			my-app
Namespace:		user-3
Labels:			app=my-app
Annotations:		kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"my-app"},"name":"my-app","namespace":"user-3"},"spec":{"ports":[{"nam...
Selector:		app=my-app
Type:			NodePort
IP:			10.0.173.234
Port:			http	8080/TCP
NodePort:		http	31461/TCP
Endpoints:
Session Affinity:	None
Events:			<none>
```

---

Test if the deployment was successful:

```
$ curl $EXTERNAL_IP:[NodePort]
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 \
   - Host: my-app-177300127-sbd1d, Version: v1.0.0
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
$ kubectl apply -f ./resources/recreate-app-v2.yaml
```

---

Now test the second deployment progress. 

(N.B. Since we are not removing the service, the NodePort will not change)

```
$ export SERVICE_URL=$EXTERNAL_IP:[NodePort]
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

---

Cleanup:

```
$ kubectl delete all -l app=my-app
```

---

## Ramped deployment
Release a new version on a rolling update fashion, one after the other.

Deploy the first application

```
$ kubectl apply -f ./resources/ramped-app-v1.yaml
```

---

Test if the deployment was successful

```
$ curl $EXTERNAL_IP:[NodePort]
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 \
   - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

---

To see the deployment in action, open a new terminal and run the following command:

```
$ kubectl get po -w
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/ramped-app-v2.yaml
```

---

Next, test the second deployment progress:

```
$ export SERVICE_URL=$EXTERNAL_IP:[NodePort]
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

We can also verify the version via the pods
```
$ kubectl get pods
NAME                     READY     STATUS    RESTARTS   AGE
my-app-873192836-0pvw3   1/1       Running   0          1m
my-app-873192836-ddfp4   1/1       Running   0          55s
my-app-873192836-wjwcc   1/1       Running   0          1m

$ kubectl describe pod my-app-873192836-0pvw3 
...
Containers:
  my-app:
    Container ID:	docker://5d87a32691adaa933f1a6a956fb76a77ade3f9da669810c66f954b94786724a3
...
```

---

In case you discover an issue with the new version, you can undo the rollout:

```
$ kubectl rollout undo deploy my-app
```

---

If you can also pause the rollout if you want to run the application for a subset of users:

```
$ kubectl rollout pause deploy my-app
```

Then if you are satisfy with the result, rollout:

```
$ kubectl rollout resume deploy my-app
```

---

Cleanup:

```
$ kubectl delete all -l app=my-app
```

---

## Blue/Green Deployment
Release a new version alongside the old version then switch traffic


Deploy the first application

```
$ kubectl apply -f ./resources/blue-green-app-v1.yaml
```

---

Test if the deployment was successful

```
$ curl $EXTERNAL_IP:[NodePort]
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 \
   - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

---

To see the deployment in action, open a new terminal and run the following command:

```
$ kubectl get po -w
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/blue-green-app-v2.yaml
```

---

Side by side, 3 pods are running with version 2 but the service still send traffic to the first deployment.

Try manually test one of the new pods by port-forwarding it to your local environment.

---

Once your are ready, you can switch the traffic to the new version by patching the service to send traffic
to all pods with label version=v2.0.0:

```
$ kubectl patch service my-app -p \
'{"spec":{"selector":{"version":"v2.0.0"}}}'
```

Alternatively you can use

```
$ kubectl edit service my-app
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

Cleanup:

```
$ kubectl delete all -l app=my-app
```

---

## Canary Deployment
Release a new version to a subset of users, then proceed to a full rollout

Deploy the first application:

```
$ kubectl apply -f ./resources/canary-app-v1.yaml
```

---

Test if the deployment was successful:

```
$ curl $EXTERNAL_IP:[NodePort]
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - \
   Host: my-app-177300127-sbd1d, \Version: v1.0.0
```

---

Again, let's watch the deploment from a new terminal:

```
$ kubectl get po -w
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/canary-app-v2.yaml
```

---

Only one pod with the new version should be running.

You can test if the second deployment was successful:

```
$ export SERVICE_URL=$EXTERNAL_IP:[NodePort]
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

---

If you are happy with it, scale up the version 2 to 3 replicas:

```
$ kubectl scale --replicas=3 deploy my-app-v2
```

Then, when all pods are running, let's verify that we are only hitting the new version:

```
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

You can now safely delete the old deployment:

```
$ kubectl delete deploy my-app-v1
```

---

Cleanup:

```
$ kubectl delete all -l app=my-app
```

---

### Exercise - Advanced deploys of the Deals microservice

* Perform a canary deployment for the Deals service
* Observe both versions of the application running
* Revert to a single version running

---

[Next up Setting up a cluster](../08_cluster.md)

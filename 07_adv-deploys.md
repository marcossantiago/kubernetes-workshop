## Advanced Deployments
## With Kubernetes

---

## Advanced Deployments

Kubernetes offers a variety of ways to release an application. 

The correct option depends on your specific requirements and use case.

In this section we will try out each option and look at some possible use cases.

---

These exercises created by my colleague Etienne Tremel. He also wrote a great blog post explaining advanced deployments with Kubernetes which you can [read here](https://container-solutions.com/kubernetes-deployment-strategies/)

---

## Deployment Strategies

- **recreate**: Terminate the old version then release a new one
- **ramped**: Release a new version via a rolling update
- **blue/green**: Release a new version alongside the old version then switch traffic
- **canary**: Release a new version to a subset of users, then proceed to a full rollout
- **a/b testing**: release a new version to a subset of users in a precise way (HTTP headers, cookie, weight, etc.). We will cover this in a later lab.


---

## Recreate deployment
First terminate the old version and then release the new one.

Deploy the first application

```
$ kubectl apply -f ./resources/recreate-app-v1.yaml
```

---

The app will be exposed via NodePort service. Retrieve the IP of one of the nodes:
```
$ export EXTERNAL_IP=$(kubectl get nodes \
-o jsonpath='{.items[1].status.addresses[?(@.type=="ExternalIP")].address}')
```

Test if the deployment was successful:

```
$ curl $EXTERNAL_IP:[NodePort]
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 \
   - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

---

To see the deployment in action, open a new terminal and run the following command:

```
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/recreate-app-v2.yaml
```

---

Now test the second deployment progress:

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
$ watch -n1 kubectl get po
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
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/blue-green-app-v2.yaml
```

---

Side by side, 3 pods are running with version 2 but the service still send traffic to the first deployment.

If necessary, you can manually test one of the pods by port-forwarding it to your local environment.

---

Once your are ready, you can switch the traffic to the new version by patching the service to send traffic
to all pods with label version=v2.0.0:

```
$ kubectl patch service my-app -p \
'{"spec":{"selector":{"version":"v2.0.0"}}}'
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

To see the deployment in action, open a new terminal and run a watch command to have a nice view on the progress:

```
$ watch -n1 kubectl get po
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
kubectl scale --replicas=3 deploy my-app-v2
```

Then, when all pods are running, you can safely delete the old deployment:

```
kubectl delete deploy my-app-v1
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

[Next up Autoscaling](../08_autoscaling.md)

# Advanced Deployments

Kubernetes offers a variety of ways to release an application. 

The correct option depends on your specific requirements and use case.

In this section we will try out each option and look at some possible use cases.

----

These exercises created by my colleague Etienne Tremel. He also wrote a great blog post explaining advanced deployments with Kubernetes which you can [read here](https://container-solutions.com/kubernetes-deployment-strategies/)

----

##Recreate deployment - First terminate the old version and then release the new one.

Deploy the first application

```
$ kubectl apply -f ./resources/app-v1.yaml
```

----

Test if the deployment was successful:

```
$ curl $(minikube service my-app --url)
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

----

To see the deployment in action, open a new terminal and run the following command:

```
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/app-v2.yaml
```

----

Now test the second deployment progress:

```
$ export SERVICE_URL=$(minikube service my-app --url)
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

----

Cleanup:

```
$ kubectl delete all -l app=my-app
```

----

##Ramped deployment - Release a new version on a rolling update fashion, one after the other.

Deploy the first application

```
$ kubectl apply -f ./resources/app-v1.yaml
```

----

Test if the deployment was successful

```
$ curl $(minikube service my-app --url)
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

----

To see the deployment in action, open a new terminal and run the following command:

```
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f ./resources/app-v2.yaml
```

----

Next, test the second deployment progress:

```
$ export SERVICE_URL=$(minikube service my-app --url)
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

----

In case you discover an issue with the new version, you can undo the rollout:

```
$ kubectl rollout undo deploy my-app
```

----

If you can also pause the rollout if you want to run the application for a subset of users:

```
$ kubectl rollout pause deploy my-app
```

Then if you are satisfy with the result, rollout:

```
$ kubectl rollout resume deploy my-app
```

----

Cleanup:

```
$ kubectl delete all -l app=my-app
```

----

##Blue/Green Deployment = Release a new version alongside the old version then switch traffic


Deploy the first application

```
$ kubectl apply -f app-v1.yaml
```

----

Test if the deployment was successful

```
$ curl $(minikube service my-app --url)
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

----

To see the deployment in action, open a new terminal and run the following command:

```
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f app-v2.yaml
```

----

Side by side, 3 pods are running with version 2 but the service still send traffic to the first deployment.

If necessary, you can manually test one of the pods by port-forwarding it to your local environment.

Once your are ready, you can switch the traffic to the new version by patching the service to send traffic
to all pods with label version=v2.0.0:

```
$ kubectl patch service my-app -p '{"spec":{"selector":{"version":"v2.0.0"}}}'
```

----

Test if the second deployment was successful:

```
$ export SERVICE_URL=$(minikube service my-app --url)
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

In case you need to rollback to the previous version:

```
$ kubectl patch service my-app -p '{"spec":{"selector":{"version":"v1.0.0"}}}'
```

----

If everything is working as expected, you can then delete the v1.0.0 deployment:

```
$ kubectl delete deploy my-app-v1
```

----

Cleanup:

```
$ kubectl delete all -l app=my-app
```

----

##Canary Deployment - Release a new version to a subset of users, then proceed to a full rollout

Deploy the first application:

```
$ kubectl apply -f app-v1.yaml
```

----

Test if the deployment was successful:

```
$ curl $(minikube service my-app --url)
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

----

To see the deployment in action, open a new terminal and run a watch command to have a nice view on the progress:

```
$ watch -n1 kubectl get po
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f app-v2.yaml
```

----

Only one pod with the new version should be running.

You can test if the second deployment was successful:

```
$ export SERVICE_URL=$(minikube service my-app --url)
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

----

If you are happy with it, scale up the version 2 to 3 replicas:

```
kubectl scale --replicas=3 deploy my-app-v2
```

Then, when all pods are running, you can safely delete the old deployment:

```
kubectl delete deploy my-app-v1
```

----

Cleanup:

```
$ kubectl delete all -l app=my-app
```

----

## A/B Testing

Release a new version to a subset of users in a precise way (HTTP headers, cookie, weight, etc.). 

This doesn’t come out of the box with Kubernetes, it imply extra work to setup a more advanced infrastructure (Istio, Linkerd, Traeffik, custom nginx/haproxy, etc).

----

##A/B testing using Istio


Deploy Istio to the cluster using Helm:

```
$ helm init
$ helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
$ helm install --name service-mesh incubator/istio
```

----

Deploy the service and ingress:

```
$ kubectl apply -f ./service.yaml
$ kubectl apply -f ./ingress.yaml
```

----

Deploy the first application and use istioctl to inject a sidecar container to proxy all in and out
requests:

```
$ kubectl apply -f <(istioctl kube-inject -f app-v1.yaml)
```

----

Test if the deployment was successful:

```
$ curl $(minikube service istio-ingress --url | head -n1)
> 2017-09-20 12:42:33.416123892 +0000 UTC m=+55.563375310 - Host: my-app-177300127-sbd1d, Version: v1.0.0
```

Then deploy the version 2 of the application:

```
$ kubectl apply -f <(istioctl kube-inject -f app-v2.yaml)
```

----

Apply the load balancing rule:

```
$ istioctl create -f ./rules.yaml
```

You can now test if the traffic is correctly splitted amongst both versions:

```
$ export SERVICE_URL=$(minikube service istio-ingress --url | head -n1)
$ while sleep 0.1; do curl $SERVICE_URL; done;
```

----

You should see 1 request on 10 ending up in the version 2.

In the rules.yaml file, you can edit the weight of each route and apply the changes as follow:

```
$ istioctl replace -f ./rules.yaml
```

----

Cleanup:

```
$ kubectl delete all -l app=my-app
$ helm delete service-mesh
$ helm reset
```
----

[Next up Autoscaling](../08_autoscaling.md)

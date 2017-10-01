### Resource Management

---

In this section we'll discuss:
* Quality of Service (QoS)
* Resource Quotas

---

### The resource model

What exactly is a resource in Kubernetes?

**CPU & Memory**
* Accounted
* Scheduled
* Isolated

**Local Storage (Disk or SSD)**
* Accounted (restriction to single partition /)

**Nvidia GPU**
* Alpha support (1 GPU per-node)

---

### From the Top

There are two main kinds of resource management for pods.

* Requests
* Limits

---

### Requests

A Request must be fullfilled for a Pod to be scheduled.  
Assume I have a pod I want to schedule.

* Request 512Mb RAM & 300m CPU
* Kubernetes looks for nodes with space.
* If a node exists, then the Pod is scheduled.

---

### Requests

* A Pod will not be deployed if there is no suitable node.
* Requested resources can be exceeded (with _potential_ consequences).
* Is more of a guideline to resource usage within a Pod.

---

### Limits

A Limit is ignored in the scheduling phase of the Pod lifecycle.

* It is hard, if a Pod tries to exceed it, there are immediate consequences.
* Must be equal to, or greater than the Request (if defined).

---


### Requests and Limits

Repercussions:
* Usage > Request: Pod may 'scavenge' excess resources.
* Usage > Limit: Pod is killed (RAM) or throttled (CPU).

---

### Types of Resource

There are two kinds of Resources

* Compressible
* Incompressible

---

### Compressible Resource Guarantees

Kubernetes only supports CPU at the moment.

* Pods are guaranteed to get the amount of CPU they request.
* Excess CPU resources will be distributed based on the amount of CPU requested.

---

### Incompressible Resource Guarantees

Kubernetes only supports memory at the moment.

* Pods will get the amount of memory they request.
* If They exceed their request, they may be killed.
  * only if another pod needs the memory.
* If Pods consume less memory than requested, they will not be killed
  * Except where QoS comes in to play.

---

### There are still some issues

* CPU isolation is at the container level.
* Cannot limit on Pod scale.
https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#planned-improvements

---

### Quality of Service

* In an overcommitted system, need to prioritise killing Containers.
* Mark containers as more/less important.
* Kill off less important Containers first.


---

### Quality of Service

There are three teirs, in decreasing order of priority.

* Guaranteed
* Burstable
* Best-Effort

---

### QoS - Guaranteed

* `limits` is non-zero, and set across all containers.
* `requests` are optionally defined and equal `limits`.
* last to be killed off.
```yaml
containers:
    name: foo
        resources:
            requests:
                cpu: 100m
                memory: 100Mi
            limits:
                cpu: 100m
                memory: 100Mi
```

---

### QoS - Burstable

* `requests` is non-zero, and set for at least one container.
* `limits` optionally defined and greater than `requests`.
* if `limits` is undefined, default to capacity of node.
```yaml
containers:
    name: foo
        resources:
            requests:
                cpu: 100m
                memory: 100Mi
```

---

### QoS - Best-Effort

* `requests` and `limits` are not set.
* uses whatever resources it needs/are available.
* first to be killed when resources become scarce.
```yaml
containers:
    name: foo
```

---

### Examples

Each container in a Pod may specify the amount of CPU it requests on a node.

CPU requests are used at schedule time, and represent a minimum amount of CPU that should be reserved for your container to run.

---

Let's demonstrate this concept using a simple container that will consume as much CPU as possible.
```
$ kubectl run cpustress --image=busybox --requests=cpu=100m \
-- md5sum /dev/urandom
```
This will create a single Pod on your node that requests 1/10 of a CPU, but it has no limit on how much CPU it may actually consume on the node.

---

To demonstrate this, you can use `kubectl top Pod <PODNAME>` to view the the used CPU shares.
(This may take some time before metrics are available)

```
$ kubectl get pods
NAME                         READY     STATUS    RESTARTS   AGE
cpustress-4101692926-zqw2p   1/1       Running   0          1m
$ kubectl top pod cpustress-4101692926-zqw2p
NAME                         CPU(cores)   MEMORY(bytes)
cpustress-4101692926-zqw2p   924m         0Mi
```

As you can see it uses 924m of a 1vCPU machine.

---

If you scale your application, we should see that each Pod is given an equal proportion of CPU time.

```
$ kubectl scale deployment cpustress --replicas=9
```
Once all the Pods are running, you will see that each Pod is getting approximately an equal proportion of CPU time.
Note: it can take a moment for the top output to reflect the changes to your deployment.

```
$ kubectl top pods
NAME                         CPU(cores)   MEMORY(bytes)
cpustress-1696410962-08kf9   314m         0Mi
cpustress-1696410962-r123x   310m         0Mi
cpustress-1696410962-5r61m   314m         0Mi
cpustress-1696410962-177g3   322m         0Mi
cpustress-1696410962-zgfqc   322m         0Mi
cpustress-1696410962-rh2kn   317m         0Mi
cpustress-1696410962-nmlvs   321m         0Mi
cpustress-1696410962-8gl26   318m         0Mi
cpustress-1696410962-fn667   317m         0Mi
```

Each container is getting 30% of the CPU time per their scheduling request, and we are unable to schedule more.

---

### Cleanup

```
$ kubectl delete deployment cpustress
```

---

### CPU limit

Setting a limit will allow you to control the maximum amount of CPU that your container can burst to.

```
$ kubectl run cpustress --image=busybox --requests=cpu=100m \
--limits=cpu=200m -- md5sum /dev/urandom
```

You can verify that by using `kubectl top pod`:
```
$ kubectl top pod
NAME                         CPU(cores)   MEMORY(bytes)
cpustress-1437538636-wkzh7   199m         0Mi
```

---

If you scale your application, we should see that each Pod is consuming a maximum of 200m CPU shares.

```
$ kubectl scale deployment cpustress --replicas=9
```
Once all the Pods are running, you will see that each Pod is getting approximately an equal proportion of CPU time.

```
$ kubectl top pod
NAME                         CPU(cores)   MEMORY(bytes)
cpustress-2801690769-895wj   198m         0Mi
cpustress-2801690769-735dt   198m         0Mi
cpustress-2801690769-gm9cz   199m         0Mi
cpustress-2801690769-ljt1w   199m         0Mi
cpustress-2801690769-wt54n   199m         0Mi
cpustress-2801690769-7c3tc   198m         0Mi
cpustress-2801690769-f2blv   199m         0Mi
cpustress-2801690769-7fm9n   201m         0Mi
cpustress-2801690769-6ssdk   198m         0Mi
```

---

### Memory requests

By default, a container is able to consume as much memory on the node as possible. It is recommended to specify the amount of memory your container will require to run.

Let's demonstrate this by creating a Pod that runs a single container which requests 100Mi of memory. The container will allocate and write to 200MB of memory every 2 seconds.

```
$ kubectl run memhog --image=derekwaynecarr/memhog --requests=memory=100Mi \
--command -- /bin/sh -c "while true; do memhog -r100 200m; sleep 1; done"
```

---

Verify the usage with `kubectl top pod`
```
$ kubectl top pod
NAME                     CPU(cores)   MEMORY(bytes)
memhog-328396322-dh03t   772m         200Mi
```

We request 100Mi, but have burst our memory usage to a greater value. That's called Burstable.

---

### Memory limits

If you specify a memory limit, you can constrain the amount of memory your container can use.

For example, let's limit our container to 200Mi of memory, and just consume 100MB.

```
$ kubectl run memhog --image=derekwaynecarr/memhog --limits=memory=200Mi \
--command -- /bin/sh -c "while true; do memhog -r100 100m; sleep 1; done"
```

```
$ kubectl top pod
NAME                      CPU(cores)   MEMORY(bytes)
memhog-4201114837-svfjl   632m         100Mi
```
As you can see we are only consuming 100MB on the node.

---

Let's demonstrate what happens if you exceed your allowed memory usage by creating a replication controller whose Pod will keep being OOM killed because it attempts to allocate 300MB of memory, but is limited to 200Mi.

```
$ kubectl run memhog-oom --image=derekwaynecarr/memhog --limits=memory=200Mi \
--command -- memhog -r100 300m
```

---

If we describe the created Pod, you will see that it keeps restarting until it goes into a CrashLoopBackOff.


```
$ kubectl get po
NAME                          READY     STATUS      RESTARTS   AGE
memhog-4201114837-svfjl       1/1       Running     0          11m
memhog-oom-3179143800-gmdbc   0/1       OOMKilled   5          3m
kubectl describe po memhog-oom-3179143800-gmdbc |grep -C 3 "Terminated"
      memory:        200Mi
    State:        Waiting
      Reason:        CrashLoopBackOff
    Last State:        Terminated
      Reason:        OOMKilled
      Exit Code:    137
      Started:        Mon, 20 Mar 2017 22:51:57 +0100

```

---

### What if my node runs out of memory?

With Guaranteed resources you are not in major danger of causing an OOM event on your node.

If any individual container consumes more than their specified limit, it will be killed.

---

With *BestEffort* and *Burstable* resources it is possible that a container will request more memory than is actually available on the node.

If this happens:
* The system will attempt to prioritize the containers that are killed based on their quality of service.
* This is done by using the OOMScoreAdjust feature in the Linux kernel
* Processes with lower values are preserved in favor of processes with higher values.
* The system daemons (kubelet, kube-proxy, docker) all run with low OOMScoreAdjust values.

---

Containers with *Guaranteed* memory are given a lower value than *Burstable* containers which have a lower value than *BestEffort* containers. As a consequence, containers with *BestEffort* should be killed before the other tier.

---

### Example

This may cause unwanted behavior...

```
$ kubectl run mem-guaranteed --image=derekwaynecarr/memhog --replicas=2 \
    --requests=cpu=10m --limits=memory=600Mi --command \
    -- memhog -r100000 500m
$ kubectl run mem-burstable --image=derekwaynecarr/memhog --replicas=2 \
    --requests=cpu=10m,memory=600Mi --command -- memhog -r100000 100m
$ kubectl run mem-besteffort --replicas=10 --image=derekwaynecarr/memhog \
    --requests=cpu=10m --command -- memhog -r10000 500m
```

---

This will force a SystemOOM.
```
$ kubectl get events | grep OOM
{kubelet gke-cluster-1-default-pool-312d7520-c4db}      System OOM encountered
```

The process relies on the Kernel to react to system OOM events. Depending on how the host operating system is configured, and which process the Kernel ultimately decides to kill on your Node, you may experience unstable results.

---

### Resource Quota

Quotas can be set per-namespace.
* Maximum request and limit across all Pods.
* Applies to each type of resource (CPU, mem).
* User must specify request or limit.
* Maximum number of a particular kind of object.
Ensure no user/app/department abuses the cluster.

Applied at admission time.

Pods which explicitly specify resource limits and requests will not pick up the namespace default values.

---

Before you continue ensure you have checked out the training materials repository:
```
$ git clone https://github.com/ContainerSolutions/$MODULE_NAME
```

Create a namespace
```
$ kubectl create namespace limit-example
```

Apply the LimitRange (`configs/limits.yaml`) to the new namespace
```
$ kubectl create -f configs/limits.yaml -n limit-example
```

```
apiVersion: v1
kind: LimitRange
metadata:
  name: mylimits
spec:
  limits:
  - max:
      cpu: "2"
      memory: 1Gi
    min:
      cpu: 200m
      memory: 6Mi
    type: Pod
  - default:
      cpu: 300m
      memory: 200Mi
    defaultRequest:
      cpu: 200m
      memory: 100Mi
    max:
      cpu: "2"
      memory: 1Gi
    min:
      cpu: 100m
      memory: 3Mi
    type: Container
```

---

Create a deployment in this namespace

```
$ kubectl run nginx --image=nginx --replicas=1 --namespace=limit-example
deployment "nginx" created
```

---

The default values of the namespace limit will be applied to this Pod
(`kubectl describe pod <pod_name> --namespace=<namespace_name>`)
```
$ kubectl get pods --namespace=limit-example
NAME                     READY     STATUS    RESTARTS   AGE
nginx-2371676037-tfncs   1/1       Running   0          4m
$ kubectl describe pod nginx-2371676037-tfncs --namespace=limit-example
...
Containers:
  nginx:
    Container ID:       docker://dece4453779a2664c045ea7edc21f41382d5b067552f5d41f4d65fa984628314
    Image:              nginx
    Image ID:           docker://sha256:e4e6d42c70b3f79c5d57c170526592168992eb3303a6594c439302fabd92d9a3
    Port:               <none>
    State:              Running
      Started:          Fri, 14 Jul 2017 12:44:24 +0000
    Ready:              True
    Restart Count:      0
    Limits:
      cpu:      300m
      memory:   200Mi
    Requests:
      cpu:              200m
      memory:           100Mi
...
```

---

If we deploy a Pod which exceeds the limit, using the following yaml:

```
apiVersion: v1
kind: Pod
metadata:
  name: invalid-pod
spec:
  containers:
  - name: kubernetes-serve-hostname
    image: gcr.io/google_containers/serve_hostname
    resources:
      limits:
        cpu: "3"
        memory: 100Mi
```
The output of kubectl is
```
$ kubectl create -f configs/invalid-cpu-pod.yaml -n limit-example
Error from server (Forbidden): error when creating "configs/invalid-cpu-pod.yaml": pods "invalid-pod" is forbidden: [maximum cpu usage per Pod is 2, but limit is 3., maximum cpu usage per Container is 2, but limit is 3.]
```

Clean up the resources used during this module:
```
$ kubectl delete --all pods
$ kubectl delete --all deployments
$ kubectl delete namespace limit-example
```

---

### Exercise - Add Resources to the Deals microservice

* Add resources section to the Deals Deployment
* Add limits and requests for memory and cpu
* Apply the udpated config

---  

[Next up advanced deployments...](./07_adv-deploys.md)

## Setup

---

### Connect to your VM

We will use `kubectl`, and `docker`. You have been provided with a VM containing the necessary tools.

Login to the provided cloud VM via ssh:

```
$ ssh csuser@111.111.111.111
csuser@111.111.111.111's password:
```

---

Retrieve config files: https://storage.googleapis.com/qcon-k8s/k8s-configs.tar.gz

```
$ wget https://storage.googleapis.com/qcon-k8s/k8s-configs.tar.gz
$ tar -xvf k8s-configs.tar.gz
```


Retrieve permissions file: 

```
$ wget https://storage.googleapis.com/qcon-k8s/ca.pem
```

---

Configure kubectl (Replace *user-x* and *password*)

```bash
$ kubectl config set-cluster workshop \
  --server=https://35.195.126.56 \
  --certificate-authority=/path/to/ca.pem
$ kubectl config set-credentials workshop-user \
  --username=user-X \
  --password=<password>
$ kubectl config set-context workshop \
  --cluster=workshop \
  --user=workshop-user \
  --namespace=user-X
$ kubectl config use-context workshop
```

---

### Verify you can access the cluster:

```
$ kubectl cluster-info
```

---

### Using the in-cluster Registry

We will use `port-forward` to expose the registry locally

```bash
$ POD=$(kubectl get pods --namespace kube-system -l k8s-app=kube-registry-upstream \
            -o template --template '{{range .items}}{{.metadata.name}} {{.status.phase}}{{"\n"}}{{end}}' \
            | grep Running | head -1 | cut -f1 -d' ')

$ kubectl port-forward --namespace kube-system $POD 5000:5000 &
```

Tag images as `localhost:5000/[user]/image`

---

Docker for Mac fix/hack

Add `docker.for.mac.localhost:5000` to Insecure Registries

---

[On to the first lab...](../01_intro.md)
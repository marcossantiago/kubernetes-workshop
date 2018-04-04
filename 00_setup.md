## Workshop Setup

---

We have provided everything you need to participate in this workshop.

We have provisioned you a Virtual Machine running a Kubernetes cluster.

The following steps will get you access to your cluster.

---

### Connect to your Virtual Machine

Open a Terminal window and login to the VM via SSH

```
$ ssh csuser@[PROVIDED IP]
```

Enter the password provided to you.

---

### Retrieve configuration files

Download cluster config files onto your VM

```
$ wget https://storage.googleapis.com/qcon-k8s/k8s-configs.tar.gz
$ tar -xvf k8s-configs.tar.gz
```

Retrieve the permissions file

```
$ wget https://storage.googleapis.com/qcon-k8s/ca.pem
```

---

### Configure your Kubernetes cluster

Execute the following commands to setup your cluster

```bash
$ kubectl config set-cluster workshop \
  --server=https://[CLUSTER IP] \
  --certificate-authority=ca.pem
```

```bash
$ kubectl config set-credentials workshop-user \
  --username=[PROVIDED USERNAME] \
  --password=[PROVIDED PASSWORD]
```

```bash
$ kubectl config set-context workshop \
  --cluster=workshop \
  --user=workshop-user \
  --namespace=[PROVIDED USERNAME]
```

```bash
$ kubectl config use-context workshop
```

---

### Verify the Configuration

Check the cluster is up and running correctly by retrieving information about it

```
$ kubectl cluster-info

Kubernetes master is running at https://...
....
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

### Docker for Mac fix

Add `docker.for.mac.localhost:5000` to Insecure Registries

---

[On to the first lab...](../01_intro.md)
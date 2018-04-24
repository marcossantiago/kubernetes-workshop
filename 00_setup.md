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

Retrieve the permissions file

```
$ wget https://storage.googleapis.com/goto-chicago/ca.pem
```

---

### Configure your Kubernetes cluster

Execute the following commands to setup your cluster

```bash
$ kubectl config set-cluster workshop \
  --server=https://35.184.157.223 \
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

[On to the first lab...](../01_intro.md)
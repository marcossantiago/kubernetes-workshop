## Setup

---

We will use `kubectl`, and `docker`.

You can either use your own laptop, or a provided cloud VM.

---

### Installing Locally

kubectl - https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-binary-via-curl

docker - https://www.docker.com/community-edition#/download

---

To view the slides locally:

```
docker run -d -p 8000:1948 -v $(pwd):/usr/src/app/ \
   containersol/reveal-md
```

Open a browser to localhost:8000

---

Retrieve permissions file: https://storage.googleapis.com/goto-k8s-workshop/ca.pem

Configure kubectl (Replace *user-x* and *password*)

```bash
$ kubectl config set-cluster workshop \
  --server=https://35.195.195.187 \
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
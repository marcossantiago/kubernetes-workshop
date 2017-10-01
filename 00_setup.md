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


---

Retrieve permissions file: https://storage.googleapis.com/goto-k8s-workshop/ca.pem

Configure kubectl (Replace *user-x* and *password*)

```bash
$ kubectl config set-cluster workshop --server=$endpoint \
--certificate-authority=/path/to/ca.pem
$ kubectl config set-credentials workshop-user --username=user-X \
--password=<password>
$ kubectl config set-context workshop --cluster=workshop \
--user=workshop-user --namespace=user-X
$ kubectl config use-context workshop
```

---

[On to the first lab...](../01_intro.md)
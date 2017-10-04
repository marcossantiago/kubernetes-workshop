## Setting up a 
## Kubernetes 
## Cluster

---

### In this section we will 

* Look at options for running Kubernetes
* Set up our own HA Kubernetes cluster
* Deploy Sock Shop
* Integrate our Deals service
* Visualize and monitor our application

---

There are **many** options for running Kubernetes

* Manual
* Terraform / Ansible
* kops
* kube-adm
* Tectonic
* Hosted (GKE, ACS)
* Openshift
* ...

---

### Depends on your needs

* Are you running in the cloud or onprem?
* Do you have dedicated infra/ops team?
* What are your security requirements?
* Do you like to do things the hard way :)?

---

### Exercise - setting up a 
### fully HA cluster

We will use `kops` to setup a Kubernetes cluster.

You can view the kops documentation on [GitHub](https://github.com/kubernetes/kops/)

---

Log-in to the Cloud VM

ssh csuser@[IP]

---

Some setup

```bash
$ sudo chown csuser /home/csuser
$ export KOPS_STATE_STORE=s3://goto-k8s
```

---

Generate an ssh keypair

```
$ mkdir .ssh
$ ssh-keygen -t rsa
Generating public/private rsa key pair.
Enter file in which to save the key (/root/.ssh/id_rsa):
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /root/.ssh/id_rsa.
Your public key has been saved in /root/.ssh/id_rsa.pub.
The key fingerprint is:
...
```

---

Create the cluster configuration

```
$ kops create cluster --zones eu-west-1c goto.user<X>.k8s.local
```

Or HA version

```
$ kops create cluster --node-count 3 --zones eu-west-1a,eu-west-1b,eu-west-1c --master-zones eu-west-1a,eu-west-1b,eu-west-1c --node-size t2.medium --master-size t2.small goto.user12.k8s.local
```

View and verify the output. You can modify the configuration by running

```
$ kops edit cluster goto.user<X>.k8s.local
```

---

Once you are happy with the config you can create the cluster by running

```
$ kops update cluster goto.user<X>.k8s.local --yes
```

This will take a few minutes to create... (coffee anyone?) we can then verify the cluster is up and healthy

```
$ kops validate cluster
```

And kubectl should be configured to point to our new cluster

```
$ kubectl cluster-info
```

---

### Run the Sock Shop

Read the documentation located at: https://microservices-demo.github.io/microservices-demo

```bash
$ git clone https://github.com/microservices-demo/microservices-demo.git

$ kubectl create ns sock-shop
$ kubectl apply -f microservices-demo/deploy/kubernetes/complete-demo.yaml
```

---

## Exercise #1

* Modify the front-end service of the Sock Shop to be of type LoadBalancer.

* Then access the application via the public dns name.

* Buy some socks!

---

## Exercise #2

Deploy the Deals service alongside the Sock Shop (hint: namespace)
* Add a deployment
* Add a service to expose it (on port 80)

Then update the front-end deployment to point to the image:  icrosby/front-end:v1

---

[Next up Monitoring...](../09_monitoring.md)

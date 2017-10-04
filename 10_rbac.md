## RBAC

---

### A bit about roles

Role Based Access Control (RBAC) is a common approach to managing users’ access to resources or operations.

Permissions specify exactly which resources and actions can be accessed. 

The basic principle is: instead of separately managing the permissions of each user, permissions are given to roles, which are then assigned to users, or better - groups of users.

---

### Roles Bundle Permissions

- Managing permissions per user can be a tedious task when many users are involved. 
- As users are added to the system, maintaining user permissions becomes harder and more prone to errors. 
- Incorrect assignment of permissions can block users’ access to required systems, or worse - allow unauthorized users to access restricted areas or perform risky operations.

---

* A regular user can only perform a limited number of actions (e.g. get, watch, list). 
* A closer look into these user actions can reveal that some actions tend to go together e.g. checking logs.
* Once roles are identified and assigned to each user, permissions can then be assigned to roles, instead of users. 

Managing the permissions of a small number of roles is a much easier task.

---

### Basic concepts

**Rule**: grants permission
* Applies to resource types
* Grants verbs (create, edit, view, delete)

* (Cluster)Role
  * Cluster wide / within a namespace
  * List of rules

* (Cluster)RoleBinding
  * Connects (Cluster)Role to User
  * Both human & service account

---

### API overview

The RBAC API declares four top-level types which will be covered in this section:
* Role
* ClusterRole
* RoleBinding
* ClusterRoleBinding

---

### Role
A `Role` contains rules that represent a set of permissions. Permissions are additive (there are no “deny” rules). 

A `Role` can be defined within a namespace, or cluster-wide (`Role` vs `ClusterRole`)

---

Here’s an example `Role` in the “default” namespace that can be used to grant read access to pods:

```
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

---

### Cluster Role

* A `ClusterRole` can be used to grant the same permissions as a `Role`, but because they are cluster-scoped, they can also be used to grant access to:
* cluster-scoped resources (like nodes)
* namespaced resources (like pods) across all namespaces

---

### Cluster Role

The following ClusterRole can be used to grant read access to secrets in any particular namespace, or across all namespaces

```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  # "namespace" omitted since ClusterRoles are not namespaced
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
```

---

### RoleBinding

A role binding grants the permissions defined in a role to a user

Permissions can be granted within a namespace with a `RoleBinding`, or cluster-wide with a `ClusterRoleBinding`.

---

A `RoleBinding` may reference a `Role` in the same namespace. The following `RoleBinding` grants the “pod-reader” role to the user “jane” within the “default” namespace. 

```
# This role binding allows "jane" to read pods in the "default" namespace.
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

---

In this example, even though the following RoleBinding refers to a `ClusterRole`, 
**dave** will only be able read secrets in the **development** namespace (the namespace of the RoleBinding).

```
# This role binding allows "dave" to read secrets in the "development" namespace.
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: read-secrets
  namespace: development # This only grants permissions within the "development" namespace.
subjects:
- kind: User
  name: dave
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

---

### Cluster Role Binding
A `ClusterRoleBinding` may be used to grant permission at the cluster level and in all namespaces. The following `ClusterRoleBinding` allows any user in the group “manager” to read secrets in any namespace.

```
# This cluster role binding allows anyone in the "manager" group to read secrets in any namespace.
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: read-secrets-global
subjects:
- kind: Group
  name: manager
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

---

### Refering to resources

Most resources are represented by a string, such as “pods”. However, some Kubernetes APIs involve a “subresource”, such as the logs for a pod.

In this case, “pods” is the namespaced resource, and “log” is a subresource of pods. Subresources can be represented with a slash:

```
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  namespace: default
  name: pod-and-pod-logs-reader
rules:
- apiGroups: [""]
  resources: ["pods/log"]
  verbs: ["get", "list"]

```

---

### Roles example

Allow reading the resource “pods” in the core API group:

```
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

---

Allow reading/writing “deployments” in both the “extensions” and “apps” API groups:

```
rules:
- apiGroups: ["extensions", "apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

---

### Role Binding Examples

Only the subjects section of a RoleBinding is shown in the following examples.

For a user named “alice@example.com”:

```
subjects:
- kind: User
  name: "alice@example.com"
  apiGroup: rbac.authorization.k8s.io
```

---

For a group named “frontend-admins”:

```
subjects:
- kind: Group
  name: "frontend-admins"
  apiGroup: rbac.authorization.k8s.io
```

---

### Example

We'll use 3 separate users. Some steps are taken to allow all users to work on the same logged in account. An alias will be created for each one so that it is easy to see who is taking the action.

* admin
* alice
* bob

---

Each of these corresponds to a different google account. Initially, we'll only use admin-kubectl, given that the other accounts don't have any permissions yet.

There are different ways to make authenticated entities, or subjects, work with the cluster. In this case, GCP Service Accounts (different from Kubernetes Service Accounts) will be used.

---

### prerequisites

** Requires gcloud cli tools and google account**

Verify the version of the cluster. The client version and the server version must be at least 1.6.

Disable the legacy authorization mechanism. It is safe to re-run this command if the legacy authorization mechanism has already been disabled.

```
gcloud beta container clusters update <cluster-name> \
    --no-enable-legacy-authorization
```

---

### Create an admin alias

```
gcloud config set account $admin
```

Create the admin-kubectl alias, an alias to kubectl that uses the token of the GCP project master account to authenticate.

```
export admin="<primary-email-address-of-gcp-project>"
alias admin-kubectl='kubectl --token="$(gcloud auth print-access-token --account=$admin)"'
```

---

### Creating the first service account

Creating First Service Account

Create the `admin-kubectl` alias, an alias to `kubectl` that uses a token associated with the new `cluster-user-1` GCP Service Account to authenticate.

---

* Create a GCP service account.
```sh
gcloud iam service-accounts create cluster-user-1 \
  --display-name=cluster-user-1
```
* Capture the full service account name.
```sh
alice=$(gcloud iam service-accounts list \
  --format='value(email)' --filter='displayName:cluster-user-1')
```
* Create a key for the GCP service account.
```sh
gcloud iam service-accounts keys create --iam-account \
  $alice cluster-user-1.json
```

---

* Use the GCP service account key to activate the service account.
```sh
gcloud auth activate-service-account $alice \
  --key-file=cluster-user-1.json
```
* Create an alias to make it easy to use `kubectl` authenticating as the new
   service account.
```sh
alias alice-kubectl='kubectl \
  --token="$(gcloud auth print-access-token \
  --account=$alice)"'
```
* Reset the active account to be ready for the next steps.
```sh
gcloud config set account $admin
```

---

### Creating the other service accounts

Create the remaining aliases, an alias to kubectl that uses a token associated with cluster-user-[2-3] GCP Service Account to authenticate.

This is just a repetition of the same steps for the second service account.

Note: Use new cluster-user

---

```
$ gcloud iam service-accounts create cluster-user-2 \
  --display-name=cluster-user-2
$ bob=$(gcloud iam service-accounts list \
  --format='value(email)' \
  --filter='displayName:cluster-user-2')
$ gcloud iam service-accounts keys create \
  --iam-account $bob cluster-user-2.json
$ gcloud auth activate-service-account $bob \
  --key-file=cluster-user-2.json
$ alias bob-kubectl='kubectl \
  --token="$(gcloud auth print-access-token \
  --account=$bob)"'
$ gcloud config set account $admin
```

---

Enable GCP IAM Cluster Admin Roles

In order for this new GCP Service Accounts to be able to do anything on clusters, they must have GCP IAM container engine permissions.

* Navigate to https://console.cloud.google.com
* Select the GCP project that contains your GKE cluster from the drop down list on the top.
* Expand the left menu and select IAM & Admin and the IAM.

...

---

* Click Add
* Enter the full email address of the user account that you are using as the admin user.
* Select Container, then Container Engine Admin from the Role menu.
* Add a second role by selecting Container, then Container Engine Cluster Admin from the Role menu.
* Then click Add to add the new IAM roles to your user

---

### Disable default authentication

The `~/.kube/config` file contains the configuration for your kubectl with all the authentication details. The aliases configured above use the `--token` parameter to authenticate. However, if there is a valid `auth-provider` section in the
`~/.kube/config` for your cluster, it will override the `--token` parameter and all requests will be authenticated using the settings `~/.kube/config`.

---

In order to test the settings, edit the `~/.kube/config` file and comment out the `auth-provider` section associated with your test cluster. It should look something like this:

```yaml
- name: gke_swisscom-bigdata_europe-west1-b_cluster-1
  user:
#    auth-provider:
#      config:
#        access-token: ya29.Gl1pBBvxyUpilEZPisyfplF4nYd6eVKmPBDvK21FgqBFqXryQF3lIQYhNqnRus-HLg6xgXzjDLxc3cr21_iNiEf2v3hygCh-X6ivwPjFahnvGhfy0UtINg41gVhGc2M
#        cmd-args: config config-helper --format=json
#        cmd-path: /Users/michael.mueller/google-cloud-sdk/bin/gcloud
#        expiry: 2017-06-14T15:06:36Z
#        expiry-key: '{.credential.token_expiry}'
#        token-key: '{.credential.access_token}'
#      name: gcp
```

---

### Validate

When the above configuration steps are complete, the admin alias should be able to list nodes, the other shouldn't be able to do anything.

```
admin-kubectl get nodes
NAME                                       STATUS    AGE       VERSION
gke-cluster-1-default-pool-75e0e5d2-4qwq   Ready     1h        v1.6.4
gke-cluster-1-default-pool-75e0e5d2-5qq1   Ready     1h        v1.6.4
gke-cluster-1-default-pool-75e0e5d2-xl86   Ready     1h        v1.6.4
```
```
bob-kubectl get nodes
Error from server (Forbidden): User "cluster-user-2@swisscom-bigdata.iam.gserviceaccount.com" cannot list nodes at the cluster scope.: "Required \"container.nodes.list\" permission." (get nodes)
```

---

### Example

We'll perform the following steps

- Creating namespaces
- Defining roles
- Creating role bindings
- Verifying that it is working

---

### Create the namespaces
Create the `production` and `test` namespaces with `admin-kubectl create namespace <NAME>`

#### Alice the `pod-reader`

Roles are resources in Kubernetes, just like Pods and Deployments. 
Typically, they are written in a text file, and applied to the Kubernetes cluster with `kubectl apply`.
Alternatively, you can use `kubectl create role`.

---

Let's create the `pod-reader` role first. This role will be able to list, get the details, create and delete deployments on the production namespace.

```
admin-kubectl create role pod-reader \
    --verb=get \
    --verb=list \
    --verb=watch \
    --resource=pods \
    --namespace=production
```

Note: If you get an error at this point:
```
Error from server (Forbidden): roles.rbac.authorization.k8s.io "pod-reader" is forbidden: attempt to grant extra privileges:
```
There is currently a known issue where IAM Service Accounts are not automatically granted cluster admin authorization. To correct the issue:
```
admin-kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole=cluster-admin \
    --user=$admin
```
---

Validate that Alice is not yet able to list pods in the production namespace

```
alice-kubectl get pods --namespace=production
Error from server (Forbidden): User "cluster-user-1@swisscom-bigdata.iam.gserviceaccount.com" cannot list pods in the namespace "production".: "Required \"container.pods.list\" permission." (get pods)
```

---

Let's create the role binding.

```
admin-kubectl create rolebinding alice-pod-reader-binding \
    --role=pod-reader \
    --user=$alice \
    --namespace=production
```

---

### Verifying that Alice can list pods
Alice should now be able to list pods.
```
alice-kubectl get pods --namespace=production
```

---

### Do it yourself

* Re-use the role created above and create a new roleBinding for Bob
* Bob should be able to do the same in the test namespace, but not in production

---

### Creating Cluster Roles

Additional permissions will be granted to Alice and Bob for specific namespaces using a common role.

Service Account Alice listing deployments in namespace production

alice-kubectl get deployments --namespace=production

As reading (checking) deployments is a common task and not specific to a namespace, let's create a clusterRole.
```
kubectl create clusterrole deployment-reader \
    --verb=get \
    --verb=list \
    --verb=watch \
    --resource=deployments
```

---

### Create a role binding

```
admin-kubectl create rolebinding alice-deployment-reader-binding \
    --clusterrole=deployment-reader \
    --user=$alice \
    --namespace=production
```

Validate:
```
alice-kubectl get deployments --namespace=production
```

---

### Do it yourself

* Create a roleBinding (deployment-reader) for Bob in the test namespace

---

### Using pre-defined roles

In Kubernetes exist some default roles. This roles can help to e.g. split cluster between multiple teams. Below are the steps needed to configure users with admin priviledges in their own namespaces, and viewing permissions in the other namespace.

The pre-configured cluster role to grant admin priviledge to a namespace is the admin cluster role.

To view the clusterroles issue:
```
admin-kubectl get clusterroles
```

---

### Creating namespace admins

Make the users administrators in their own namespaces by binding the cluster role to the user

```
admin-kubectl create rolebinding alice-admin \
    --clusterrole=admin \
    --user=$alice  \
    --namespace=production
admin-kubectl create rolebinding bob-admin \
    --clusterrole=admin \
    --user=$bob  \
    --namespace=test
```

---

### Create namespace viewers

All users should be allowed to view resources in other namespaces, this can also be achieved with default roles

```
admin-kubectl create rolebinding alice-view \
    --clusterrole=view \
    --user=$alice \
    --namespace=test
admin-kubectl create rolebinding bob-view \
    --clusterrole=view \
    --user=$bob \
    --namespace=production
```

### Cleanup

```
gcloud iam service-accounts delete cluster-user-1 -q
gcloud iam service-accounts delete cluster-user-2 -q
gcloud container clusters delete <cluster-name>
```

---

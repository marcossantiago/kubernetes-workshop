## Exercise

- Write Deployment config
- Write Service config
- Deploy our deals service on Kubernetes
- Call service
- Modify service, perform rolling update / rollback
- Integrate with Socks Shop

----

## A real microservice (v3)

`open v3/`

----

## Build new image

`cd v3/`
`docker build -t z2h-zurich/deals:v3 .`

----

## Add Kubernetes Deployments

`vi deals-dep.yaml`

```bash
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: deals
  labels:
    name: deals
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: deals
    spec:
      containers:
      - name: deals
        image: z2h-zurich/deals:v3
        imagePullPolicy: Never
        ports:
        - containerPort: 8080
```

----

`vi deals-db-dep.yaml'

```bash
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: deals-db
  labels:
    name: deals-db
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: deals-db
    spec:
      containers:
      - name: deals-db
        image: mongo
        ports:
        - name: mongo
          containerPort: 27017
```

----

## Add Kubernetes Services

`vi deals-svc.yaml`

```bash
apiVersion: v1
kind: Service
metadata:
  name: deals
  labels:
    name: deals
spec:
  type: NodePort
  ports:
    # the port that this service should serve on
  - port: 8081
    targetPort: 8080
  selector:
    name: deals
   type: NodePort
```

----

`vi deals-db-svc.yaml`

```bash
apiVersion: v1
kind: Service
metadata:
  name: deals-db
  labels:
    name: deals-db
spec:
  ports:
  - port: 27017
    targetPort: 27017
  selector:
    name: deals-db
```

----

## Deploy

may need `sudo` or `--validate=false`

```bash
kubectl create -f deals-dep.yaml
kubectl create -f deals-svc.yaml
kubectl create -f deals-db-dep.yaml
kubectl create -f deals-db-svc.yaml
kubectl get pods -w
```

----

## Call service

Find port

`kubectl describe service deals`

`curl localhost:[port]/deals?id=1`

----

## Logging and Monitoring

`open v4/`

----

## Build new version

`cd v4/`
`docker build -t z2h-zurich/deals:v4 .`

----

## Update k8s deploy

`vi deals-dep-v2.yaml`

```bash
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: deals
  labels:
    name: deals
spec:
  replicas: 3
  strategy:
    rollingUpdate:
      maxUnavailable: 1
  template:
    metadata:
      labels:
        name: deals
    spec:
      containers:
      - name: deals
        image: z2h-zurich/deals:v3
        ports:
        - containerPort: 8080
```

----

Apply changes and watch pods...

`kubectl apply -f deals-dep-v2.yaml`
`kubectl get pods -w`

----

## Update to latest image

`kubectl set image deployment/deals deals=z2h-zurich/deals:v4`
`kubectl rollout status deployment/deals`
`kubectl get pods`

----

## Rollback?

`kubectl rollout undo deployment/deals`
`kubectl rollout status deployment/deals`
`kubectl get pods`

----

## Launch Sock Shop

`kubectl apply -f complete-demo.yaml --validate=false`

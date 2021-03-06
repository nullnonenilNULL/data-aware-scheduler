# Data Aware Scheduler (PoC)

This Scheduler works (and is tested) with Kubernetes 1.5

## Dependency Managment

For the dependency managment was (dep)[https://github.com/golang/dep] used.

```bash
$ dep ensure -update
```

## Build

To build the binary (for linux):

```
$ GOOS=linux go build -o scheduler -a --ldflags '-extldflags "-static" -s -w' -tags netgo -installsuffix netgo .
```

to build a Docker container:

```
$ VERSION=0.0.2
$ docker build -t johscheuer/data-aware-scheduler:$VERSION .
$ docker push johscheuer/data-aware-scheduler:$VERSION
```

## Run it 

### Prerequisites

- Running Quobyte Cluster (you can use this [tool](https://github.com/inovex/quobyte-kubernetes-operator) for a quick start)
- Running Kubernetes Cluster

### Inside Kubernetes

Adjust config `deployments/scheduler_config.yaml` to your need:

```
kubectl create --namespace=kube-system -f deployments/scheduler_config.yaml
kubectl create --namespace=kube-system -f deployments/scheduler.yaml
```

### Outside Kubernetes

At the moment you need a local mount of Quobyte with xattr enabled (`-o user_xattr`). As soon as the same information is available over the API this requierement drops.

Copy the example configuration and adjust it to your needs:

```bash
copy config.yaml.example config.yaml
```

Now start the Scheduler:

```bash
./scheduler --config=./config.yaml

```

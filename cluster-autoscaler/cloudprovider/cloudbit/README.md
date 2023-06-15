# Cluster Autoscaler for Cloudbit

The cluster autoscaler for Cloudbit Cloud scales worker nodes within any specified Cloudbit Kubernetes cluster.

# Configuration
As there is no concept of a node group within Cloudbit's Kubernetes offering, the configuration required is quite 
simple. You need to set:

- `CLOUDBIT_CLUSTER_ID`: the ID of the cluster (a UUID)
- `CLOUDBIT_API_KEY`: the Cloudbit access token literally defined
- `CLOUDBIT_API_URL`: the Cloudbit URL (optional; defaults to `https://api.cloudbit.ch/`)
- The minimum and maximum number of **worker** nodes you want (the master is excluded)

See the [cluster-autoscaler-standard.yaml](examples/cluster-autoscaler-standard.yaml) example configuration, but to 
summarise you should set a `nodes` startup parameter for cluster autoscaler to specify a node group called `workers` 
e.g. `--nodes=1:10:workers`.

# Development

Make sure you're inside the root path of the [autoscaler
repository](https://github.com/kubernetes/autoscaler)

1.) Build the `cluster-autoscaler` binary:


```
make build-in-docker
```

2.) Build the docker image:

```
docker build -t cloudbit/cluster-autoscaler:dev .
```

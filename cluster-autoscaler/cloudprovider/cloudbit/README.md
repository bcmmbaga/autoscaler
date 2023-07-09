# Cluster Autoscaler for Cloudbit

The cluster autoscaler for Cloudbit Cloud scales worker nodes within any specified Cloudbit Kubernetes cluster.

# Configuration
As there is no concept of a node group within Cloudbit's Kubernetes offering, the configuration required is quite 
simple. You need to set:

- The ID of the cluster
- The Cloudbit access token literally defined
- The Cloudbit URL (optional; defaults to `https://api.cloudbit.ch/`)
- The minimum and maximum number of **worker** nodes you want (the master is excluded)

See the [cluster-autoscaler-standard.yaml](examples/cluster-autoscaler-standard.yaml) example configuration, but to 
summarise you should set a `nodes` startup parameter for cluster autoscaler to specify a node group called `workers` 
e.g. `--nodes=3:10:workers`.

The remaining parameters can be set via environment variables (`CLOUDBIT_CLUSTER_ID`, `CLOUDBIT_API_TOKEN` and `CLOUDBIT_API_URL`) as in the
example YAML.

It is also possible to get these parameters through a YAML file mounted into the container
(for example via a Kubernetes Secret). The path configured with a startup parameter e.g.
`--cloud-config=/etc/kubernetes/cloud.config`. In this case the YAML keys are `api_url`, `api_token` and `cluster_id`.


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

3.) Push the docker image to Docker hub:

```
docker push cloudbit/cluster-autoscaler:dev
```
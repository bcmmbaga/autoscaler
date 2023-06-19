/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"
	"io"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/klog/v2"
)

var (
	version = "dev"
)

type nodeGroupClient interface {
	// List lists all the node pools found in a Kubernetes cluster.
	List(ctx context.Context, cursor goclient.Cursor) (kubernetes.NodeList, error)

	// PerformAction updates the details of an existing node pool.
	PerformAction(ctx context.Context, nodeID int, req kubernetes.NodePerformAction) (kubernetes.Node, error)

	// Delete deletes a specific node in a node pool.
	Delete(ctx context.Context, nodeID int) error
}

// Manager handles Cloudbit communication and data caching of
// node groups
type Manager struct {
	client        nodeGroupClient
	clusterID     int
	nodeGroups    []*NodeGroup
	discoveryOpts cloudprovider.NodeGroupDiscoveryOptions
}

// Config is the configuration of the Cloudbit cloud provider
type Config struct {
	// ClusterID is the id associated with the cluster where Cloudbit
	// Cluster Autoscaler is running.
	ClusterID int `json:"cluster_id"`

	// Token is the User's Access Token associated with the cluster where
	// Cloudbit Cluster Autoscaler is running.
	Token string `json:"token"`

	// URL points to Cloudbit API. If empty, defaults to
	// https://api.cloudbit.ch/
	ApiURL string `json:"api_url"`
}

func newManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*Manager, error) {
	cfg := &Config{}
	if configReader != nil {
		body, err := io.ReadAll(configReader)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Token == "" {
		return nil, errors.New("access token is not provided")
	}
	if cfg.ClusterID == 0 {
		return nil, errors.New("cluster ID is not provided")
	}

	opts := []goclient.Option{}
	if cfg.ApiURL != "" {
		opts = append(opts, goclient.WithBase(cfg.ApiURL))
	}
	opts = append(opts, goclient.WithUserAgent("cluster-autoscaler-cloudbit/"+version))
	opts = append(opts, goclient.WithToken(cfg.Token))

	doClient := goclient.NewClient(opts...)
	kubernetes.NewNodeService(doClient, cfg.ClusterID)

	m := &Manager{
		client:        kubernetes.NewNodeService(doClient, cfg.ClusterID),
		clusterID:     cfg.ClusterID,
		nodeGroups:    make([]*NodeGroup, 0),
		discoveryOpts: discoveryOpts,
	}

	return m, nil
}

// Refresh refreshes the cache holding the nodegroups. This is called by the CA
// based on the `--scan-interval`. By default it's 10 seconds.
func (m *Manager) Refresh() error {
	var (
		minSize           int
		maxSize           int
		workerConfigFound = false
		poolConfigFound   = false
		poolGroups        []*NodeGroup
		workerGroups      []*NodeGroup
	)

	ctx := context.Background()
	pools, err := m.client.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return fmt.Errorf("couldn't list Kubernetes cluster pools: %s", err)
	}

	klog.V(4).Infof("refreshing workers node group kubernetes cluster: %q", m.clusterID)

	for _, specString := range m.discoveryOpts.NodeGroupSpecs {
		spec, err := dynamic.SpecFromString(specString, true)
		if err != nil {
			return fmt.Errorf("failed to parse node group spec: %v", err)
		}

		if spec.Name == "workers" {
			minSize = spec.MinSize
			maxSize = spec.MaxSize
			workerConfigFound = true
			klog.V(4).Infof("found configuration for workers node group: min: %d max: %d", minSize, maxSize)
		} else {
			poolConfigFound = true
			pool := m.getNodeGroupConfig(spec, pools.Items)
			if pool != nil {
				poolGroups = append(poolGroups, pool)
			}
			klog.V(4).Infof("found configuration for pool node group: min: %d max: %d", minSize, maxSize)
		}
	}

	if poolConfigFound {
		m.nodeGroups = poolGroups
	} else if workerConfigFound {
		for _, nodePool := range pools.Items {
			np := nodePool
			klog.V(4).Infof("adding node pool: %q", nodePool.ID)

			workerGroups = append(workerGroups, &NodeGroup{
				id:        nodePool.ID,
				clusterID: m.clusterID,
				client:    m.client,
				nodePool:  &np,
				minSize:   minSize,
				maxSize:   maxSize,
			})
		}
		m.nodeGroups = workerGroups
	} else {
		return fmt.Errorf("no workers node group configuration found")
	}

	// If both config found, pool config get precedence
	if poolConfigFound && workerConfigFound {
		m.nodeGroups = poolGroups
	}

	if len(m.nodeGroups) == 0 {
		klog.V(4).Info("cluster-autoscaler is disabled. no node pools are configured")
	}

	return nil
}

// getNodeGroupConfig get the node group configuration from the cluster pool configuration
func (m *Manager) getNodeGroupConfig(spec *dynamic.NodeGroupSpec, pools []kubernetes.Node) *NodeGroup {
	for _, nodePool := range pools {
		if spec.Name == nodePool.Name {
			np := nodePool
			klog.V(4).Infof("adding node pool: %q min: %d max: %d", nodePool.ID, spec.MinSize, spec.MaxSize)

			return &NodeGroup{
				id:        nodePool.ID,
				clusterID: m.clusterID,
				client:    m.client,
				nodePool:  &np,
				minSize:   spec.MinSize,
				maxSize:   spec.MaxSize,
			}
		}
	}
	return nil
}

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
	"errors"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/kubernetes"
	autoscaler "k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/klog/v2"
	schedulerframework "k8s.io/kubernetes/pkg/scheduler/framework"
)

type NodeGroup struct {
	id        int
	clusterID int
	client    nodeGroupClient
	nodes     []kubernetes.Node

	minSize    int
	maxSize    int
	getOptions *autoscaler.NodeGroupAutoscalingOptions
}

// MaxSize returns maximum size of the node group.
func (n *NodeGroup) MaxSize() int {
	return n.maxSize
}

// MinSize returns minimum size of the node group.
func (n *NodeGroup) MinSize() int {
	return n.minSize
}

// GetOptions returns the options used to create this node group.
func (n *NodeGroup) GetOptions(autoscaler.NodeGroupAutoscalingOptions) (*autoscaler.NodeGroupAutoscalingOptions, error) {
	return n.getOptions, nil
}

// TargetSize returns the current target size of the node group. It is possible
// that the number of nodes in Kubernetes is different at the moment but should
// be equal to Size() once everything stabilizes (new nodes finish startup and
// registration or removed nodes are deleted completely). Implementation
// required.
func (n *NodeGroup) TargetSize() (int, error) {
	return len(n.nodes), nil
}

// IncreaseSize increases the size of the node group. To delete a node you need
// to explicitly name it and use DeleteNode. This function should wait until
// node group size is updated. Implementation required.
func (n *NodeGroup) IncreaseSize(delta int) error {
	if delta <= 0 {
		return fmt.Errorf("delta must be positive, have: %d", delta)
	}

	targetSize := len(n.nodes) + delta

	if targetSize > n.MaxSize() {
		return fmt.Errorf("size increase is too large. current: %d desired: %d max: %d",
			len(n.nodes), targetSize, n.MaxSize())
	}

	req := kubernetes.ClusterUpdateFlavor{
		Worker: kubernetes.ClusterWorkerUpdate{
			ProductID: n.nodeProductID(),
			Count:     targetSize,
		},
	}

	ctx := context.Background()
	cluster, err := n.client.UpdateCluster(ctx, n.clusterID, req)
	if err != nil {
		return err
	}

	if cluster.NodeCount.Expected.Worker != targetSize {
		return fmt.Errorf("couldn't increase size to %d (delta: %d). Current size is: %d",
			targetSize, delta, cluster.NodeCount.Expected.Worker)
	}

	return nil
}

// DeleteNodes deletes nodes from this node group (and also increasing the size
// of the node group with that). Error is returned either on failure or if the
// given node doesn't belong to this node group. This function should wait
// until node group size is updated. Implementation required.
func (n *NodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	for _, node := range nodes {
		instanceID, err := toNodeID(node.Spec.ProviderID)
		if err != nil {
			return fmt.Errorf("deleting node failed for cluster: %q node pool: %q node: %q: %s",
				n.clusterID, n.id, node.Name, err)
		}

		klog.V(4).Info("deleting node: %q", instanceID)

		ctx := context.Background()
		err = n.client.DeleteClusterNode(ctx, instanceID)
		if err != nil {
			return fmt.Errorf("deleting node failed for cluster: %q node pool: %q node: %q: %s",
				n.clusterID, n.id, node.Name, err)
		}

	}

	return nil
}

// DecreaseTargetSize decreases the target size of the node group. This function
// doesn't permit to delete any existing node and can be used only to reduce the
// request for new nodes that have not been yet fulfilled. Delta should be negative.
// It is assumed that cloud provider will not delete the existing nodes when there
// is an option to just decrease the target. Implementation required.
func (n *NodeGroup) DecreaseTargetSize(delta int) error {
	if delta >= 0 {
		return fmt.Errorf("delta must be negative, have: %d", delta)
	}

	targetSize := len(n.nodes) + delta
	if targetSize < n.MinSize() {
		return fmt.Errorf("size decrease is too small. current: %d desired: %d min: %d",
			len(n.nodes), targetSize, n.MinSize())
	}

	req := kubernetes.ClusterUpdateFlavor{
		Worker: kubernetes.ClusterWorkerUpdate{
			ProductID: n.nodeProductID(),
			Count:     targetSize,
		},
	}

	ctx := context.Background()
	cluster, err := n.client.UpdateCluster(ctx, n.clusterID, req)
	if err != nil {
		return err
	}

	if cluster.NodeCount.Expected.Worker != targetSize {
		return fmt.Errorf("couldn't decrease size to %d (delta: %d). Current size is: %d",
			targetSize, delta, cluster.NodeCount.Expected.Worker)
	}

	return nil
}

// Id returns an unique identifier of the node group.
func (n *NodeGroup) Id() string {
	return fmt.Sprintf("%d", n.id)
}

// Debug returns a string containing all information regarding this node group.
func (n *NodeGroup) Debug() string {
	return fmt.Sprintf("cluster ID: %s (min:%d max:%d)", n.Id(), n.MinSize(), n.MaxSize())
}

// Nodes returns a list of all nodes that belong to this node group.  It is
// required that Instance objects returned by this method have Id field set.
// Other fields are optional.
func (n *NodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	if n.nodes == nil {
		return nil, errors.New("node pool instance is not created")
	}

	return toInstances(n.nodes), nil
}

// TemplateNodeInfo returns a schedulernodeinfo.NodeInfo structure of an empty
// (as if just started) node. This will be used in scale-up simulations to
// predict what would a new node look like if a node group was expanded. The
// returned NodeInfo is expected to have a fully populated Node object, with
// all the labels, capacity and allocatable information as well as all pods
// that are started on the node by default, using manifest (most likely only
// kube-proxy). Implementation optional.
func (n *NodeGroup) TemplateNodeInfo() (*schedulerframework.NodeInfo, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// Exist checks if the node group really exists on the cloud provider side.
// Allows to tell the theoretical node group from the real one. Implementation
// required.
func (n *NodeGroup) Exist() bool {
	return n.nodes != nil
}

// Create creates the node group on the cloud provider side. Implementation
// optional.
func (n *NodeGroup) Create() (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// Delete deletes the node group on the cloud provider side.  This will be
// executed only for autoprovisioned node groups, once their size drops to 0.
// Implementation optional.
func (n *NodeGroup) Delete() error {
	return cloudprovider.ErrNotImplemented
}

// Autoprovisioned returns true if the node group is autoprovisioned. An
// autoprovisioned group was created by CA and can be deleted when scaled to 0.
func (n *NodeGroup) Autoprovisioned() bool {
	return false
}

// nodeProductID returns the cloudbit product id from existing worker node.
func (n *NodeGroup) nodeProductID() int {
	if totalNodes := len(n.nodes); totalNodes > 0 {
		return n.nodes[totalNodes-1].Product.ID
	}
	return 0
}

// toInstances converts a slice of []*kubernetes.Node to
// cloudprovider.Instance
func toInstances(nodes []kubernetes.Node) []cloudprovider.Instance {
	instances := make([]cloudprovider.Instance, 0, len(nodes))
	for _, node := range nodes {
		instances = append(instances, toInstance(node))
	}
	return instances
}

// toInstance converts the given *kubernetes.Node to a
// cloudprovider.Instance
func toInstance(node kubernetes.Node) cloudprovider.Instance {
	return cloudprovider.Instance{
		Id:     toProviderID(node.ID),
		Status: toInstanceStatus(&node.Status),
	}
}

// toInstanceStatus converts the given *godo.KubernetesNodeStatus to a
// cloudprovider.InstanceStatus
func toInstanceStatus(nodeState *kubernetes.NodeStatus) *cloudprovider.InstanceStatus {
	if nodeState == nil {
		return nil
	}

	st := &cloudprovider.InstanceStatus{}
	switch nodeState.Key {
	case "creating":
		st.State = cloudprovider.InstanceCreating
	case "healthy":
		st.State = cloudprovider.InstanceRunning
	case "draining", "deleting":
		st.State = cloudprovider.InstanceDeleting
	default:
		st.ErrorInfo = &cloudprovider.InstanceErrorInfo{
			ErrorClass:   cloudprovider.OtherErrorClass,
			ErrorCode:    "no-code-cloudbit",
			ErrorMessage: nodeState.Key,
		}
	}

	return st
}

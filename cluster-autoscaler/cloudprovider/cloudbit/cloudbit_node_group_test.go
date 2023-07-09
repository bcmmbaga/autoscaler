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
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/kubernetes"
	"testing"
)

func TestNodeGroup_TargetSize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		numberOfNodes := 2

		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)
		size, err := ng.TargetSize()
		assert.NoError(t, err)
		assert.Equal(t, numberOfNodes, size, "target size is not correct")
	})
}

func TestNodeGroup_IncreaseSize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		numberOfNodes := 3
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := 2

		newCount := numberOfNodes + delta
		client.On("UpdateCluster",
			context.Background(),
			ng.clusterID,
			kubernetes.ClusterUpdateFlavor{
				Worker: kubernetes.ClusterWorkerUpdate{
					ProductID: 0,
					Count:     newCount,
				},
			},
		).Return(
			kubernetes.Cluster{
				ID: 1,
				NodeCount: struct {
					Current struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"current"`
					Expected struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"expected"`
				}{
					Expected: struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					}{
						Worker: 5,
					},
				},
			},
			nil,
		).Once()

		err := ng.IncreaseSize(delta)
		assert.NoError(t, err)
	})

	t.Run("successful increase to maximum", func(t *testing.T) {
		// Increase from 3 nodes to 4 (but 2 worker nodes which is the max)
		numberOfNodes := 3
		maxNodes := 4

		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, maxNodes)

		delta := 1

		newCount := numberOfNodes + delta
		client.On("UpdateCluster",
			context.Background(),
			ng.clusterID,
			kubernetes.ClusterUpdateFlavor{
				Worker: kubernetes.ClusterWorkerUpdate{
					ProductID: 0,
					Count:     newCount,
				},
			},
		).Return(
			kubernetes.Cluster{
				ID: 1,
				NodeCount: struct {
					Current struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"current"`
					Expected struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"expected"`
				}{
					Expected: struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					}{
						Worker: 4,
					},
				},
			},
			nil,
		).Once()

		err := ng.IncreaseSize(delta)
		assert.NoError(t, err)
	})

	t.Run("negative increase", func(t *testing.T) {
		numberOfNodes := 3
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := -1
		err := ng.IncreaseSize(delta)
		exp := fmt.Errorf("delta must be positive, have: %d", delta)
		assert.EqualError(t, err, exp.Error(), "size increase must be positive")
	})

	t.Run("zero increase", func(t *testing.T) {
		numberOfNodes := 3
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := 0
		err := ng.IncreaseSize(delta)
		exp := fmt.Errorf("delta must be positive, have: %d", delta)
		assert.EqualError(t, err, exp.Error(), "size increase must be positive")
	})

	t.Run("large increase above maximum", func(t *testing.T) {
		numberOfNodes := 95
		maxNodes := 100
		delta := 10
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, maxNodes)

		exp := fmt.Errorf("size increase is too large. current: %d desired: %d max: %d",
			numberOfNodes, numberOfNodes+delta, ng.MaxSize())
		err := ng.IncreaseSize(delta)
		assert.EqualError(t, err, exp.Error(), "size increase is too large")
	})
}

func TestNodeGroup_DecreaseTargetSize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		numberOfNodes := 5

		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := -2

		newCount := numberOfNodes + delta
		client.On("UpdateCluster",
			context.Background(),
			ng.clusterID,
			kubernetes.ClusterUpdateFlavor{
				Worker: kubernetes.ClusterWorkerUpdate{
					ProductID: 0,
					Count:     newCount,
				},
			},
		).Return(
			kubernetes.Cluster{
				ID: 1,
				NodeCount: struct {
					Current struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"current"`
					Expected struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					} `json:"expected"`
				}{
					Expected: struct {
						ControlPlane int `json:"control-plane"`
						Worker       int `json:"worker"`
					}{
						Worker: 3,
					},
				},
			},
			nil,
		).Once()

		err := ng.DecreaseTargetSize(delta)
		assert.NoError(t, err)
	})

	t.Run("positive decrease", func(t *testing.T) {
		numberOfNodes := 5
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := 1
		err := ng.DecreaseTargetSize(delta)

		exp := fmt.Errorf("delta must be negative, have: %d", delta)
		assert.EqualError(t, err, exp.Error(), "size decrease must be negative")
	})

	t.Run("zero decrease", func(t *testing.T) {
		numberOfNodes := 5
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 10)

		delta := 0
		exp := fmt.Errorf("delta must be negative, have: %d", delta)

		err := ng.DecreaseTargetSize(delta)
		assert.EqualError(t, err, exp.Error(), "size decrease must be negative")
	})

	t.Run("small decrease below minimum", func(t *testing.T) {
		delta := -2
		numberOfNodes := 2
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, numberOfNodes), 3, 5)

		exp := fmt.Errorf("size decrease is too small. current: %d desired: %d min: %d",
			numberOfNodes, numberOfNodes+delta, ng.MinSize())
		err := ng.DecreaseTargetSize(delta)
		assert.EqualError(t, err, exp.Error(), "size decrease is too small")
	})
}

func TestNodeGroup_DeleteNodes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, 3), 3, 10)

		nodes := []*apiv1.Node{
			{Spec: apiv1.NodeSpec{ProviderID: "flow://1"}},
			{Spec: apiv1.NodeSpec{ProviderID: "flow://2"}},
			{Spec: apiv1.NodeSpec{ProviderID: "flow://3"}},
		}

		// this should be called three times (the number of nodes)
		client.On("DeleteClusterNode",
			context.Background(),
			1,
		).Return(
			nil,
		).Once()

		client.On("DeleteClusterNode",
			context.Background(),
			2,
		).Return(
			nil,
		).Once()

		client.On("DeleteClusterNode",
			context.Background(),
			3,
		).Return(
			nil,
		).Once()

		err := ng.DeleteNodes(nodes)
		assert.NoError(t, err)
	})

	t.Run("client deleting node fails", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, 3), 3, 10)

		nodes := []*apiv1.Node{
			{Spec: apiv1.NodeSpec{ProviderID: "flow://1"}},
			{Spec: apiv1.NodeSpec{ProviderID: "flow://2"}},
			{Spec: apiv1.NodeSpec{ProviderID: "flow://3"}},
		}

		// this should be called three times (the number of nodes)
		client.On("DeleteClusterNode",
			context.Background(),
			1,
		).Return(
			nil,
		).Once()

		client.On("DeleteClusterNode",
			context.Background(),
			2,
		).Return(
			errors.New("random error"),
		).Once()

		err := ng.DeleteNodes(nodes)
		assert.Error(t, err)
	})
}

func TestNodeGroup_Nodes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(
			client,
			[]kubernetes.Node{
				{
					ID: 1,
					Status: kubernetes.NodeStatus{
						ID:  1,
						Key: "healthy",
					},
				},
				{
					ID: 2,
					Status: kubernetes.NodeStatus{
						ID:  2,
						Key: "creating",
					},
				},
				{
					ID: 3,
					Status: kubernetes.NodeStatus{
						ID:  3,
						Key: "draining",
					},
				},
				{
					ID: 4,
					Status: kubernetes.NodeStatus{
						ID:  4,
						Key: "random",
					},
				},
			},
			3,
			10)

		exp := []cloudprovider.Instance{
			{
				Id: "flow://1",
				Status: &cloudprovider.InstanceStatus{
					State: cloudprovider.InstanceRunning,
				},
			},
			{
				Id: "flow://2",
				Status: &cloudprovider.InstanceStatus{
					State: cloudprovider.InstanceCreating,
				},
			},
			{
				Id: "flow://3",
				Status: &cloudprovider.InstanceStatus{
					State: cloudprovider.InstanceDeleting,
				},
			},
			{
				Id: "flow://4",
				Status: &cloudprovider.InstanceStatus{
					ErrorInfo: &cloudprovider.InstanceErrorInfo{
						ErrorClass:   cloudprovider.OtherErrorClass,
						ErrorCode:    "no-code-cloudbit",
						ErrorMessage: "random",
					},
				},
			},
		}

		nodes, err := ng.Nodes()
		assert.NoError(t, err)
		assert.Equal(t, exp, nodes, "nodes do not match")
	})

	t.Run("failure (nil node pool)", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, nil, 3, 10)

		_, err := ng.Nodes()
		assert.Error(t, err, "Nodes() should return an error")
	})
}

func TestNodeGroup_Debug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, 2), 3, 200)
		d := ng.Debug()
		exp := "cluster ID: 1 (min:3 max:200)"
		assert.Equal(t, exp, d, "debug string do not match")
	})
}

func TestNodeGroup_Exist(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, make([]kubernetes.Node, 3), 3, 200)
		exist := ng.Exist()
		assert.Equal(t, true, exist, "node group should exist")
	})

	t.Run("failure", func(t *testing.T) {
		client := &cloudbitClientMock{}
		ng := testNodeGroup(client, nil, 1, 200)
		exist := ng.Exist()
		assert.Equal(t, false, exist, "node group should not exist")
	})
}

func testNodeGroup(client nodeGroupClient, nodes []kubernetes.Node, min int, max int) *NodeGroup {
	return &NodeGroup{
		id:        1,
		clusterID: 1,
		client:    client,
		nodes:     nodes,
		minSize:   min,
		maxSize:   max,
	}
}

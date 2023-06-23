package cloudbit

import (
	"bytes"
	"context"
	"errors"
	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"
	"github.com/stretchr/testify/assert"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"testing"
)

func TestNewManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.cloudbit.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)
		assert.Equal(t, 123456, manager.clusterID, "cluster ID does not match")
		assert.Equal(t, nodeGroupDiscoveryOptions, manager.discoveryOpts, "node group discovery options do not match")
	})

	t.Run("empty api_token", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "", "api_url": "https://api.cloudbit.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		_, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.EqualError(t, err, errors.New("cloudbit access token is not provided").Error())
	})

	t.Run("empty cluster ID", func(t *testing.T) {
		cfg := `{"api_token": "123-123-123", "api_url": "https://api.cloudbit.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		_, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.EqualError(t, err, errors.New("cloudbit cluster ID is not provided").Error())
	})
}

func TestCloudbitManager_Refresh(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.cloudbit.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)

		client := &cloudbitClientMock{}
		ctx := context.Background()
		cursor := goclient.Cursor{NoFilter: 1}

		client.On("List", ctx, cursor).Return(
			kubernetes.NodeList{
				Items: []kubernetes.Node{
					{
						ID:   1,
						Name: "worker1",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   1,
						Name: "worker2",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
				},
			},
			nil,
		).Once()

		manager.client = client
		err = manager.Refresh()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(manager.nodeGroups), "number of node groups do not match")
	})
}

func TestCloudbitManager_RefreshWithNodeSpec(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := `{"cluster_id": 123456, "api_token": "123-123-123", "api_url": "https://api.cloudbit.ch/"}`

		nodeGroupSpecs := []string{"1:10:workers"}
		nodeGroupDiscoveryOptions := cloudprovider.NodeGroupDiscoveryOptions{NodeGroupSpecs: nodeGroupSpecs}

		manager, err := newManager(bytes.NewBufferString(cfg), nodeGroupDiscoveryOptions)
		assert.NoError(t, err)

		client := &cloudbitClientMock{}
		ctx := context.Background()
		cursor := goclient.Cursor{NoFilter: 1}

		client.On("List", ctx, cursor).Return(
			kubernetes.NodeList{
				Items: []kubernetes.Node{
					{
						ID:   1,
						Name: "worker1",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
					{
						ID:   1,
						Name: "worker2",
						Status: kubernetes.NodeStatus{
							ID:   1,
							Key:  "healthy",
							Name: "Healthy",
						},
					},
				},
			},
			nil,
		).Once()

		manager.client = client
		err = manager.Refresh()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(manager.nodeGroups), "number of node groups do not match")
		assert.Equal(t, 1, manager.nodeGroups[0].minSize, "minimum node for node group does not match")
		assert.Equal(t, 10, manager.nodeGroups[0].maxSize, "maximum node for node group does not match")
	})
}

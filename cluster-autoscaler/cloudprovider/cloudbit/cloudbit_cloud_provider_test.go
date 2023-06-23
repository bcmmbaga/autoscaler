package cloudbit

import (
	"context"
	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"
	"github.com/stretchr/testify/mock"
)

type cloudbitClientMock struct {
	mock.Mock
}

func (c *cloudbitClientMock) List(ctx context.Context, cursor goclient.Cursor) (kubernetes.NodeList, error) {
	args := c.Called(ctx, cursor)
	return args.Get(0).(kubernetes.NodeList), args.Error(1)
}

func (c *cloudbitClientMock) PerformAction(ctx context.Context, nodeID int, req kubernetes.NodePerformAction) (kubernetes.Node, error) {
	args := c.Called(ctx, nodeID, req)
	return args.Get(0).(kubernetes.Node), args.Error(1)
}

func (c *cloudbitClientMock) Delete(ctx context.Context, nodeID int) error {
	args := c.Called(ctx, nodeID)
	return args.Error(0)
}

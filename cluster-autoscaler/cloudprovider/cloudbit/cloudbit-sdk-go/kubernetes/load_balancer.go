package kubernetes

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/compute"
)

type LoadBalancer = compute.LoadBalancer
type LoadBalancerList = compute.LoadBalancerList

type LoadBalancerService struct {
	client    cloudbitgo.Client
	clusterID int
}

func NewLoadBalancerService(client cloudbitgo.Client, clusterID int) LoadBalancerService {
	return LoadBalancerService{
		client:    client,
		clusterID: clusterID,
	}
}

func (v LoadBalancerService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list LoadBalancerList, err error) {
	list.Pagination, err = v.client.List(ctx, getLoadBalancerPath(v.clusterID), cursor, &list.Items)
	return
}

const loadBalancerSegment = "load-balancers"

func getLoadBalancerPath(clusterID int) string {
	return cloudbitgo.Join(clusterSegment, clusterID, loadBalancerSegment)
}

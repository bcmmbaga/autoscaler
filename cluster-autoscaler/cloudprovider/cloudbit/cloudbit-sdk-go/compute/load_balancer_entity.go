package compute

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

const (
	LoadBalancerStatusActive = iota + 1
	LoadBalancerStatusDisabled
	LoadBalancerStatusWorking
	LoadBalancerStatusDegraded
	LoadBalancerStatusError
)

type LoadBalancerAlgorithm struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerAlgorithmList struct {
	Items      []LoadBalancerAlgorithm
	Pagination cloudbitgo.Pagination
}

type LoadBalancerProtocol struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerProtocolList struct {
	Items      []LoadBalancerProtocol
	Pagination cloudbitgo.Pagination
}

type LoadBalancerHealthCheckType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerHealthCheckTypeList struct {
	Items      []LoadBalancerHealthCheckType
	Pagination cloudbitgo.Pagination
}

type LoadBalancerStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type LoadBalancerEntityService struct {
	client cloudbitgo.Client
}

func NewLoadBalancerEntityService(client cloudbitgo.Client) LoadBalancerEntityService {
	return LoadBalancerEntityService{client: client}
}

func (l LoadBalancerEntityService) ListAlgorithms(ctx context.Context, cursor cloudbitgo.Cursor) (list LoadBalancerAlgorithmList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-algorithms", cursor, &list.Items)
	return
}

func (l LoadBalancerEntityService) ListProtocols(ctx context.Context, cursor cloudbitgo.Cursor) (list LoadBalancerProtocolList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-protocols", cursor, &list.Items)
	return
}

func (l LoadBalancerEntityService) ListHealthCheckTypes(ctx context.Context, cursor cloudbitgo.Cursor) (list LoadBalancerHealthCheckTypeList, err error) {
	list.Pagination, err = l.client.List(ctx, "/v4/entities/compute/load-balancer-health-check-types", cursor, &list.Items)
	return
}

package objectstorage

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/common"
)

type Instance struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Location common.Location `json:"location"`
}

type InstanceList struct {
	Items      []Instance
	Pagination cloudbitgo.Pagination
}

type InstanceCreate struct {
	LocationID int `json:"location_id"`
}

type InstanceService struct {
	client cloudbitgo.Client
}

func NewInstanceService(client cloudbitgo.Client) InstanceService {
	return InstanceService{
		client: client,
	}
}

func (i InstanceService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list InstanceList, err error) {
	list.Pagination, err = i.client.List(ctx, getInstancePath(), cursor, &list.Items)
	return
}

func (i InstanceService) Create(ctx context.Context, body InstanceCreate) (instance Instance, err error) {
	err = i.client.Create(ctx, getInstancePath(), body, &instance)
	return
}

func (i InstanceService) Delete(ctx context.Context, id int) (err error) {
	err = i.client.Delete(ctx, getSpecificInstancePath(id))
	return
}

const instanceSegment = "/v4/object-storage/instances"

func getInstancePath() string {
	return instanceSegment
}

func getSpecificInstancePath(loadBalancerID int) string {
	return cloudbitgo.Join(instanceSegment, loadBalancerID)
}

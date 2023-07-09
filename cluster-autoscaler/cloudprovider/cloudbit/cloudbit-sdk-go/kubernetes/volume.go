package kubernetes

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/compute"
)

type Volume = compute.Volume
type VolumeList = compute.VolumeList

type VolumeService struct {
	client    cloudbitgo.Client
	clusterID int
}

func NewVolumeService(client cloudbitgo.Client, clusterID int) VolumeService {
	return VolumeService{
		client:    client,
		clusterID: clusterID,
	}
}

func (v VolumeService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list VolumeList, err error) {
	list.Pagination, err = v.client.List(ctx, getVolumePath(v.clusterID), cursor, &list.Items)
	return
}

func (v VolumeService) Delete(ctx context.Context, id int) (err error) {
	err = v.client.Delete(ctx, getSpecificVolumePath(v.clusterID, id))
	return
}

const volumeSegment = "volumes"

func getVolumePath(clusterID int) string {
	return cloudbitgo.Join(clusterSegment, clusterID, volumeSegment)
}

func getSpecificVolumePath(clusterID, volumeID int) string {
	return cloudbitgo.Join(clusterSegment, clusterID, volumeSegment, volumeID)
}

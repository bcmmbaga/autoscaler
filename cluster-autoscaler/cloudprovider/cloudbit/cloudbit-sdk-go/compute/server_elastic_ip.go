package compute

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

type ElasticIPAttach struct {
	ElasticIPID        int `json:"elastic_ip_id"`
	NetworkInterfaceID int `json:"network_interface_id"`
}

type ServerElasticIPService struct {
	client   cloudbitgo.Client
	serverID int
}

func NewServerElasticIPService(client cloudbitgo.Client, serverID int) ServerElasticIPService {
	return ServerElasticIPService{
		client:   client,
		serverID: serverID,
	}
}

func (s ServerElasticIPService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list ElasticIPList, err error) {
	list.Pagination, err = s.client.List(ctx, getServerElasticIPsPath(s.serverID), cursor, &list.Items)
	return
}

func (s ServerElasticIPService) Attach(ctx context.Context, body ElasticIPAttach) (elasticIP ElasticIP, err error) {
	err = s.client.Create(ctx, getServerElasticIPsPath(s.serverID), body, &elasticIP)
	return
}

func (s ServerElasticIPService) Detach(ctx context.Context, id int) (err error) {
	err = s.client.Delete(ctx, getSpecificServerElasticIPPath(s.serverID, id))
	return
}

const serverElasticIPsSegment = "elastic-ips"

func getServerElasticIPsPath(serverID int) string {
	return cloudbitgo.Join(serversSegment, serverID, serverElasticIPsSegment)
}

func getSpecificServerElasticIPPath(serverID, elasticIPID int) string {
	return cloudbitgo.Join(serversSegment, serverID, serverElasticIPsSegment, elasticIPID)
}

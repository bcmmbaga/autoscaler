package macbaremetal

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/common"
)

type Network struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Subnet              string          `json:"cidr"`
	Location            common.Location `json:"location"`
	DomainName          string          `json:"domain_name"`
	DomainNameServers   []string        `json:"domain_name_servers"`
	AllocationPoolStart string          `json:"allocation_pool_start"`
	AllocationPoolEnd   string          `json:"allocation_pool_end"`
	GatewayIP           string          `json:"gateway_ip"`
	UsedIPs             int             `json:"used_ips"`
	TotalIPs            int             `json:"total_ips"`
}

type NetworkList struct {
	Items      []Network
	Pagination cloudbitgo.Pagination
}

type NetworkCreate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	LocationID  int    `json:"location_id,omitempty"`
}

type NetworkUpdate struct {
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	DomainName        string   `json:"domain_name,omitempty"`
	DomainNameServers []string `json:"domain_name_servers,omitempty"`
}

type NetworkService struct {
	client cloudbitgo.Client
}

func NewNetworkService(client cloudbitgo.Client) NetworkService {
	return NetworkService{client: client}
}

func (n NetworkService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list NetworkList, err error) {
	list.Pagination, err = n.client.List(ctx, getNetworksPath(), cursor, &list.Items)
	return
}

func (n NetworkService) Get(ctx context.Context, id int) (network Network, err error) {
	err = n.client.Get(ctx, getSpecificNetworkPath(id), &network)
	return
}

func (n NetworkService) Create(ctx context.Context, body NetworkCreate) (network Network, err error) {
	err = n.client.Create(ctx, getNetworksPath(), body, &network)
	return
}

func (n NetworkService) Update(ctx context.Context, id int, body NetworkUpdate) (network Network, err error) {
	err = n.client.Update(ctx, getSpecificNetworkPath(id), body, &network)
	return
}

func (n NetworkService) Delete(ctx context.Context, id int) (err error) {
	err = n.client.Delete(ctx, getSpecificNetworkPath(id))
	return
}

const networksSegment = "/v4/macbaremetal/networks"

func getNetworksPath() string {
	return networksSegment
}

func getSpecificNetworkPath(id int) string {
	return cloudbitgo.Join(networksSegment, id)
}

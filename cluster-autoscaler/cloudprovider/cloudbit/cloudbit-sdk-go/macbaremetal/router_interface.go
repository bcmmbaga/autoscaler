package macbaremetal

import (
	"context"
)

type RouterInterface struct {
	ID        int     `json:"id"`
	PrivateIP string  `json:"private_ip"`
	Network   Network `json:"network"`
}

type RouterInterfaceList struct {
	Items      []RouterInterface
	Pagination cloudbitgo.Pagination
}

type RouterInterfaceService struct {
	client   cloudbitgo.Client
	routerID int
}

func NewRouterInterfaceService(client cloudbitgo.Client, routerID int) RouterInterfaceService {
	return RouterInterfaceService{client: client, routerID: routerID}
}

func (r RouterInterfaceService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list RouterInterfaceList, err error) {
	list.Pagination, err = r.client.List(ctx, getRouterInterfacesPath(r.routerID), cursor, &list.Items)
	return
}

const routerInterfacesSegment = "router-interfaces"

func getRouterInterfacesPath(id int) string {
	return cloudbitgo.Join(routersSegment, id, routerInterfacesSegment)
}

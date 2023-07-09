package common

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

type Module struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Parent    *Module    `json:"parent"`
	Sorting   int        `json:"sorting"`
	Locations []Location `json:"locations"`
}

type ModuleList struct {
	cloudbitgo.Pagination
	Items []Module
}

type ModuleService struct {
	client cloudbitgo.Client
}

func NewModuleService(client cloudbitgo.Client) ModuleService {
	return ModuleService{client: client}
}

func (l ModuleService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list ModuleList, err error) {
	list.Pagination, err = l.client.List(ctx, getModulesPath(), cursor, &list.Items)
	return
}

func (l ModuleService) Get(ctx context.Context, id int) (module Module, err error) {
	err = l.client.Get(ctx, getSpecificModulePath(id), &module)
	return
}

const modulesSegment = "/v4/entities/modules"

func getModulesPath() string {
	return modulesSegment
}

func getSpecificModulePath(moduleID int) string {
	return cloudbitgo.Join(modulesSegment, moduleID)
}

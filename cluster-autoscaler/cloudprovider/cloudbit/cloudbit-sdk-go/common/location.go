package common

import (
	"context"

	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

type Location struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Key     string   `json:"key"`
	City    string   `json:"city"`
	Modules []Module `json:"available_modules"`
}

type LocationList struct {
	cloudbitgo.Pagination
	Items []Location
}

type LocationService struct {
	client cloudbitgo.Client
}

func NewLocationService(client cloudbitgo.Client) LocationService {
	return LocationService{client: client}
}

func (l LocationService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list LocationList, err error) {
	list.Pagination, err = l.client.List(ctx, getLocationsPath(), cursor, &list.Items)
	return
}

func (l LocationService) Get(ctx context.Context, id int) (location Location, err error) {
	err = l.client.Get(ctx, getSpecificLocationPath(id), &location)
	return
}

const locationsSegment = "/v4/entities/locations"

func getLocationsPath() string {
	return locationsSegment
}

func getSpecificLocationPath(locationID int) string {
	return cloudbitgo.Join(locationsSegment, locationID)
}

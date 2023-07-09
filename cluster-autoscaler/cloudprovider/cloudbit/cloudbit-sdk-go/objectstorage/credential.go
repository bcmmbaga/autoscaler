package objectstorage

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go/common"
)

type Credential struct {
	ID        int             `json:"id"`
	Location  common.Location `json:"location"`
	Endpoint  string          `json:"endpoint"`
	AccessKey string          `json:"access_key"`
	SecretKey string          `json:"secret_key"`
}

type CredentialList struct {
	Items      []Credential
	Pagination cloudbitgo.Pagination
}

type CredentialService struct {
	client cloudbitgo.Client
}

func NewCredentialService(client cloudbitgo.Client) CredentialService {
	return CredentialService{
		client: client,
	}
}

func (i CredentialService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list CredentialList, err error) {
	list.Pagination, err = i.client.List(ctx, getCredentialSegment(), cursor, &list.Items)
	return
}

const credentialSegment = "/v4/object-storage/credentials"

func getCredentialSegment() string {
	return credentialSegment
}

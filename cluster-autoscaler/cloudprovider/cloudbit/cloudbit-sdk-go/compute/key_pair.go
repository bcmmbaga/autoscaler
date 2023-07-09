package compute

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

type KeyPair struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type KeyPairList struct {
	Items      []KeyPair
	Pagination cloudbitgo.Pagination
}

type KeyPairCreate struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type KeyPairService struct {
	client cloudbitgo.Client
}

func NewKeyPairService(client cloudbitgo.Client) KeyPairService {
	return KeyPairService{client: client}
}

func (k KeyPairService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list KeyPairList, err error) {
	list.Pagination, err = k.client.List(ctx, getKeyPairsPath(), cursor, &list.Items)
	return
}

func (k KeyPairService) Create(ctx context.Context, body KeyPairCreate) (keyPair KeyPair, err error) {
	err = k.client.Create(ctx, getKeyPairsPath(), body, &keyPair)
	return
}

func (k KeyPairService) Delete(ctx context.Context, id int) (err error) {
	err = k.client.Delete(ctx, getSpecificKeyPairPath(id))
	return
}

const keyPairsSegment = "/v4/compute/key-pairs"

func getKeyPairsPath() string {
	return keyPairsSegment
}

func getSpecificKeyPairPath(id int) string {
	return cloudbitgo.Join(keyPairsSegment, id)
}

package compute

import (
	"context"
	cloudbitgo "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/cloudbit/cloudbit-sdk-go"
)

type LoadBalancerMember struct {
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	Address string             `json:"address"`
	Port    int                `json:"port"`
	Status  LoadBalancerStatus `json:"status"`
}

type LoadBalancerMemberList struct {
	Items      []LoadBalancerMember
	Pagination cloudbitgo.Pagination
}

type LoadBalancerMemberCreate struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type LoadBalancerMemberService struct {
	client cloudbitgo.Client

	loadBalancerID int
	poolID         int
}

func NewLoadBalancerMemberService(client cloudbitgo.Client, loadBalancerID, poolID int) LoadBalancerMemberService {
	return LoadBalancerMemberService{client: client, loadBalancerID: loadBalancerID, poolID: poolID}
}

func (m LoadBalancerMemberService) List(ctx context.Context, cursor cloudbitgo.Cursor) (list LoadBalancerMemberList, err error) {
	list.Pagination, err = m.client.List(ctx, getLoadBalancerMembersPath(m.loadBalancerID, m.poolID), cursor, &list.Items)
	return
}

func (m LoadBalancerMemberService) Create(ctx context.Context, body LoadBalancerMemberCreate) (member LoadBalancerMember, err error) {
	err = m.client.Create(ctx, getLoadBalancerMembersPath(m.loadBalancerID, m.poolID), body, &member)
	return
}

func (m LoadBalancerMemberService) Delete(ctx context.Context, id int) (err error) {
	err = m.client.Delete(ctx, getSpecificLoadBalancerMemberPath(m.loadBalancerID, m.poolID, id))
	return
}

const loadBalancerMembersSegment = "members"

func getLoadBalancerMembersPath(loadBalancerID, poolID int) string {
	return cloudbitgo.Join(loadBalancerSegment, loadBalancerID, loadBalancerPoolsSegment, poolID, loadBalancerMembersSegment)
}

func getSpecificLoadBalancerMemberPath(loadBalancerID, poolID, memberID int) string {
	return cloudbitgo.Join(loadBalancerSegment, loadBalancerID, loadBalancerPoolsSegment, poolID, loadBalancerMembersSegment, memberID)
}

package cloudbit

import (
	"github.com/flowswiss/goclient/kubernetes"
)

type NodeGroup struct {
	id        int
	clusterID int
	client    nodeGroupClient
	nodePool  *kubernetes.Node

	minSize int
	maxSize int
}

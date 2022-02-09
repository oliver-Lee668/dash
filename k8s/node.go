package k8s

import (
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/core/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeClient struct {
	clientset  *kubernetes.Clientset
	nodeLister v1.NodeLister
}

func NewNodeClient(clientset *kubernetes.Clientset, nodeLister v1.NodeLister) *NodeClient {
	return &NodeClient{
		clientset:  clientset,
		nodeLister: nodeLister,
	}
}

func (cli *NodeClient) List(label string) ([]*corev1.Node, error) {
	return cli.nodeLister.List(labels.Everything())
}

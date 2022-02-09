package k8s

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type DynamicClient struct {
	client dynamic.Interface
}

func NewDynamicClient(client dynamic.Interface) DynamicClient {
	return DynamicClient{
		client: client,
	}
}

// DynamicList 列表
func (c *DynamicClient) DynamicList(namespace, group, version, resource string) ([]unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	s, err := c.client.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return s.Items, nil
}

// GetDynamic 获取dynamic
func (c *DynamicClient) GetDynamic(namespace, name, group, version, resource string) (*unstructured.Unstructured,
	error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	s, err := c.client.Resource(gvr).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// AddDynamic 添加dynamic
func (c *DynamicClient) AddDynamic(namespace, group, version, resource string,
	obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	s, err := c.client.Resource(gvr).Namespace(namespace).Create(context.Background(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

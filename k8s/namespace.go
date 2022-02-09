package k8s

import (
	"context"
	v1 "k8s.io/client-go/listers/core/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceClient struct {
	clientset *kubernetes.Clientset
	nsLister  v1.NamespaceLister
}

func NewNamespaceClient(clientset *kubernetes.Clientset, nsLister v1.NamespaceLister) *NamespaceClient {
	return &NamespaceClient{
		clientset: clientset,
		nsLister:  nsLister,
	}
}

func (cli *NamespaceClient) AddNamespace(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	ns, err := cli.clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (cli *NamespaceClient) List(labels string) ([]corev1.Namespace, error) {
	opts := metav1.ListOptions{}
	if labels != "" {
		opts.LabelSelector = labels
	}
	namespaceList, err := cli.clientset.CoreV1().Namespaces().List(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return namespaceList.Items, nil
}

func (cli *NamespaceClient) DeleteNamespace(namespace string) error {
	return cli.clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
}

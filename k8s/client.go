package k8s

import (
	"flag"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// kubeclient 定义包含所有需要操作的client
type kubeclient struct {
	Pod       *PodClient
	Node      *NodeClient
	Namespace *NamespaceClient
	Dynamic   DynamicClient
}

var Client *kubeclient

func initK8sClient() (*rest.Config, error) {
	var err error
	var config *rest.Config
	// inCluster（Pod）、Kubeconfig（kubectl）
	// 通过flag传递kubeconfig参数
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选) kubeconfig 文件的绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的绝对路径")
	}
	flag.Parse()
	// 首先使用 inCluster 模式（RBAC -> list|get node）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 kubeconfig 模式
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func NewKubeClient() error { //*kubernetes.Clientset,
	config, err := initK8sClient()
	if err != nil {
		return err
	}
	// 创建clientSet对象
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	podLister := informerFactory.Core().V1().Pods().Lister()
	nodeLister := informerFactory.Core().V1().Nodes().Lister()
	nsLister := informerFactory.Core().V1().Namespaces().Lister()

	stopper := make(chan struct{})
	informerFactory.Start(stopper)

	// 创建dynamicClient对象
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	Client = &kubeclient{
		Pod:       NewPodClient(clientSet, config, podLister),
		Node:      NewNodeClient(clientSet, nodeLister),
		Namespace: NewNamespaceClient(clientSet, nsLister),
		Dynamic:   NewDynamicClient(dynamicClient),
	}
	return nil
}

//func onAdd(obj interface{}) {
//	deploy := obj.(*v1.Deployment)
//	fmt.Println("add a deployment:", deploy.Name)
//}
//
//func onUpdate(old, new interface{}) {
//	oldDeploy := old.(*v1.Deployment)
//	newDeploy := new.(*v1.Deployment)
//	fmt.Println("update deployment:", oldDeploy.Name, newDeploy.Name)
//}
//
//func onDelete(obj interface{}) {
//	deploy := obj.(*v1.Deployment)
//	fmt.Println("delete a deployment:", deploy.Name)
//}

package controllers

import (
	"github.com/cnych/dash/k8s"
	"github.com/cnych/dash/models"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// AddNamespace 添加namespace
func AddNamespace(c *gin.Context) {
	input := models.NamespaceDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.NamespaceDto to json failed", "controller", "AddNamespace")
		writeOK(c, gin.H{})
		return
	}

	ns := corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              input.Name,
			CreationTimestamp: metav1.Time{},
			Labels:            map[string]string{"kubernetes.io/metadata.name": input.Name},
		},
		Spec: corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{
			Phase: "Active",
		},
	}
	res, err := k8s.Client.Namespace.AddNamespace(&ns)
	if err != nil {
		klog.V(2).ErrorS(err, "add namespace in cluster failed", "controller", "AddNamespace")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"namespace": res})
}

// GetNamespaceList 获取namespace列表
func GetNamespaceList(c *gin.Context) {
	namespaces, err := k8s.Client.Namespace.List("")
	if err != nil {
		klog.V(2).ErrorS(err, "get namespace list failed", "controller", "GetNamespaceList")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"namespaces": namespaces})
}

// DeleteNamespace 删除namespace
func DeleteNamespace(c *gin.Context) {
	namespace := c.Param("name")
	err := k8s.Client.Namespace.DeleteNamespace(namespace)
	if err != nil {
		klog.V(2).ErrorS(err, "delete namespace failed", "controller", "DeleteNamespace")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"namespace": namespace})
}

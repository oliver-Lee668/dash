package controllers

import (
	"fmt"
	"github.com/cnych/dash/k8s"
	"github.com/cnych/dash/models"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"
)

// GetDynamicList 列表
func GetDynamicList(c *gin.Context) {
	input := models.DynamicListDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.DynamicListDto to json failed", "controller",
			"GetDynamicList")
		writeOK(c, gin.H{})
		return
	}
	s, err := k8s.Client.Dynamic.DynamicList(input.Namespace, input.Group, input.Version, input.Resource)
	if err != nil {
		klog.V(2).ErrorS(err, "get dynamic list failed", "controller", "GetDynamicList")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"dynamics": s})
}

// GetDynamic 获取dynamic
func GetDynamic(c *gin.Context) {
	input := models.DynamicListDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.DynamicListDto to json failed", "controller",
			"GetDynamicList")
		writeOK(c, gin.H{})
		return
	}
	s, err := k8s.Client.Dynamic.GetDynamic(input.Namespace, input.Name, input.Group, input.Version, input.Resource)
	if err != nil {
		klog.V(2).ErrorS(err, "get dynamic failed", "controller", "GetDynamic")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"dynamic": s})
}

// AddDynamic 添加dynamic
func AddDynamic(c *gin.Context) {
	input := models.DynamicListDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.DynamicListDto to json failed", "controller",
			"GetDynamicList")
		writeOK(c, gin.H{})
		return
	}
	obj := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", input.Group, input.Version),
			"kind":       input.Kind,
			"metadata": map[string]interface{}{
				"name": input.Name,
			},
			"spec": map[string]interface{}{
				"replicas": 2,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "web",
								"image": "nginx:1.12",
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 80,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	s, err := k8s.Client.Dynamic.AddDynamic(input.Namespace, input.Group, input.Version, input.Resource, &obj)
	if err != nil {
		klog.V(2).ErrorS(err, "add dynamic failed", "controller", "AddDynamic")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"dynamic": s})
}

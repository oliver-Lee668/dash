package controllers

import (
	"fmt"
	"github.com/cnych/dash/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"

	"github.com/cnych/dash/k8s"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// AddPod 添加pod
func AddPod(c *gin.Context) {
	input := models.PodDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.PodDto to json failed", "controller", "AddPod")
		writeOK(c, gin.H{})
		return
	}

	namespace := "test"
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: namespace,
			Labels:    input.Label,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  input.ContainerName,
					Image: input.ContainerImage,
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
	res, err := k8s.Client.Pod.Add(namespace, &pod)
	if err != nil {
		klog.V(2).ErrorS(err, "add pod in cluster failed", "controller", "AddPod")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"pods": res})
}

// GetKubeLogs 实时获取Pod的日志
func GetKubeLogs(c *gin.Context) {
	// /api/v1/namespaces/kube-system/pods/traefik-7cb4cb6bf5-p779x/logs?tailLines=500&timestamps=true&previous=false&container=traefik'
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.Query("container")
	tailLines, _ := strconv.ParseInt(c.DefaultQuery("tailLines", "500"), 10, 64)
	timestamps, _ := strconv.ParseBool(c.DefaultQuery("timestamps", "true"))
	previous, _ := strconv.ParseBool(c.DefaultQuery("previous", "false"))

	klog.V(2).InfoS("get kube logs request params", "namespace", namespace, "pod", podName,
		"container", container, "tailLines", tailLines, "timestamps", timestamps, "previous", previous)

	if namespace == "" || podName == "" || container == "" {
		c.String(http.StatusBadRequest, "must specific namespace、pod and container query params")
		return
	}

	// 获取pod的日志（websocket）
	// 把当前的get http request -> upgrade websocket
	kubeLogger, err := k8s.NewKubeLogger(c.Writer, c.Request, nil)
	if err != nil {
		klog.Error(err, "upgrade websocket failed")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 构造获取日志的结构体
	opts := corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		TailLines:  &tailLines,
		Timestamps: timestamps,
		Previous:   previous,
	}
	if err := k8s.Client.Pod.LogsStream(podName, namespace, &opts, kubeLogger); err != nil {
		klog.Error(err, "GetLogs stream failed")
		_, _ = kubeLogger.Write([]byte(err.Error()))
	}
}

func HandleTerminal(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.Query("container")
	cmd := []string{
		"/bin/sh", "-c", "clear;(bash || sh)",
	}
	klog.V(2).InfoS("get kube logs request params",
		"namespace", namespace, "pod", podName, "container", container, "cmd", cmd)

	if _, err := k8s.Client.Pod.Get(podName, namespace); err != nil {
		klog.Error(err, "get pod failed")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// todo，校验下 pod
	kubeShell, err := k8s.NewKubeShell(c.Writer, c.Request, nil)
	if err != nil {
		klog.Error(err, "init kube shell failed")
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		_ = kubeShell.Close()
	}()

	if err := k8s.Client.Pod.Exec(cmd, kubeShell, namespace, podName, container); err != nil {
		klog.Error(err, "exec pod failed")
		c.String(http.StatusBadRequest, err.Error())
	}
}

func DeletePod(c *gin.Context) {
	input := models.PodDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		klog.V(2).ErrorS(err, "bind models.PodDto to json failed", "controller", "DeletePod")
		writeOK(c, gin.H{})
		return
	}

	if err := k8s.Client.Pod.Delete(input.Name, input.Namespace); err != nil {
		klog.V(2).ErrorS(err, "delete pod failed", "controller", "DeletePod")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"pods": input.Name})
}

// GetPodList 获取集群所有Pod列表
func GetPodList(c *gin.Context) {
	pods, err := k8s.Client.Pod.List("")
	if err != nil {
		klog.V(2).ErrorS(err, "get pod list in cluster failed", "controller", "GetPodList")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"pods": pods})
}

// GetKubePods 获取指定namespace下的Pod列表
func GetKubePods(c *gin.Context) {
	namespace := c.Param("namespace")
	pods, err := k8s.Client.Pod.List(namespace)
	if err != nil {
		klog.V(2).ErrorS(err, fmt.Sprintf("get pod list in %s failed", namespace), "controller", "GetKubePods")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"pods": pods})
}

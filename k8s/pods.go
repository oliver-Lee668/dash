package k8s

import (
	"bufio"
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/listers/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type TtyHandler interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Tty() bool
	remotecommand.TerminalSizeQueue
	Done()
}

type PodClient struct {
	clientset *kubernetes.Clientset
	config    *restclient.Config
	podLister v1.PodLister
}

func NewPodClient(clientset *kubernetes.Clientset, config *restclient.Config, podLister v1.PodLister) *PodClient {
	return &PodClient{
		clientset: clientset,
		config:    config,
		podLister: podLister,
	}
}

func (cli *PodClient) Add(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	res, err := cli.clientset.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cli *PodClient) Delete(name, namespace string) error {
	return cli.clientset.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (cli *PodClient) Get(name, namespace string) (*corev1.Pod, error) {
	return cli.podLister.Pods(namespace).Get(name)
}

func (cli *PodClient) List(namespace string) ([]*corev1.Pod, error) {
	return cli.podLister.Pods(namespace).List(labels.Everything())
}

func (cli *PodClient) Logs(name, namespace string, opts *corev1.PodLogOptions) *restclient.Request {
	return cli.clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
}

func (cli *PodClient) LogsStream(name, namespace string, opts *corev1.PodLogOptions, writer io.Writer) error {
	req := cli.Logs(name, namespace, opts)
	stream, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer stream.Close()

	buf := bufio.NewReader(stream)
	for { // 一直从buffer中读取数据去
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_, err = writer.Write(bytes)
			}
			return err
		}
		_, err = writer.Write(bytes)
		if err != nil {
			return err
		}
	}

}

func (cli *PodClient) Exec(cmd []string, handler TtyHandler, namespace, pod, container string) error {
	defer func() {
		handler.Done()
	}()

	// 构造请求
	req := cli.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(pod).SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: container,
		Command:   cmd,
		Stdin:     handler.Stdin() != nil,
		Stdout:    handler.Stdout() != nil,
		Stderr:    handler.Stderr() != nil,
		TTY:       handler.Tty(),
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(cli.config, "POST", req.URL())
	if err != nil {
		return err
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler.Stdin(),
		Stdout:            handler.Stdout(),
		Stderr:            handler.Stderr(),
		Tty:               handler.Tty(),
		TerminalSizeQueue: handler,
	})
}

package exec

import (
	"bytes"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// 入参为kubeconfig、pod名字、命名空间、命令字符串
func cmdExecuter(config *restclient.Config, podName, namespace, cmd string) (map[string]string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// 构造执行命令请求
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{cmd},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)
		// 执行命令
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	// 使用bytes.Buffer变量接收标准输出和标准错误
	var stdout, stderr bytes.Buffer
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  strings.NewReader(""),
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		return nil, err
	}
	// 返回数据
	ret := map[string]string{"stdout": stdout.String(), "stderr": stderr.String(), "pod_name": podName}
	return ret, nil
}

package get

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	version "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
)

type GetResourceWorker struct {
	Client *kubernetes.Clientset
}

type K8sInfo struct {
	GitVersion string
	GoVersion  string
	Platform   string
}

// get node info by name
type NodeInfo struct {
	K8sInfo

	address string
	CPU     int64
	Memory  int64
}

// namespace

// namespace list
func (grw *GetResourceWorker) GetNameSpaceList(ctx context.Context) (nameSpaceList *corev1.NamespaceList, err error) {
	if nameSpaceList, err = grw.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return nameSpaceList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get namespace list")
}

// is namespace exist?
func (grw *GetResourceWorker) IsNameSpaceExist(namespace string) (b bool) {
	var (
		list *corev1.NamespaceList
		ctx  = context.Background()
		err  error
	)
	if list, err = grw.GetNameSpaceList(ctx); err != nil {
		return false
	}
	for _, v := range list.Items {
		if v.ObjectMeta.Name == namespace {
			return true
		}
	}
	return false
}

// deployment
func (grw *GetResourceWorker) GetdeploymentList(ctx context.Context) (deploymentList *appsv1.DeploymentList, err error) {
	if deploymentList, err = grw.Client.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return deploymentList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get deployment list")
}

// pod
func (grw *GetResourceWorker) GetPodList(ctx context.Context) (podList *corev1.PodList, err error) {
	if podList, err = grw.Client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return podList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get deployment list")

}

// node
func (grw *GetResourceWorker) GetNodeList(ctx context.Context) (nodeList *corev1.NodeList, err error) {
	if nodeList, err = grw.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return nodeList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get node list")
}

// GetK8sIpAddress
// cluster ip
func (grw *GetResourceWorker) GetK8sIpAddress(ctx context.Context) (addresses []string, err error) {
	var (
		nodeList *corev1.NodeList
	)
	if nodeList, err = grw.GetNodeList(ctx); err != nil {
		goto FAIL
	}
	for _, v := range nodeList.Items {
		addresses = append(addresses, v.Status.Addresses[0].Address)
	}
	return
FAIL:
	return nil, errors.Wrap(err, "fail to get k8s ip address")
}

func (grw *GetResourceWorker) GetNodeInfoByName(ctx context.Context, name string) (nodeinfo NodeInfo, err error) {
	var (
		nl      *corev1.NodeList
		node    corev1.Node
		k8sinfo *version.Info
	)
	set := make(map[string]interface{})
	if nl, err = grw.GetNodeList(ctx); err != nil {
		goto FAIL
	}
	// found node by name
	for _, v := range nl.Items {
		set[v.ObjectMeta.Name] = v
	}
	if v, ok := set[name]; ok {
		node = v.(corev1.Node)
	} else {
		err = errors.Errorf("node [%s] not found", name)
		goto FAIL
	}
	// id node ready?
	for _, v := range node.Status.Conditions {
		if v.Type == corev1.NodeReady {
			if v.Status == corev1.ConditionFalse {
				err = errors.Errorf("node [%s] is unhealthy", name)
				goto FAIL
			}
		}
	}
	// is address type "NodeInternalIP"
	for _, v := range node.Status.Addresses {
		set[string(v.Type)] = struct{}{}
	}
	if _, ok := set[string(corev1.NodeInternalIP)]; !ok {
		err = errors.Errorf("node [%s] type is not InternalIP", name)
		goto FAIL
	}
	// get k8s info
	if k8sinfo, err = grw.GetK8sInfo(ctx); err != nil {
		goto FAIL
	}
	// get node info
	nodeinfo.address = node.Status.Addresses[0].Address
	nodeinfo.CPU = node.Status.Capacity.Cpu().Value()
	nodeinfo.Memory = node.Status.Capacity.Memory().Value()
	nodeinfo.K8sInfo.GitVersion = k8sinfo.GitVersion
	nodeinfo.K8sInfo.GoVersion = k8sinfo.GoVersion
	nodeinfo.K8sInfo.Platform = k8sinfo.Platform
	return nodeinfo, nil
FAIL:
	return nodeinfo, errors.Wrap(err, "get node info by name error")
}

func (grw *GetResourceWorker) GetK8sInfo(context.Context) (k8sinfo *version.Info, err error) {
	if k8sinfo, err = grw.Client.Discovery().ServerVersion(); err != nil {
		goto FAIL
	}
	return k8sinfo, nil
FAIL:
	return nil, errors.Wrap(err, "get k8s info error")
}

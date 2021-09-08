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

	CPU    int64
	Memory int64
}

func (grw *GetResourceWorker) GetdeploymentList(ctx context.Context) (deploymentList *appsv1.DeploymentList, err error) {
	if deploymentList, err = grw.Client.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return deploymentList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get deployment list")
}

func (grw *GetResourceWorker) GetPodList(ctx context.Context) (podList *corev1.PodList, err error) {
	if podList, err = grw.Client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return podList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get deployment list")

}

func (grw *GetResourceWorker) GetNodeList(ctx context.Context) (nodeList *corev1.NodeList, err error) {
	if nodeList, err = grw.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	return nodeList, nil
FAIL:
	return nil, errors.Wrap(err, "fail to get node list")
}

func (grw *GetResourceWorker) GetNodeInfoByName(ctx context.Context, name string) (nodeinfo NodeInfo, err error) {
	var (
		nl      *corev1.NodeList
		node    corev1.Node
		k8sinfo *version.Info
	)
	if nl, err = grw.GetNodeList(ctx); err != nil {
		goto FAIL
	}
	// found node by name
	for _, v := range nl.Items {
		if v.ObjectMeta.Name == name {
			node = v
			break
		}
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
		if v.Type == corev1.NodeInternalIP {
			break
		}
	}
	// get k8s info
	if k8sinfo, err = grw.GetK8sInfo(ctx); err != nil {
		goto FAIL
	}
	// get node info
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

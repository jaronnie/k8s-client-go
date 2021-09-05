package get

import (
	"context"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type GetResourceWorker struct {
	Client *kubernetes.Clientset
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

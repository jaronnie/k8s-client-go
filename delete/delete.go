package delete

import (
	"context"
	"encoding/json"
	"k8test/constant"
	"k8test/get"
	"k8test/util"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8labels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type DeleteWorker struct {
	Client    *kubernetes.Clientset
	NameSpace string
}

// DeleteByName delete resource by name
// namespace name
// deployment name
// pod name
// service name
// ...
func (dw *DeleteWorker) DeleteByName(ctx context.Context, resourceType string, name string) {}

// DeleteByYAML delete resource by yaml
func (dw *DeleteWorker) DeleteByYAML(ctx context.Context, yaml string) error {
	var (
		b   [][]byte
		err error
	)
	if b, err = util.Yaml2Jsons([]byte(yaml)); err != nil {
		goto FAIL
	}
	if err = deleteOperator(dw.Client, dw.NameSpace, b); err != nil {
		goto FAIL
	}
FAIL:
	return errors.Wrap(err, "fail to delete resource")
}

// DeleteByLabels delete resource by labels
// labels
// for example
// app=kube-go-app
func (dw *DeleteWorker) DeleteByLabels(ctx context.Context, resourceType string, labels map[string]string) error {
	var (
		set k8labels.Selector
		options metav1.ListOptions
		deploymentList *appsv1.DeploymentList
		serviceList *corev1.ServiceList
		err error
	)
	set = k8labels.SelectorFromSet(labels)
	options = metav1.ListOptions{
		LabelSelector: set.String(),
	}
	switch resourceType {
	case constant.Deployment:
		if deploymentList, err = dw.Client.AppsV1().Deployments(dw.NameSpace).List(ctx, options); err != nil {
			goto FAIL
		}
		if len(deploymentList.Items) == 0 {
			err = errors.Errorf("not found resource by labels [%v]", labels)
			goto FAIL
		}
		for _, v := range deploymentList.Items {
			if err = dw.Client.AppsV1().Deployments(dw.NameSpace).Delete(ctx, v.Name, metav1.DeleteOptions{}); err != nil {
				goto FAIL
			}
		}
	case constant.Service:
		if serviceList, err = dw.Client.CoreV1().Services(dw.NameSpace).List(ctx, options); err != nil {
			goto FAIL
		}
		if len(serviceList.Items) == 0 {
			err = errors.Errorf("not found resource by labels [%v]", labels)
			goto FAIL
		}
		for _, v := range serviceList.Items {
			if err = dw.Client.CoreV1().Services(dw.NameSpace).Delete(ctx, v.Name, metav1.DeleteOptions{}); err != nil {
				goto FAIL
			}
		}
	}
FAIL:
	return errors.Wrap(err, "fail to delete by labels")
}

func deleteOperator(client *kubernetes.Clientset, namespace string, jsons [][]byte) (err error) {
	var (
		grw = &get.GetResourceWorker{
			Client: client,
		}
	)
	var (
		deployment appsv1.Deployment
		service    corev1.Service
	)
	var (
		ctx = context.Background()
	)
	// is namespace exist
	if !grw.IsNameSpaceExist(namespace) {
		err = errors.Errorf("not found namespace [%s]", namespace)
		goto FAIL
	}
	for _, v := range jsons {
		var entity map[string]interface{}
		if err := json.Unmarshal(v, &entity); err != nil {
			goto FAIL
		}
		switch entity["kind"] {
		case constant.Deployment:
			if err = json.Unmarshal(v, &deployment); err != nil {
				goto FAIL
			}
			if err = client.AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
				goto FAIL
			}
		case constant.Service:
			if err = json.Unmarshal(v, &service); err != nil {
				goto FAIL
			}
			if err = client.CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
				goto FAIL
			}
		}
	}
	return nil
FAIL:
	return errors.Wrap(err, "delete by yaml error")
}

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
	"k8s.io/client-go/kubernetes"
)

type DeleteWorker struct {
	Client    *kubernetes.Clientset
	NameSpace string
}

// delete resource by name
// namesapce name
// deployment name
// pod name
// service name
// ...
func (dw *DeleteWorker) DeleteByName(ctx context.Context, resourceType string, name string) {}

// delete resource by yaml
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

// delete resource by labels
func (dw *DeleteWorker) DeleteByLabels(ctx context.Context, labels string) {}

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

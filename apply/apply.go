package apply

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

type ApplyWorker struct {
	Client    *kubernetes.Clientset
	NameSpace string
}

// already available yaml
func (aw *ApplyWorker) ApplyByYAML(yaml string) error {
	var (
		b   [][]byte
		err error
	)
	if b, err = util.Yaml2Jsons([]byte(yaml)); err != nil {
		goto FAIL
	}
	if err = applyOperator(aw.Client, aw.NameSpace, b); err != nil {
		goto FAIL
	}
FAIL:
	return errors.Wrap(err, "fail to apply yaml")
}

func applyOperator(client *kubernetes.Clientset, namespace string, jsons [][]byte) (err error) {
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
			if _, err = client.AppsV1().Deployments(namespace).Create(ctx, &deployment, metav1.CreateOptions{}); err != nil {
				goto FAIL
			}
		case constant.Service:
			if err = json.Unmarshal(v, &service); err != nil {
				goto FAIL
			}
			if _, err = client.CoreV1().Services(namespace).Create(ctx, &service, metav1.CreateOptions{}); err != nil {
				goto FAIL
			}
		}
	}
	return nil
FAIL:
	return errors.Wrap(err, "apply yaml error")
}

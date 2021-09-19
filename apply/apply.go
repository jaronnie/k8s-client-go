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
func (aw *ApplyWorker) ApplyByYAML(yaml string) (res map[string]interface{}, err error) {
	var (
		b [][]byte
	)
	if b, err = util.Yaml2Jsons([]byte(yaml)); err != nil {
		goto FAIL
	}
	if res, err = applyOperator(aw.Client, aw.NameSpace, b); err != nil {
		goto FAIL
	}
	return
FAIL:
	return nil, errors.Wrap(err, "fail to apply yaml")
}

func applyOperator(client *kubernetes.Clientset, namespace string, jsons [][]byte) (res map[string]interface{}, err error) {
	var (
		grw = &get.GetResourceWorker{
			Client: client,
		}
	)
	var (
		deploymentReq *appsv1.Deployment
		deploymentRes *appsv1.Deployment
		serviceReq    *corev1.Service
		serviceRes    *corev1.Service
	)
	var (
		ctx = context.Background()
	)
	res = make(map[string]interface{})
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
			if err = json.Unmarshal(v, &deploymentReq); err != nil {
				goto FAIL
			}
			if deploymentRes, err = client.AppsV1().Deployments(namespace).Create(ctx, deploymentReq, metav1.CreateOptions{}); err != nil {
				goto FAIL
			}
			res[constant.Deployment] = deploymentRes
		case constant.Service:
			if err = json.Unmarshal(v, &serviceReq); err != nil {
				goto FAIL
			}
			if serviceRes, err = client.CoreV1().Services(namespace).Create(ctx, serviceReq, metav1.CreateOptions{}); err != nil {
				goto FAIL
			}
			res[constant.Service] = serviceRes
		}
	}
	return
FAIL:
	return nil, errors.Wrap(err, "apply yaml error")
}

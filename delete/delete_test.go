package delete

import (
	"context"
	"k8test/connect"
	"k8test/constant"
	"log"
	"testing"

	"k8s.io/client-go/kubernetes"
)

var (
	client *kubernetes.Clientset
	ctx    = context.Background()
)

var GoBackendApp = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-go-app-deployment
  labels:
    app: kube-go-app
spec:
  selector:
    matchLabels:
      app: kube-go-app
  replicas: 1
  template:
    metadata:
      labels:
        app: kube-go-app
    spec:
      containers:
        - name: kube-go-app
          image: "gocloudcoder/kube-go-app:v1"
          ports:
            - name: http
              containerPort: 8888
---
apiVersion: v1
kind: Service
metadata:
  name: kube-go-app-service
  labels:
    app: kube-go-app
spec:
  selector:
    app: kube-go-app
  ports:
    - name: go-app
      protocol: TCP
      port: 8888
      targetPort: 8888
  type: NodePort
`

func TestDeleteByYAML(t *testing.T) {
	var (
		dw = &DeleteWorker{
			Client:    client,
			NameSpace: "default",
		}
	)
	if err := dw.DeleteByYAML(ctx, GoBackendApp); err != nil {
		t.Fatal(err)
	}
	t.Log("success delete resouce by yaml")
}

func TestDeleteByLabels(t *testing.T) {
	var (
		dw = &DeleteWorker{
			Client:    client,
			NameSpace: "default",
		}
	)
	labels := map[string]string{
		"app": "kube-go-app",
	}
	// delete service
	if err := dw.DeleteByLabels(ctx, constant.Service, labels); err != nil {
		t.Fatal(err)
	}
	// delete deployment
	if err := dw.DeleteByLabels(ctx, constant.Deployment, labels); err != nil {
		t.Fatal(err)
	}
	t.Log("pass")
}

func init() {
	var (
		err error
	)
	// connect with kubeconfig file
	if client, err = connect.DefaultConnect(); err != nil {
		goto FAIL
	}
	// connect from url, ca, token
	// if client, err = connect.DefaultConnect(); err != nil {
	// 	goto FAIL
	// }
	return
FAIL:
	log.Fatal(err)
}

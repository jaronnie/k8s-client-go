package apply

import (
	"context"
	"k8test/connect"
	"log"
	"testing"

	"k8s.io/client-go/kubernetes"
)

var (
	client *kubernetes.Clientset
	aw     ApplyWorker
	ctx    = context.Background()
)

var goBackendApp = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend1
spec:
  selector:
    matchLabels:
      app: kube-go-app
      tier: backend1
      track: stable
  replicas: 3
  template:
    metadata:
      labels:
        app: kube-go-app
        tier: backend1
        track: stable
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
  name: kube-go-app
spec:
  selector:
    app: kube-go-app
    tier: backend1
  ports:
    - name: go-app
      protocol: TCP
      port: 8888
      targetPort: 8888
  type: LoadBalancer
`

func TestApplyByYAML(t *testing.T) {
	aw = ApplyWorker{
		Client:    client,
		NameSpace: "default",
	}
	err := aw.ApplyByYAML(goBackendApp)
	if err != nil {
		t.Log(err)
	}
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

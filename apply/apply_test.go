package apply

import (
	"encoding/json"
	"fmt"
	"k8test/connect"
	"log"
	"testing"

	"k8s.io/client-go/kubernetes"
)

var (
	client *kubernetes.Clientset
	aw     ApplyWorker
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
          resources:
            requests:
              cpu: "100m"
              memory: "100Mi"
            limits:
              cpu: "100m"
              memory: "500Mi"
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

func TestApplyByYAML(t *testing.T) {
	aw = ApplyWorker{
		Client:    client,
		NameSpace: "default",
	}
	aps, err := aw.ApplyByYAML(GoBackendApp)
	if err != nil {
		t.Log(err)
	}
	b, _ := json.Marshal(aps["Deployment"])
	fmt.Println(string(b))
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
	// if client, err = connect.Connect(); err != nil {
	// 	goto FAIL
	// }
	return
FAIL:
	log.Fatal(err)
}

package util

import (
	"fmt"
	"testing"
)

var goBackendApp = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  selector:
    matchLabels:
      app: kube-go-app
      tier: backend
      track: stable
  replicas: 3
  template:
    metadata:
      labels:
        app: kube-go-app
        tier: backend
        track: stable
    spec:
      containers:
        - name: kube-go-app
          image: gocloudcoder/kube-go-app:v1
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
    tier: backend
  ports:
    - name: go-app
      protocol: TCP
      port: 8888
      targetPort: 8888
  type: LoadBalancer
`

var goTemplateBackendApp = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  selector:
    matchLabels:
      app: kube-go-app
      tier: backend
      track: stable
  replicas: 3
  template:
    metadata:
      labels:
        app: kube-go-app
        tier: backend
        track: stable
    spec:
      containers:
        - name: kube-go-app
          image: "{{.Image}}"
          ports:
            - name: http
              containerPort: {{.Port}}
---
apiVersion: v1
kind: Service
metadata:
  name: kube-go-app
spec:
  selector:
    app: kube-go-app
    tier: backend
  ports:
    - name: go-app
      protocol: TCP
      port: {{.Port}}
      targetPort: {{.Port}}
  type: LoadBalancer
`

func TestYaml2Jsons(t *testing.T) {
	b, err := Yaml2Jsons([]byte(goBackendApp))
	if err != nil {
		t.Log(err)
	}
	for _, v := range b {
		t.Log(string(v))
	}
}

func TestParsreTemplateYAML(t *testing.T) {
	data := struct {
		Image string
		Port  int
	}{
		Image: "gocloudcoder/kube-go-app:v1",
		Port:  8888,
	}
	t.Run("test", func(t *testing.T) {
		b, err := ParseTemplateYAML(data, []byte(goTemplateBackendApp))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))
	})
}

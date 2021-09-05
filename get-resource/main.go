package main

import (
	"context"
	"fmt"
	"k8test/connect"

	appsv1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		client         *kubernetes.Clientset
		ctx            = context.Background()
		podList        *core_v1.PodList
		deploymentList *appsv1.DeploymentList
		err            error
	)
	if client, err = connect.Connect(); err != nil {
		goto FAIL
	}
	// get deployments list
	if deploymentList, err = client.AppsV1().Deployments("default").List(ctx, meta_v1.ListOptions{}); err != nil {
		goto FAIL
	}
	fmt.Printf("deployment list: %+v", deploymentList)
	// get pods list
	if podList, err = client.CoreV1().Pods("default").List(ctx, meta_v1.ListOptions{}); err != nil {
		goto FAIL
	}
	fmt.Printf("pod list: %+v", podList)
	return
FAIL:
	fmt.Println(err)
}

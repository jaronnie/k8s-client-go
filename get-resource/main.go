package main

import (
	"context"
	"fmt"
	"k8test/connect"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		client  *kubernetes.Clientset
		ctx     = context.Background()
		podList *core_v1.PodList
		err     error
	)
	if client, err = connect.Connect(); err != nil {
		goto FAIL
	}
	// get pods list
	if podList, err = client.CoreV1().Pods("docker-desktop").List(ctx, meta_v1.ListOptions{}); err != nil {
		goto FAIL
	}
	fmt.Println("pod list", podList)
	return
FAIL:
	fmt.Println(err)
}

package get

import (
	"context"
	"fmt"
	"k8test/connect"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
)

var (
	client *kubernetes.Clientset
	grw    GetResourceWorker
	ctx    = context.Background()
)

func TestGetDeploymentList(t *testing.T) {
	grw = GetResourceWorker{
		Client: client,
	}
	dl, err := grw.GetdeploymentList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dl)
}

func TestGetPodList(t *testing.T) {
	grw = GetResourceWorker{
		Client: client,
	}
	pl, err := grw.GetPodList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pl)
}

func TestGetNodeList(t *testing.T) {
	grw = GetResourceWorker{
		Client: client,
	}
	nl, err := grw.GetNodeList(ctx)
	if !assert.NotNil(t, nl) {
		t.Log(err)
		return
	}
	t.Log(nl)
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

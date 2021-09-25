package exec

import (
	"k8test/connect"
	"log"
	"testing"
)

func TestCmdExecter(t *testing.T) {
	config, err := connect.NewRestConfig()
	if err != nil {
		log.Fatal(err)
	}
	m, err := cmdExecuter(config, "kube-go-app-deployment-577978df75-glvbp", "default", "date")
	if err != nil {
		t.Log(err)
	}
	t.Log(m)
}

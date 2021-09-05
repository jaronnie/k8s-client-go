package main

import (
	"context"
	"fmt"
	"k8test/connect"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Credential struct {
	Url   string
	Ca    string
	Token string
}

func main() {
	var (
		crt = &Credential{}
		err error
	)
	if crt, err = readConfig(); err != nil {
		log.Fatal(err)
	}
	client, err := connect.GetClient("https://kubernetes.docker.internal:6443", crt.Token, crt.Ca)
	if err != nil {
		fmt.Println("get client error", err)
		return
	}
	ctx := context.Background()
	// get pods list
	podList, err := client.CoreV1().Pods("docker-desktop").List(ctx, meta_v1.ListOptions{})
	if err != nil {
		fmt.Println("get pod list error", err)
		return
	}
	fmt.Println("pod list:", podList)
}

func readConfig() (*Credential, error) {
	var (
		ca       []byte
		token    []byte
		tokenstr string
		err      error
	)
	if ca, err = os.ReadFile("../config/credential.pem"); err != nil {
		goto FAIL
	}
	if token, err = os.ReadFile("../config/token.txt"); err != nil {
		goto FAIL
	}
	tokenstr = strings.TrimSpace(string(token))
	return &Credential{
		Url:   "https://kubernetes.docker.internal:6443",
		Ca:    string(ca),
		Token: tokenstr,
	}, nil
FAIL:
	return nil, errors.Wrap(err, "fail to read config")
}

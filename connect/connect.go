package connect

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Credential struct {
	Url   string
	Ca    string
	Token string
}

func DefaultConnect() (*kubernetes.Clientset, error) {
	var (
		client     *kubernetes.Clientset
		kubeconfig []byte
		restConf   *rest.Config
		err        error
	)
	home := homedir.HomeDir()
	if kubeconfig, err = ioutil.ReadFile(home + "/.kube/config"); err != nil {
		goto FAIL
	}
	// 生成rest client配置
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		goto FAIL
	}
	if client, err = kubernetes.NewForConfig(restConf); err != nil {
		goto FAIL
	}
	return client, nil
FAIL:
	return nil, errors.Wrap(err, "connect k8s server")
}

func Connect() (*kubernetes.Clientset, error) {
	var (
		crt    *Credential
		client *kubernetes.Clientset
		err    error
	)
	if crt, err = readConfig(); err != nil {
		goto FAIL
	}
	if client, err = getClient(crt.Url, crt.Token, crt.Ca); err != nil {
		goto FAIL
	}
	return client, nil
FAIL:
	return nil, errors.Wrap(err, "connect k8s server")
}

func getClient(apiURL string, token string, caCert string) (*kubernetes.Clientset, error) {
	config, err := newConfig(apiURL, token, caCert)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return clientSet, nil

}

func newConfig(apiURL string, token string, caCert string) (*rest.Config, error) {

	tlsClientConfig := rest.TLSClientConfig{}

	if _, err := newCertPool(caCert); err != nil {
		return nil, err
	} else {
		tlsClientConfig.CAData = []byte(caCert)
	}
	return &rest.Config{
		Host:            apiURL,
		TLSClientConfig: tlsClientConfig,
		BearerToken:     token,
	}, nil
}

func newCertPool(caCert string) (*x509.CertPool, error) {
	certs, err := ParseCertsPEM([]byte(caCert))
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	for _, c := range certs {
		pool.AddCert(c)
	}
	return pool, nil
}

func ParseCertsPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false
	certs := []*x509.Certificate{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		// Only use PEM "CERTIFICATE" blocks without extra headers
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}

		certs = append(certs, cert)
		ok = true
	}

	if !ok {
		return certs, errors.New("data does not contain any valid RSA or ECDSA certificates")
	}
	return certs, nil
}

func readConfig() (*Credential, error) {
	var (
		ca       []byte
		token    []byte
		tokenstr string
		url      []byte
		urlstr   string
		err      error
	)
	if ca, err = os.ReadFile("../config/credential.pem"); err != nil {
		goto FAIL
	}
	if token, err = os.ReadFile("../config/token.txt"); err != nil {
		goto FAIL
	}
	if url, err = os.ReadFile("../config/url.txt"); err != nil {
		goto FAIL
	}
	tokenstr = strings.TrimSpace(string(token))
	urlstr = strings.TrimSpace(string(url))
	return &Credential{
		Url:   urlstr,
		Ca:    string(ca),
		Token: tokenstr,
	}, nil
FAIL:
	return nil, errors.Wrap(err, "fail to read config")
}

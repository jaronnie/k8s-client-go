package connect

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetClient(apiURL string, token string, caCert string) (*kubernetes.Clientset, error) {
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

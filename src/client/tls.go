package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TLSCerts struct {
	UseTLS              bool
	ClientCert          string
	ClientCertKey       string
	CACert              string
	SkipTLSVerification bool
}

func appendRootCert(certFilePath string, caCertPool *x509.CertPool) error {
	cert, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return err
	}
	ok := caCertPool.AppendCertsFromPEM(cert)
	if !ok {
		return fmt.Errorf("could not add %v to Root Certificates", certFilePath)
	}
	return nil
}

func setUpTLS(certs TLSCerts) error {
	if !certs.UseTLS {
		return nil
	}
	caCertPool, _ := x509.SystemCertPool()
	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}
	certificates := make([]tls.Certificate, 0)

	if certs.CACert != "" {
		err := appendRootCert(certs.CACert, caCertPool)
		if err != nil {
			return err
		}
	}

	if certs.ClientCert != "" && certs.ClientCertKey != "" {
		cert, err := tls.LoadX509KeyPair(certs.ClientCert, certs.ClientCertKey)
		if err != nil {
			return err
		}
		certificates = append(certificates, cert)
	}

	tlsConfig := &tls.Config{
		Certificates: certificates,
		RootCAs:      caCertPool,
	}
	if certs.SkipTLSVerification {
		tlsConfig.InsecureSkipVerify = true
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig
	return nil
}

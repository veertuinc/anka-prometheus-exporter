package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

type ClientTLSCerts struct {
	UseTLS              bool
	Cert                string
	CertKey             string
	CACert              string
	SkipTLSVerification bool
}

func appendRootCert(certFilePath string, caCertPool *x509.CertPool) error {
	cert, err := os.ReadFile(certFilePath)
	if err != nil {
		return err
	}
	ok := caCertPool.AppendCertsFromPEM(cert)
	if !ok {
		return fmt.Errorf("could not add %v to Root Certificates", certFilePath)
	}
	return nil
}

func setUpTLS(certs ClientTLSCerts) error {
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

	if certs.Cert != "" && certs.CertKey != "" {
		cert, err := tls.LoadX509KeyPair(certs.Cert, certs.CertKey)
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

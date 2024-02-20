package utils

import (
	"crypto/tls"
	"crypto/x509"

	"os"
)

type Secure struct {
	certFile string
	keyFile  string
	caFile   string

	tlsConfig *tls.Config
}

func NewSecure(certFile, keyFile, caFile string) *Secure {
	return &Secure{
		certFile: certFile,
		keyFile:  keyFile,
		caFile:   caFile,
	}
}

func (s *Secure) GetTLSConfig(mtls bool) (*tls.Config, error) {
	if s.tlsConfig != nil {
		if mtls {
			s.tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			s.tlsConfig.ClientAuth = tls.NoClientCert
		}

		return s.tlsConfig, nil
	}

	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(s.caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	s.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
	}

	if mtls {
		s.tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		s.tlsConfig.ClientAuth = tls.NoClientCert
	}

	return s.tlsConfig, nil
}

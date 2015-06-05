package common

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/olebedev/config"
)

// DefaultAddress Bind to all available network interfaces by default
const DefaultAddress = ""

// DefaultPort Default port number
const DefaultPort = "44123"

type GozwConfig struct {
	// Shared
	PortalAddress string
	Certificate   []tls.Certificate
	RootCAs       *x509.CertPool

	// Gateway
	Device string
	Baud   int
}

func stringList(list []interface{}) []string {
	strings := make([]string, len(list))

	for i, v := range list {
		strings[i] = v.(string)
	}

	return strings
}

func loadCertPair(prefix string, config *config.Config) ([]tls.Certificate, error) {
	cert, err := config.String(prefix + ".cert")
	if err != nil {
		return nil, err
	}

	key, err := config.String(prefix + ".key")
	if err != nil {
		return nil, err
	}

	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	certList := make([]tls.Certificate, 1)
	certList[0] = certificate

	return certList, nil
}

func loadCertPoolFromFiles(certPaths []string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for i := 0; i < len(certPaths); i++ {
		pemCerts, err := ioutil.ReadFile(certPaths[i])
		if err != nil {
			return nil, err
		}

		for len(pemCerts) > 0 {
			var block *pem.Block
			block, pemCerts = pem.Decode(pemCerts)

			if block == nil {
				break
			}

			if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
				continue
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}

			certPool.AddCert(cert)
		}
	}

	return certPool, nil
}

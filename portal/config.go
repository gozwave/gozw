package portal

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

type portalConfig struct {
	listenAddress     string
	serverCertificate []tls.Certificate
	rootCAs           *x509.CertPool
}

func stringList(list []interface{}) []string {
	strings := make([]string, len(list))

	for i, v := range list {
		strings[i] = v.(string)
	}

	return strings
}

func LoadConfigFromYaml(configFile string) (*portalConfig, error) {
	config, err := config.ParseYamlFile(configFile)
	if err != nil {
		return nil, err
	}

	bindPort, err := config.String("gozwd.bindPort")
	if err != nil || bindPort == "" {
		bindPort = DefaultPort
	}

	bindAddress, err := config.String("gozwd.bindAddress")
	if err != nil {
		bindAddress = DefaultAddress
	}

	cert, err := config.String("gozwd.cert")
	if err != nil {
		return nil, err
	}

	key, err := config.String("gozwd.key")
	if err != nil {
		return nil, err
	}

	ca, err := config.List("gozwd.ca")
	if err != nil {
		return nil, err
	}

	// clientCA, err := config.List("gozwd.clientCA")
	// if err != nil {
	// 	return nil, err
	// }

	serverCertificate, err := loadCertificateList(cert, key)
	if err != nil {
		return nil, err
	}

	rootCAs, err := loadCertPoolFromFiles(stringList(ca))
	if err != nil {
		return nil, err
	}

	// clientCAs, err := loadCertPoolFromFiles(stringList(clientCA))
	// if err != nil {
	//   return nil, err
	// }

	portalConfig := portalConfig{
		listenAddress:     bindAddress + ":" + bindPort,
		serverCertificate: serverCertificate,
		rootCAs:           rootCAs,
		// clientCAs:         clientCAs,
	}

	return &portalConfig, nil
}

func (config *portalConfig) GetTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: config.serverCertificate,
		// @todo change to RequireAndVerifyClientCert when I understand SSL certs
		// ClientAuth: tls.RequireAndVerifyClientCert,
		RootCAs: config.rootCAs,
		// ClientCAs: config.clientCAs,
	}
}

func (config *portalConfig) GetListenAddress() string {
	return config.listenAddress
}

func loadCertificateList(certPath string, keyPath string) ([]tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	certList := make([]tls.Certificate, 1)
	certList[0] = cert

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

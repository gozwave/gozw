package common

import (
	"errors"

	"github.com/olebedev/config"
)

func LoadPortalConfig(configFile string) (*GozwConfig, error) {

	config, err := config.ParseYamlFile(configFile)
	if err != nil {
		return nil, err
	}

	portalAddress, err := config.String("settings.portal.address")
	if err != nil || portalAddress == "" {
		return nil, errors.New("Could not read required 'settings.portal.address'")
	}

	certificate, err := loadCertPair("settings.portal", config)
	if err != nil {
		return nil, err
	}

	ca, err := config.List("settings.ca")
	if err != nil {
		return nil, err
	}

	// clientCA, err := config.List("settings.clientCA")
	// if err != nil {
	// 	return nil, err
	// }

	rootCAs, err := loadCertPoolFromFiles(stringList(ca))
	if err != nil {
		return nil, err
	}

	// clientCAs, err := loadCertPoolFromFiles(stringList(clientCA))
	// if err != nil {
	//   return nil, err
	// }

	clientConfig := GozwConfig{
		PortalAddress: portalAddress,
		Certificate:   certificate,
		RootCAs:       rootCAs,
		// clientCAs:         clientCAs,
	}

	return &clientConfig, nil
}

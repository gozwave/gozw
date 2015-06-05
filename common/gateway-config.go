package common

import (
	"errors"

	"github.com/olebedev/config"
)

func LoadGatewayConfig(configFile string) (*GozwConfig, error) {

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

	device, err := config.String("settings.gateway.device")
	if err != nil {
		return nil, err
	}

	baud, err := config.Int("settings.gateway.baud")
	if err != nil {
		return nil, err
	}

	clientConfig := GozwConfig{
		PortalAddress: portalAddress,
		Certificate:   certificate,
		RootCAs:       rootCAs,

		Device: device,
		Baud:   baud,
	}

	return &clientConfig, nil
}

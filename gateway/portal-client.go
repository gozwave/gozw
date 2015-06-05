package gateway

import "crypto/tls"

func NewPortalClient() {
	tlsConfig := tls.Config{}
	tls.Dial("tcp", "127.0.0.1:44123", &tlsConfig)
}

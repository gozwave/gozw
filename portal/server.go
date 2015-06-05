package portal

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/bjyoungblood/gozw/common"
)

// PortalServer is the server that Z/IP gateways connect to
type PortalServer struct {
	listener net.Listener
	clients  map[net.Conn]Client
}

func NewPortalServer(config *common.GozwConfig) (*PortalServer, error) {
	server := new(PortalServer)

	tlsConfig := tls.Config{
		Certificates: config.Certificate,
		// @todo change to RequireAndVerifyClientCert when I understand SSL certs
		// ClientAuth: tls.RequireAndVerifyClientCert,
		RootCAs: config.RootCAs,
		// ClientCAs: config.clientCAs,
	}

	listener, err := tls.Listen("tcp4", config.PortalAddress, &tlsConfig)
	if err != nil {
		return nil, err
	}

	fmt.Println("Listening on " + config.PortalAddress)

	server.listener = listener
	return server, nil
}

func (psock *PortalServer) Start() {
	// Close the listener when the application closes
	defer psock.Close()

	for {
		// Listen for an incoming connection

		conn, err := psock.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("Got a client")
		client := NewClient(conn)
		defer client.Close()

		go client.Handle()
	}
}

func (t *PortalServer) Accept() (net.Conn, error) {
	conn, err := t.listener.Accept()
	return conn, err
}

func (t *PortalServer) Close() error {
	return t.listener.Close()
}

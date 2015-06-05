package portal

import (
	"crypto/tls"
	"fmt"
	"net"
)

// PortalServer is the server that Z/IP gateways connect to
type PortalServer struct {
	listener net.Listener
	clients  map[net.Conn]Client
}

func NewPortalServer(config *portalConfig) (*PortalServer, error) {
	server := new(PortalServer)

	listener, err := tls.Listen("tcp4", config.GetListenAddress(), config.GetTLSConfig())
	if err != nil {
		return nil, err
	}

	fmt.Println("Listening on " + config.GetListenAddress())

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

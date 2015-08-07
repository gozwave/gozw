package arcturus

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	listener    net.Listener
	connections map[uint32]*Conn
}

func NewServer() *Server {
	return &Server{
		connections: map[uint32]*Conn{},
	}
}

func (s *Server) Listen(nettype, addr string) {
	var err error

	s.listener, err = net.Listen(nettype, addr)
	if err != nil {
		panic(err)
	}

	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	c := NewConn(conn)

	select {
	case homeId := <-c.Authenticated:
		fmt.Println("Authentication ok")
		s.connections[homeId] = c
		go c.processOutgoing()
	case <-time.After(time.Second * 5):
		fmt.Println("Authentication timeout")
		c.close()
	}

	<-c.Closed
}

func (s *Server) Close() error {
	return s.listener.Close()
}

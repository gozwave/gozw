package portal

import (
	"fmt"
	"net"
)

// Client represents a connection from a Z/IP gateway
type Client struct {
	state int
	conn  net.Conn
	ch    chan []byte
}

func NewClient(conn net.Conn) Client {
	return Client{
		conn: conn,
	}
}

func (c *Client) Handle() {
	fmt.Println("got a client!")
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

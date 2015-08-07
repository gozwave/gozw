package arcturus

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/bjyoungblood/gozw/proto"
	"github.com/davecgh/go-spew/spew"
)

type ConnState int

const (
	ConnStateInit ConnState = iota
	ConnStateReady
)

type Conn struct {
	conn  net.Conn
	state ConnState

	outgoing      chan proto.Event
	Authenticated chan uint32
	Closed        chan bool
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		conn:  conn,
		state: ConnStateInit,

		outgoing:      make(chan proto.Event, 10),
		Authenticated: make(chan uint32, 1),
		Closed:        make(chan bool, 1),
	}

	go c.processIncoming()

	return c
}

func (c *Conn) close() {
	c.conn.Close()
	c.Closed <- true
}

func (c *Conn) processOutgoing() {
	encoder := gob.NewEncoder(c.conn)

	for ev := range c.outgoing {
		err := encoder.Encode(ev)
		if err != nil {
			fmt.Printf("Encoding error: %v\n", err)
		}
	}
}

func (c *Conn) processIncoming() {
	decoder := gob.NewDecoder(c.conn)

	for {
		event := proto.Event{}
		err := decoder.Decode(&event)
		if err == io.EOF {
			c.close()
			return
		} else if err != nil {
			fmt.Printf("Decoding error: %v\n", err)
			c.close()
			return
		}

		c.handleEvent(event)
	}
}

func (c *Conn) handleEvent(ev proto.Event) {
	if c.state != ConnStateReady {
		event, ok := ev.Payload.(proto.IdentEvent)
		if !ok {
			fmt.Println("Did not receive expected IdentEvent")
			c.close()
			return
		} else {
			c.state = ConnStateReady
			c.Authenticated <- event.HomeId
		}
	}

	spew.Dump(ev.Payload)
}

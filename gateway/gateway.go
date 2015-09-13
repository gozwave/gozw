package gateway

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/helioslabs/gozw/zwave/application"
	"github.com/helioslabs/gozw/zwave/frame"
	"github.com/helioslabs/gozw/zwave/serial-api"
	"github.com/helioslabs/gozw/zwave/session"
	"github.com/helioslabs/gozw/zwave/transport"
	"github.com/helioslabs/proto"
)

type GatewayOptions struct {
	CommNetType string
	CommAddress string

	ZWaveSerialPort string
	BaudRate        int
}

type Gateway struct {
	opts GatewayOptions
	app  *application.Layer
	conn net.Conn

	outgoingEvents chan proto.Event
}

func NewGateway(opts GatewayOptions) (*Gateway, error) {
	gateway := &Gateway{
		opts:           opts,
		outgoingEvents: make(chan proto.Event, 1),
	}

	if err := gateway.openCommPort(); err != nil {
		return nil, err
	}

	if err := gateway.zwaveStart(); err != nil {
		return nil, err
	}

	return gateway, nil
}

func (g *Gateway) openCommPort() error {
	conn, err := net.Dial(g.opts.CommNetType, g.opts.CommAddress)
	if err != nil {
		return err
	}

	g.conn = conn

	return nil
}

func (g *Gateway) zwaveStart() error {
	transport, err := transport.NewSerialPortTransport(g.opts.ZWaveSerialPort, g.opts.BaudRate)
	if err != nil {
		return err
	}

	frameLayer := frame.NewFrameLayer(transport)
	sessionLayer := session.NewSessionLayer(frameLayer)
	apiLayer := serialapi.NewLayer(sessionLayer)
	appLayer, err := application.NewLayer(apiLayer)
	if err != nil {
		return err
	}

	g.app = appLayer

	return nil
}

func (g *Gateway) Run() {
	g.subscribeToAppEvents()

	go g.processOutgoing()
	go g.processIncoming()

	g.outgoingEvents <- proto.Event{
		Payload: proto.IdentEvent{HomeID: g.app.Controller.HomeID},
	}
}

func (g *Gateway) Shutdown() {
	g.app.Shutdown()
	g.conn.Close()
}

func (g *Gateway) subscribeToAppEvents() {
	g.app.EventBus.SubscribeAsync("event", func(ev proto.Event) {
		g.outgoingEvents <- ev
	}, true)
}

func (g *Gateway) processOutgoing() {
	encoder := gob.NewEncoder(g.conn)

	for ev := range g.outgoingEvents {
		err := encoder.Encode(ev)
		if err != nil {
			fmt.Printf("Encoding error; %v\n", err)
		}
	}
}

func (g *Gateway) processIncoming() {
	decoder := gob.NewDecoder(g.conn)

	for {
		event := proto.Event{}
		err := decoder.Decode(&event)
		if err == io.EOF {
			// @todo initiate reconnect sequence
			fmt.Println("EOF!")
			break
		} else if err != nil {
			fmt.Printf("Decoding error: %v\n", err)
			continue
		}

		g.handleEvent(event)
	}
}

func (g *Gateway) handleEvent(ev proto.Event) {
	switch ev.Payload.(type) {
	case proto.RegisterEvent:
		g.outgoingEvents <- proto.Event{
			Payload: proto.ControllerInfoEvent{
				APIVersion:          g.app.Controller.APIVersion,
				APILibraryType:      g.app.Controller.APILibraryType,
				HomeID:              g.app.Controller.HomeID,
				NodeID:              g.app.Controller.NodeID,
				Version:             g.app.Controller.Version,
				APIType:             g.app.Controller.APIType,
				IsPrimaryController: g.app.Controller.IsPrimaryController,
				ApplicationVersion:  g.app.Controller.ApplicationVersion,
				ApplicationRevision: g.app.Controller.ApplicationRevision,
				SupportedFunctions:  g.app.Controller.SupportedFunctions,
			},
		}

		for _, node := range g.app.Nodes() {

			payload := proto.NodeInfoEvent{
				NodeID:              node.NodeID,
				Capability:          node.Capability,
				BasicDeviceClass:    node.BasicDeviceClass,
				GenericDeviceClass:  node.GenericDeviceClass,
				SpecificDeviceClass: node.SpecificDeviceClass,
				Failing:             node.Failing,
				CommandClasses:      node.CommandClasses,
				ManufacturerID:      node.ManufacturerID,
				ProductTypeID:       node.ProductTypeID,
				ProductID:           node.ProductID,
			}

			g.outgoingEvents <- proto.Event{
				Payload: payload,
			}

		}

	case proto.NodeCommandEvent:
		cmd := ev.Payload.(proto.NodeCommandEvent)
		node, err := g.app.SendData(cmd.NodeID, cmd.CommandData)

	default:
		spew.Dump(ev.Payload)
	}
}

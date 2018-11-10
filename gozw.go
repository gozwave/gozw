package gozw

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gozwave/gozw/cc"
	zwsec "github.com/gozwave/gozw/cc/security"
	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/security"
	"github.com/gozwave/gozw/serialapi"
	"github.com/gozwave/gozw/session"
	"github.com/gozwave/gozw/transport"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// MaxSecureInclusionDuration is the timeout for secure inclusion mode. If this
// timeout expires, secure inclusion will be canceled no matter how far the
// process has proceeded.
const MaxSecureInclusionDuration = time.Second * 60

// Maximum possible size (in bytes) of the plaintext payload to be sent when
// sending a secure frame, based on the SendData options. The smallest possible
// must be used based on the given option bitset (e.g. if using both no route
// and explore, the maximum size is 26).
const (
	SecurePayloadMaxSizeExplore   = 26 // in bytes
	SecurePayloadMaxSizeAutoRoute = 28
	SecurePayloadMaxSizeNoRoute   = 34
)

type Client struct {
	Controller Controller

	serialAPI     serialapi.ILayer
	securityLayer security.ILayer

	networkKey []byte
	nodes      map[byte]*Node

	// REPLACE THIS WITH A GENERIC CALLBACK FUNCTION
	// EventBus EventBus.Bus
	EventCallback func(*Client, byte, cc.Command)

	l  *zap.Logger
	db *bolt.DB

	ctx    context.Context
	cancel context.CancelFunc

	secureInclusionStep map[byte]chan error
}

func NewDefaultClient(dbName, serialPort string, baudRate int, networkKey []byte) (*Client, error) {
	logger, err := NewLogger()
	if err != nil {
		return nil, errors.Wrap(err, "initialize logger")
	}

	client := Client{
		Controller:          Controller{},
		networkKey:          networkKey,
		nodes:               map[byte]*Node{},
		EventCallback:       DefaultEventCallback,
		l:                   logger,
		secureInclusionStep: map[byte]chan error{},
	}

	client.ctx, client.cancel = context.WithCancel(context.Background())

	transport, err := transport.NewSerialPortTransport(serialPort, baudRate)
	if err != nil {
		return nil, errors.Wrap(err, "initializing transport")
	}

	frameLayer, err := frame.NewFrameLayer(client.ctx, transport, logger)
	if err != nil {
		return nil, errors.Wrap(err, "initialize frame layer")
	}

	sessionLayer := session.NewSessionLayer(client.ctx, frameLayer, logger)

	client.serialAPI = serialapi.NewLayer(client.ctx, sessionLayer, logger)

	client.securityLayer = security.NewLayer(client.networkKey, logger)

	err = client.initDb(dbName)
	if err != nil {
		return nil, errors.Wrap(err, "initialize db")
	}

	go client.handleApplicationCommands()
	go client.handleControllerUpdates()

	err = client.initZWave()
	if err != nil {
		return nil, errors.Wrap(err, "initializing z-wave")
	}

	return &client, nil
}

func (c *Client) SetLogger(logger *zap.Logger) {
	c.l = logger
}

func (c *Client) initDb(dbName string) (err error) {
	c.db, err = bolt.Open(dbName, 0600, &bolt.Options{})
	if err != nil {
		return
	}

	c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("controller"))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

//  NewLogger builds a  new logger.
func NewLogger() (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build logger")
	}

	return logger, nil
}

// Nodes will return all nodes
func (c *Client) Nodes() map[byte]*Node {
	return c.nodes
}

// Node will retrieve a single node.
func (c *Client) Node(nodeID byte) (*Node, error) {
	if node, ok := c.nodes[nodeID]; ok {
		return node, nil
	}

	return nil, errors.New("Node not found")
}

func (c *Client) initZWave() error {
	version, err := c.serialAPI.GetVersion()
	if err != nil {
		return err
	}

	c.Controller.APIVersion = version.Version
	c.Controller.APILibraryType = version.GetLibraryTypeString()

	c.Controller.HomeID, c.Controller.NodeID, err = c.serialAPI.MemoryGetID()
	if err != nil {
		return err
	}

	serialAPICapabilities, err := c.serialAPI.GetCapabilities()
	if err != nil {
		return err
	}

	c.Controller.ApplicationVersion = serialAPICapabilities.ApplicationVersion
	c.Controller.ApplicationRevision = serialAPICapabilities.ApplicationRevision
	c.Controller.SupportedFunctions = serialAPICapabilities.GetSupportedFunctions()

	initData, err := c.serialAPI.GetInitAppData()
	if err != nil {
		return err
	}

	c.Controller.Version = initData.Version
	c.Controller.APIType = initData.GetAPIType()
	c.Controller.IsPrimaryController = initData.IsPrimaryController()
	c.Controller.NodeList = initData.GetNodeIDs()

	for _, nodeID := range c.Controller.NodeList {
		node, err := NewNode(c, nodeID)

		if err != nil {
			spew.Dump(err)
			continue
		}

		c.nodes[nodeID] = node
	}

	return nil
}

// Shutdown will stop the client.
func (c *Client) Shutdown() error {
	c.cancel()
	return nil
}

func (c *Client) AddNode() (*Node, error) {
	newNodeInfo, err := c.serialAPI.AddNode()
	if err != nil {
		return nil, err
	}

	if newNodeInfo == nil {
		return nil, errors.New("Adding node failed")
	}

	node, err := NewNode(c, newNodeInfo.Source)
	if err != nil {
		return nil, err
	}

	node.setFromAddNodeCallback(newNodeInfo)
	c.nodes[node.NodeID] = node

	if node.IsSecure() {
		c.l.Debug("starting secure inclusion")
		err = c.includeSecureNode(node)
		if err != nil {
			return nil, err
		}
	}

	node.nextQueryStage()

	select {
	case <-node.queryStageVersionsComplete:
		c.l.Info("node queries complete")
	case <-time.After(time.Second * 30):
		c.l.Warn("node query timeout", zap.String("node", fmt.Sprint(node.NodeID)))
	}

	node.AddAssociation(1, 1)

	return node, nil
}

func (c *Client) RemoveNode() (byte, error) {
	result, err := c.serialAPI.RemoveNode()
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, errors.New("Removing node failed")
	}

	return result.Source, nil
}

func (c *Client) RemoveFailedNode(nodeID byte) (ok bool, err error) {
	return c.serialAPI.RemoveFailedNode(nodeID)
}

func (c *Client) handleApplicationCommands() {
	for {
		select {
		case cmd := <-c.serialAPI.ControllerCommands():
			switch cc.CommandClassID(cmd.CommandData[0]) {

			case cc.Security:
				c.interceptSecurityCommandClass(cmd)

			default:
				if node, err := c.Node(cmd.SrcNodeID); err == nil {
					go node.receiveApplicationCommand(cmd)
				} else {
					c.l.Warn("Received command for unknown node", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
				}

			}
		case <-c.ctx.Done():
			c.l.Info("stopping application commands handler")
			return
		}
	}
}

// SetEventCallback will set the event callback for any events received
func (c *Client) SetEventCallback(callback func(c *Client, nodeID byte, e cc.Command)) {
	c.EventCallback = callback
}

// DefaultEventCallback is the default callback for handling events.
func DefaultEventCallback(c *Client, nodeID byte, e cc.Command) {
	c.l.Info("event received", zap.Any("event", e), zap.Int("nodeID", int(nodeID)))
}

func (c *Client) handleControllerUpdates() {
	for {
		select {
		case update := <-c.serialAPI.ControllerUpdates():
			switch update.Status {

			case protocol.UpdateStateNodeInfoReceived,
				protocol.UpdateStateNodeInfoReqFailed:
				if node, ok := c.nodes[update.NodeID]; ok {
					node.receiveControllerUpdate(update)
				} else {
					c.l.Debug("controller update:", zap.String("data", spew.Sdump(update)))
				}

			default:
				c.l.Debug("controller update:", zap.String("data", spew.Sdump(update)))

			}
		case <-c.ctx.Done():
			c.l.Info("stopping controller updates handler")
			return
		}
	}
}

func (c *Client) SendData(dstNode byte, payload encoding.BinaryMarshaler) error {
	marshaled, err := payload.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.serialAPI.SendData(dstNode, marshaled)
	return err
}

// SendDataSecure encapsulates payload in a security encapsulation command and
// sends it to the destination node.
func (c *Client) SendDataSecure(dstNode byte, message encoding.BinaryMarshaler) error {
	// This function wraps the private sendDataSecure because no external packages
	// should ever call this while in inclusion mode (and doing so would be incorrect)
	return c.sendDataSecure(dstNode, message, false)
}

func (c *Client) requestNonceForNode(dstNode byte) (security.Nonce, error) {
	err := c.SendData(dstNode, &zwsec.NonceGet{})
	if err != nil {
		return nil, err
	}

	return c.securityLayer.WaitForExternalNonce(dstNode)
}

func (c *Client) getOrRequestNonceForNode(dstNode byte) (nonce security.Nonce, err error) {
	if nonce, err = c.securityLayer.GetExternalNonce(dstNode); err == nil {
		return nonce, nil
	}

	for i := 0; i < 3; i++ {
		nonce, err = c.requestNonceForNode(dstNode)
		if err == nil {
			break
		}

		c.l.Error("error: get nonce attempt failed", zap.Int("attempt", i))
		time.Sleep(50 * time.Millisecond)
	}

	return nonce, err
}

func (c *Client) sendDataSecure(dstNode byte, message encoding.BinaryMarshaler, inclusionMode bool) error {
	// Previously, this function would just split and prepare the payload based on
	// whether it should be split after figuring out whether to segment. For now,
	// we're just going to assume that we will never have to worry about segmenting.
	// It wasn't too hard to implement before, but since I couldn't find a real payload
	// big enough, it wasn't possible to verify the implementation, so I didn't port
	// it while refactoring (for simplicity's sake).

	payload, err := message.MarshalBinary()
	if err != nil {
		return err
	}

	// Get a nonce from the other node
	receiverNonce, err := c.getOrRequestNonceForNode(dstNode)
	if err != nil {
		return err
	}

	senderNonce, err := c.securityLayer.GenerateInternalNonce()
	if err != nil {
		return err
	}

	var securityByte byte
	// var securityByte byte = sequenceCounter & SecuritySequenceCounterMask
	// if sequenced {
	// 	securityByte |= SecuritySequenceSequencedFlag
	//
	// 	if isSecondFrame {
	// 		securityByte |= SecuritySequenceSecondFrameFlag
	// 	}
	// }

	securePayload := append([]byte{securityByte}, payload...)

	encapsulatedMessage, err := c.securityLayer.EncapsulateMessage(
		1,
		dstNode,
		zwsec.CommandMessageEncapsulation, // @todo CC should be determined by sequencing
		senderNonce,
		receiverNonce,
		securePayload,
		inclusionMode,
	)

	if err != nil {
		c.l.Error("failed to encrypt message", zap.String("err", err.Error()), zap.String("node", fmt.Sprint(dstNode)))
		return err
	}

	return c.SendData(dstNode, encapsulatedMessage)
}

func (c *Client) includeSecureNode(node *Node) error {
	c.secureInclusionStep[node.NodeID] = make(chan error)
	c.SendData(node.NodeID, &zwsec.SchemeGet{})

	defer close(c.secureInclusionStep[node.NodeID])
	defer delete(c.secureInclusionStep, node.NodeID)

	c.l.Info("requesting security scheme")

	select {
	case err := <-c.secureInclusionStep[node.NodeID]:
		if err != nil {
			return err
		}
	case <-time.After(time.Second * 10):
		return errors.New("Secure inclusion timeout")
	}

	c.l.Info("sending network key")
	node.NetworkKeySent = true

	c.sendDataSecure(
		node.NodeID,
		&zwsec.NetworkKeySet{NetworkKeyByte: c.networkKey},
		true,
	)

	select {
	case err := <-c.secureInclusionStep[node.NodeID]:
		return err
	case <-time.After(time.Second * 20):
		return errors.New("Secure inclusion timeout")
	}
}

func (c *Client) interceptSecurityCommandClass(cmd serialapi.ApplicationCommand) {
	command, err := cc.Parse(1, cmd.CommandData)
	if err != nil {
		c.l.Error(err.Error())
		return
	}

	switch command.(type) {

	case *zwsec.MessageEncapsulation, *zwsec.MessageEncapsulationNonceGet:
		c.l.Info("rx secure message", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		// @todo determine whether to bother with sequenced messages. According to
		// openzwave, they didn't bother to implement it because they never ran across
		// a situation where a frame was large enough that it needed to be sequenced.
		// in any case, the following is the following is the process to follow with
		// or without sequencing:

		// 1. decrypt message
		// 2. if it's the first half of a sequenced message, wait for the second half
		// 2.5  if it's an EncapsulationGetNonce, then send a NonceReport back to the sender
		// 3. if it's the second half of a sequenced message, reassemble the payloads
		// 4. emit the decrypted (possibly recombined) message back

		var decrypted []byte
		node, err := c.Node(cmd.SrcNodeID)
		if err != nil {
			return
		}

		if !node.NetworkKeySent {
			decrypted, err = c.securityLayer.DecryptMessage(cmd, true)
		} else {
			decrypted, err = c.securityLayer.DecryptMessage(cmd, false)
		}

		if err != nil {
			return
		}

		c.l.Info("received encapsulated message", zap.String("data", spew.Sdump(decrypted)))

		if decrypted[1] == byte(cc.Security) &&
			decrypted[2] == byte(zwsec.CommandNetworkKeyVerify) {
			c.l.Info("network key verify", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
			if ch, ok := c.secureInclusionStep[cmd.SrcNodeID]; ok {
				ch <- nil
			}
			return
		}

		if node, ok := c.nodes[cmd.SrcNodeID]; ok {
			cmd.CommandData = decrypted[1:]
			go node.receiveApplicationCommand(cmd)
		} else {
			c.l.Warn("received secure command for unknown node", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		}

	case *zwsec.NonceGet:
		c.l.Info("nonce get", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		nonce, err := c.securityLayer.GenerateInternalNonce()
		if err != nil {
		}

		reply := &zwsec.NonceReport{NonceByte: nonce}
		c.SendData(cmd.SrcNodeID, reply)

	case *zwsec.NonceReport:
		c.l.Info("nonce report", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		c.securityLayer.ReceiveNonce(cmd.SrcNodeID, *command.(*zwsec.NonceReport))

	case *zwsec.SchemeReport:
		c.l.Info("security scheme report", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		if ch, ok := c.secureInclusionStep[cmd.SrcNodeID]; ok {
			ch <- nil
		} else {
			c.l.Warn("not in secure inclusion mode", zap.String("node", fmt.Sprint(cmd.SrcNodeID)))
		}

	default:
		c.l.Warn("unexpected security command", zap.String("data", spew.Sdump(cmd)))
	}
}

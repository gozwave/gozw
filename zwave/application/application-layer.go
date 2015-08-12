package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
	"github.com/helioslabs/gozw/zwave/command-class"
	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/security"
	"github.com/helioslabs/gozw/zwave/serial-api"
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

// These are used in the security message encapsulation command to indicate how
// the message is sequenced.
// const (
// 	// securitySequenceSequencedFlag indicates that the message is sequenced if set.
// 	securitySequenceSequencedFlag byte = 0x10
//
// 	// securitySequenceSecondFrameFlag indicates that the message is the second in
// 	// the sequence if set.
// 	securitySequenceSecondFrameFlag = 0x20
//
// 	// securitySequenceCounterMask masks the non-counter bytes from the sequence
// 	// counter in the security byte
// 	securitySequenceCounterMask = 0x0f
// )

// Layer is the top-level controlling layer for the Z-Wave network. It maintains
// information about the controller itself, as well as a list of network nodes.
// It is responsible for routing messages between the Z-Wave controller and the
// in-memory representations of network nodes. This involves coordinating Z-Wave
// security functions (encrypting/decrypting messages, fetching nonces, etc.)
// and interaction with the Serial API layer.
type Layer struct {

	// APIVersion is the Z-Wave serial API Version
	APIVersion string

	// APILibraryType is the Z-Wave library type string (as returned from the
	// controller)
	APILibraryType string

	// HomeID is the controller's home ID
	HomeID uint32

	// NodeID is the controller's node ID
	NodeID byte

	// Version is returned from the GetVersion serial API call
	Version byte

	// APIType
	APIType string

	// IsPrimaryController is true when this node is the network's primary
	// controller (we don't currently support being a non-primary controller)
	IsPrimaryController bool

	// ApplicationVersion
	ApplicationVersion  byte
	ApplicationRevision byte
	SupportedFunctions  []byte

	NodeList []byte

	serialAPI     serialapi.ILayer
	securityLayer security.ILayer
	networkKey    []byte
	nodes         map[byte]*Node

	EventBus *EventBus.EventBus

	db *bolt.DB

	// maps node id to channel
	secureInclusionStep map[byte]chan error
}

// NewLayer instantiates a new application layer, handles opening and setting up
// the local database, reads (or generates) a network key, starts threads to
// handle incming Z-Wave commands and updates, and loads basic controller data
// from the Z-Wave controller.
func NewLayer(serialAPI serialapi.ILayer) (app *Layer, err error) {
	app = &Layer{
		serialAPI: serialAPI,
		nodes:     map[byte]*Node{},

		EventBus: EventBus.New(),

		secureInclusionStep: map[byte]chan error{},
	}

	err = app.initDb()
	if err != nil {
		return
	}

	networkKey, err := app.initNetworkKey()
	if err != nil {
		return
	}

	app.networkKey = networkKey
	app.securityLayer = security.NewLayer(networkKey)

	go app.handleApplicationCommands()
	go app.handleControllerUpdates()

	err = app.initZWave()

	return
}

func (a *Layer) Nodes() map[byte]*Node {
	return a.nodes
}

func (a *Layer) Node(nodeID byte) (*Node, error) {
	if node, ok := a.nodes[nodeID]; ok {
		return node, nil
	}

	return nil, errors.New("Node not found")
}

func (a *Layer) initDb() (err error) {
	a.db, err = bolt.Open("data.db", 0600, &bolt.Options{})
	if err != nil {
		return
	}

	a.db.Update(func(tx *bolt.Tx) error {
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

func (a *Layer) initNetworkKey() ([]byte, error) {
	var networkKey []byte

	err := a.db.View(func(tx *bolt.Tx) error {
		networkKey = tx.Bucket([]byte("controller")).Get([]byte("networkKey"))
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(networkKey) == 16 {
		return networkKey, nil
	}

	networkKey = security.GenerateNetworkKey()

	err = a.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("controller")).Put([]byte("networkKey"), networkKey)
	})

	if err != nil {
		return nil, err
	}

	return networkKey, nil
}

func (a *Layer) initZWave() error {
	version, err := a.serialAPI.GetVersion()
	if err != nil {
		return err
	}

	a.APIVersion = version.Version
	a.APILibraryType = version.GetLibraryTypeString()

	a.HomeID, a.NodeID, err = a.serialAPI.MemoryGetID()
	if err != nil {
		return err
	}

	serialAPICapabilities, err := a.serialAPI.GetSerialAPICapabilities()
	if err != nil {
		return err
	}

	a.ApplicationVersion = serialAPICapabilities.ApplicationVersion
	a.ApplicationRevision = serialAPICapabilities.ApplicationRevision
	a.SupportedFunctions = serialAPICapabilities.GetSupportedFunctions()

	initData, err := a.serialAPI.GetInitAppData()
	if err != nil {
		return err
	}

	a.Version = initData.Version
	a.APIType = initData.GetAPIType()
	a.IsPrimaryController = initData.IsPrimaryController()
	a.NodeList = initData.GetNodeIDs()

	for _, nodeID := range a.NodeList {
		node, err := NewNode(a, nodeID)

		if err != nil {
			spew.Dump(err)
			continue
		}

		a.nodes[nodeID] = node
	}

	return nil
}

func (a *Layer) Shutdown() error {
	return a.db.Close()
}

func (a *Layer) AddNode() (*Node, error) {
	newNodeInfo, err := a.serialAPI.AddNode()
	if err != nil {
		return nil, err
	}

	if newNodeInfo == nil {
		return nil, errors.New("Adding node failed")
	}

	node, err := NewNode(a, newNodeInfo.Source)
	if err != nil {
		return nil, err
	}

	node.setFromAddNodeCallback(newNodeInfo)
	a.nodes[node.NodeID] = node

	if node.IsSecure() {
		fmt.Println("Starting secure inclusion")
		err = a.includeSecureNode(node.NodeID)
		if err != nil {
			return nil, err
		}

		time.Sleep(time.Millisecond * 50)
		err := node.RequestSupportedSecurityCommands()
		if err != nil {
			fmt.Println(err)
		}

		select {
		case <-node.receivedSecurityInfo:
		case <-time.After(time.Second * 5):
			fmt.Println("timed out after requesting security commands")
		}
	}

	spew.Dump(node)

	node.AddAssociation(1, 1)

	return node, nil
}

func (a *Layer) RemoveNode() (byte, error) {
	result, err := a.serialAPI.RemoveNode()

	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, errors.New("Removing node failed")
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("nodes")).Delete([]byte{result.Source})
	})

	if err != nil {
		return 0, err
	}

	return result.Source, nil
}

func (a *Layer) RemoveFailedNode(nodeID byte) (ok bool, err error) {
	ok, err = a.serialAPI.RemoveFailedNode(nodeID)

	if ok && err != nil {
		err = a.db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("nodes")).Delete([]byte{nodeID})
		})
	}

	return
}

func (a *Layer) handleApplicationCommands() {
	for cmd := range a.serialAPI.ControllerCommands() {
		switch cmd.CommandData[0] {

		case commandclass.CommandClassSecurity:
			a.handleSecurityCommand(cmd)

		default:
			if node, err := a.Node(cmd.SrcNodeID); err == nil {
				go node.receiveApplicationCommand(cmd)
			} else {
				fmt.Println("Received command for unknown node", cmd.SrcNodeID)
			}

		}

	}
}

func (a *Layer) handleControllerUpdates() {
	for update := range a.serialAPI.ControllerUpdates() {

		switch update.Status {

		case protocol.UpdateStateNodeInfoReceived,
			protocol.UpdateStateNodeInfoReqFailed:
			if node, ok := a.nodes[update.NodeID]; ok {
				node.receiveControllerUpdate(update)
			} else {
				fmt.Println("controller update:", spew.Sdump(update))
			}

		default:
			fmt.Println("controller update:", spew.Sdump(update))

		}

	}
}

func (a *Layer) SendData(dstNode byte, payload []byte) error {
	_, err := a.serialAPI.SendData(dstNode, payload)
	return err
}

// SendDataSecure encapsulates payload in a security encapsulation command and
// sends it to the destination node.
func (a *Layer) SendDataSecure(dstNode byte, payload []byte) error {
	// This function wraps the private sendDataSecure because no external packages
	// should ever call this while in inclusion mode (and doing so would be incorrect)
	return a.sendDataSecure(dstNode, payload, false)
}

func (a *Layer) requestNonceForNode(dstNode byte) (security.Nonce, error) {
	err := a.SendData(dstNode, commandclass.NewSecurityNonceGet())

	if err != nil {
		return nil, err
	}

	return a.securityLayer.WaitForExternalNonce(dstNode)
}

func (a *Layer) getOrRequestNonceForNode(dstNode byte) (nonce security.Nonce, err error) {
	if nonce, err = a.securityLayer.GetExternalNonce(dstNode); err == nil {
		return nonce, nil
	}

	for i := 0; i < 3; i++ {
		nonce, err = a.requestNonceForNode(dstNode)
		if err == nil {
			break
		}

		fmt.Printf("get nonce attempt #%d failed\n", i)
		time.Sleep(50 * time.Millisecond)
	}

	return nonce, err
}

func (a *Layer) sendDataSecure(dstNode byte, payload []byte, inclusionMode bool) error {
	// Previously, this function would just split and prepare the payload based on
	// whether it should be split after figuring out whether to segment. For now,
	// we're just going to assume that we will never have to worry about segmenting.
	// It wasn't too hard to implement before, but since I couldn't find a real payload
	// big enough, it wasn't possible to verify the implementation, so I didn't port
	// it while refactoring (for simplicity's sake).

	// Get a nonce from the other node
	receiverNonce, err := a.getOrRequestNonceForNode(dstNode)
	if err != nil {
		return err
	}

	senderNonce, err := a.securityLayer.GenerateInternalNonce()
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

	encapsulatedMessage := a.securityLayer.EncapsulateMessage(
		securePayload,
		senderNonce,
		receiverNonce,
		1,
		dstNode,
		inclusionMode,
	)

	return a.SendData(dstNode, encapsulatedMessage)
}

func (a *Layer) includeSecureNode(nodeID byte) error {
	a.secureInclusionStep[nodeID] = make(chan error)
	a.SendData(nodeID, commandclass.NewSecuritySchemeGet())

	defer close(a.secureInclusionStep[nodeID])
	defer delete(a.secureInclusionStep, nodeID)

	select {
	case err := <-a.secureInclusionStep[nodeID]:
		if err != nil {
			return err
		}
	case <-time.After(time.Second * 10):
		return errors.New("Secure inclusion timeout")
	}

	a.sendDataSecure(
		nodeID,
		commandclass.NewSecurityNetworkKeySet(a.networkKey),
		true,
	)

	select {
	case err := <-a.secureInclusionStep[nodeID]:
		return err
	case <-time.After(time.Second * 20):
		return errors.New("Secure inclusion timeout")
	}
}

func (a *Layer) handleSecurityCommand(cmd serialapi.ApplicationCommand) {
	switch cmd.CommandData[1] {

	case commandclass.CommandSecurityMessageEncapsulation, commandclass.CommandSecurityMessageEncapsulationNonceGet:
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

		data := commandclass.ParseSecurityMessageEncapsulation(cmd.CommandData)
		msg, err := a.securityLayer.DecryptMessage(data)

		if err != nil {
			fmt.Println("error handling encrypted message", err)
			return
		}

		if msg[0] == commandclass.CommandClassSecurity && msg[1] == commandclass.CommandNetworkKeyVerify {
			if ch, ok := a.secureInclusionStep[cmd.SrcNodeID]; ok {
				ch <- nil
			}
			return
		}

		if node, ok := a.nodes[cmd.SrcNodeID]; ok {
			cmd.CommandData = msg
			go node.receiveApplicationCommand(cmd)
		} else {
			fmt.Println("Received secure command for unknown node", cmd.SrcNodeID)
		}

	case commandclass.CommandSecurityNonceGet:
		nonce, err := a.securityLayer.GenerateInternalNonce()
		if err != nil {
			fmt.Println("error generating internal nonce", err)
		}

		reply := commandclass.NewSecurityNonceReport(nonce)
		a.SendData(cmd.SrcNodeID, reply)

	case commandclass.CommandSecurityNonceReport:
		nonceReport := commandclass.ParseSecurityNonceReport(cmd.CommandData)
		a.securityLayer.ReceiveNonce(cmd.SrcNodeID, nonceReport)

	case commandclass.CommandSecuritySchemeReport:
		if ch, ok := a.secureInclusionStep[cmd.SrcNodeID]; ok {
			ch <- nil
		}

	default:
		fmt.Println("Unexpected security command:", spew.Sdump(cmd))
	}
}

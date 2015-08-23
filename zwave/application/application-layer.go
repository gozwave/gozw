package application

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/boltdb/bolt"
	"github.com/comail/colog"
	"github.com/davecgh/go-spew/spew"
	"github.com/helioslabs/gozw/zwave/command-class"
	zwsec "github.com/helioslabs/gozw/zwave/command-class/security"
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

	logger *log.Logger
	db     *bolt.DB

	// maps node id to channel
	secureInclusionStep map[byte]chan error
}

// NewLayer instantiates a new application layer, handles opening and setting up
// the local database, reads (or generates) a network key, starts threads to
// handle incming Z-Wave commands and updates, and loads basic controller data
// from the Z-Wave controller.
func NewLayer(serialAPI serialapi.ILayer) (app *Layer, err error) {
	applicationLogger := colog.NewCoLog(os.Stdout, "application ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	applicationLogger.ParseFields(true)

	app = &Layer{
		serialAPI: serialAPI,
		nodes:     map[byte]*Node{},

		EventBus: EventBus.New(),

		logger: applicationLogger.NewLogger(),

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
		a.logger.Println("debug: starting secure inclusion")
		err = a.includeSecureNode(node)
		if err != nil {
			return nil, err
		}

		err := node.RequestSupportedSecurityCommands()
		if err != nil {
			a.logger.Printf("error: %v\n", err)
		}

		select {
		case <-node.receivedSecurityInfo:
		case <-time.After(time.Second * 10):
			a.logger.Println("error: timed out after requesting security commands")
		}
	}

	spew.Dump(node)

	node.LoadCommandClassVersions()

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
		switch commandclass.ID(cmd.CommandData[0]) {

		case commandclass.Security:
			a.interceptSecurityCommandClass(cmd)

		default:
			if node, err := a.Node(cmd.SrcNodeID); err == nil {
				go node.receiveApplicationCommand(cmd)
			} else {
				a.logger.Println("warn: Received command for unknown node", cmd.SrcNodeID)
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
				a.logger.Println("debug: controller update:", spew.Sdump(update))
			}

		default:
			a.logger.Println("debug: controller update:", spew.Sdump(update))

		}

	}
}

func (a *Layer) SendData(dstNode byte, payload commandclass.Command) error {
	marshaled, err := payload.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = a.serialAPI.SendData(dstNode, marshaled)
	return err
}

// SendDataSecure encapsulates payload in a security encapsulation command and
// sends it to the destination node.
func (a *Layer) SendDataSecure(dstNode byte, message commandclass.Command) error {
	// This function wraps the private sendDataSecure because no external packages
	// should ever call this while in inclusion mode (and doing so would be incorrect)
	return a.sendDataSecure(dstNode, message, false)
}

func (a *Layer) requestNonceForNode(dstNode byte) (security.Nonce, error) {
	err := a.SendData(dstNode, &zwsec.NonceGet{})

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

		a.logger.Printf("error: get nonce attempt #%d failed\n", i)
		time.Sleep(50 * time.Millisecond)
	}

	return nonce, err
}

func (a *Layer) sendDataSecure(dstNode byte, message commandclass.Command, inclusionMode bool) error {
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

	encapsulatedMessage, err := a.securityLayer.EncapsulateMessage(
		1,
		dstNode,
		zwsec.CommandMessageEncapsulation, // @todo CC should be determined by sequencing
		senderNonce,
		receiverNonce,
		securePayload,
		inclusionMode,
	)

	if err != nil {
		a.logger.Printf("error: failed to encrypt message: %v node=%d", err, dstNode)
		return err
	}

	return a.SendData(dstNode, encapsulatedMessage)
}

func (a *Layer) includeSecureNode(node *Node) error {
	a.secureInclusionStep[node.NodeID] = make(chan error)
	a.SendData(node.NodeID, &zwsec.SchemeGet{})

	defer close(a.secureInclusionStep[node.NodeID])
	defer delete(a.secureInclusionStep, node.NodeID)

	a.logger.Print("info: requesting security scheme")

	select {
	case err := <-a.secureInclusionStep[node.NodeID]:
		if err != nil {
			return err
		}
	case <-time.After(time.Second * 10):
		return errors.New("Secure inclusion timeout")
	}

	a.logger.Print("info: sending network key")
	node.NetworkKeySent = true

	a.sendDataSecure(
		node.NodeID,
		&zwsec.NetworkKeySet{NetworkKeyByte: a.networkKey},
		true,
	)

	select {
	case err := <-a.secureInclusionStep[node.NodeID]:
		return err
	case <-time.After(time.Second * 20):
		return errors.New("Secure inclusion timeout")
	}
}

func (a *Layer) interceptSecurityCommandClass(cmd serialapi.ApplicationCommand) {
	command, err := commandclass.Parse(1, cmd.CommandData)
	if err != nil {
		a.logger.Printf("error: %v\n", err)
		return
	}

	switch command.(type) {

	case zwsec.MessageEncapsulation, zwsec.MessageEncapsulationNonceGet:
		a.logger.Printf("rx secure message node=%d", cmd.SrcNodeID)
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
		var err error
		node, err := a.Node(cmd.SrcNodeID)
		if err != nil {
			a.logger.Printf("error: unknown node node=%X", cmd.SrcNodeID)
			return
		}

		if !node.NetworkKeySent {
			decrypted, err = a.securityLayer.DecryptMessage(cmd, true)
		} else {
			decrypted, err = a.securityLayer.DecryptMessage(cmd, false)
		}

		if err != nil {
			a.logger.Printf("error: error handling encrypted message %v\n", err)
			return
		}

		a.logger.Printf("info: received encapsulated message %s", spew.Sdump(decrypted))

		if decrypted[1] == byte(commandclass.Security) &&
			decrypted[2] == byte(zwsec.CommandNetworkKeyVerify) {
			a.logger.Printf("network key verify node=%d", cmd.SrcNodeID)
			if ch, ok := a.secureInclusionStep[cmd.SrcNodeID]; ok {
				ch <- nil
			}
			return
		}

		if node, ok := a.nodes[cmd.SrcNodeID]; ok {
			cmd.CommandData = decrypted[1:]
			go node.receiveApplicationCommand(cmd)
		} else {
			a.logger.Println("warn: received secure command for unknown node", cmd.SrcNodeID)
		}

	case zwsec.NonceGet:
		a.logger.Printf("info: nonce get node=%d", cmd.SrcNodeID)
		nonce, err := a.securityLayer.GenerateInternalNonce()
		if err != nil {
			a.logger.Println("alert: error generating internal nonce", err)
		}

		reply := &zwsec.NonceReport{NonceByte: nonce}
		a.SendData(cmd.SrcNodeID, reply)

	case zwsec.NonceReport:
		a.logger.Printf("nonce report node=%d", cmd.SrcNodeID)
		a.securityLayer.ReceiveNonce(cmd.SrcNodeID, (command.(zwsec.NonceReport)))

	case zwsec.SchemeReport:
		a.logger.Printf("security scheme report node=%d", cmd.SrcNodeID)
		if ch, ok := a.secureInclusionStep[cmd.SrcNodeID]; ok {
			ch <- nil
		} else {
			a.logger.Printf("warn: not in secure inclusion mode node=%d", cmd.SrcNodeID)
		}

	default:
		a.logger.Println("warn: unexpected security command:", spew.Sdump(cmd))
	}
}

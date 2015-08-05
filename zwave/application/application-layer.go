package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/security"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
)

// @todo should probably be 60 seconds, but 30 is easier for dev
const MaxSecureInclusionDuration = time.Second * 30

// Note: always use the smallest size based on the options
// @todo: also, implement the ability to use different options
const (
	SecurePayloadMaxSizeExplore   = 26 // in bytes
	SecurePayloadMaxSizeAutoRoute = 28
	SecurePayloadMaxSizeNoRoute   = 34
)

const (
	SecuritySequenceSequencedFlag   byte = 0x10
	SecuritySequenceSecondFrameFlag      = 0x20
	SecuritySequenceCounterMask          = 0x0f
)

type ApplicationLayer struct {
	ApiVersion     string
	ApiLibraryType string

	HomeId uint32
	NodeId byte

	Version             byte
	ApiType             string
	IsPrimaryController bool
	ApplicationVersion  byte
	ApplicationRevision byte
	SupportedFunctions  []byte

	NodeList []byte

	serialApi     serialapi.ISerialAPILayer
	securityLayer security.ISecurityLayer
	networkKey    []byte
	nodes         map[byte]*Node

	db *bolt.DB

	// maps node id to channel
	secureInclusionStep map[byte]chan error
}

func NewApplicationLayer(serialApi serialapi.ISerialAPILayer) (app *ApplicationLayer, err error) {
	app = &ApplicationLayer{
		serialApi: serialApi,
		nodes:     map[byte]*Node{},

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
	app.securityLayer = security.NewSecurityLayer(networkKey)

	go app.handleApplicationCommands()
	go app.handleControllerUpdates()

	err = app.initZWave()

	return
}

func (a *ApplicationLayer) Nodes() map[byte]*Node {
	return a.nodes
}

func (a *ApplicationLayer) Node(nodeId byte) (*Node, error) {
	if node, ok := a.nodes[nodeId]; ok {
		return node, nil
	}

	return nil, errors.New("Node not found")
}

func (a *ApplicationLayer) initDb() (err error) {
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

func (a *ApplicationLayer) initNetworkKey() ([]byte, error) {
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

func (a *ApplicationLayer) initZWave() error {
	version, err := a.serialApi.GetVersion()
	if err != nil {
		return err
	}

	a.ApiVersion = version.Version
	a.ApiLibraryType = version.GetLibraryTypeString()

	a.HomeId, a.NodeId, err = a.serialApi.MemoryGetId()
	if err != nil {
		return err
	}

	serialApiCapabilities, err := a.serialApi.GetSerialApiCapabilities()
	if err != nil {
		return err
	}

	a.ApplicationVersion = serialApiCapabilities.ApplicationVersion
	a.ApplicationRevision = serialApiCapabilities.ApplicationRevision
	a.SupportedFunctions = serialApiCapabilities.GetSupportedFunctions()

	initData, err := a.serialApi.GetInitAppData()
	if err != nil {
		return err
	}

	a.Version = initData.Version
	a.ApiType = initData.GetApiType()
	a.IsPrimaryController = initData.IsPrimaryController()
	a.NodeList = initData.GetNodeIds()

	for _, nodeId := range a.NodeList {
		node, err := NewNode(a, nodeId)

		if err != nil {
			spew.Dump(err)
			continue
		}

		a.nodes[nodeId] = node
	}

	return nil
}

func (a *ApplicationLayer) Shutdown() error {
	return a.db.Close()
}

func (a *ApplicationLayer) AddNode() (*Node, error) {
	newNodeInfo, err := a.serialApi.AddNode()
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
	a.nodes[node.NodeId] = node

	if node.IsSecure() {
		fmt.Println("Starting secure inclusion")
		err = a.includeSecureNode(node.NodeId)
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

func (a *ApplicationLayer) RemoveNode() (byte, error) {
	result, err := a.serialApi.RemoveNode()

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

func (a *ApplicationLayer) RemoveFailedNode(nodeId byte) (ok bool, err error) {
	ok, err = a.serialApi.RemoveFailedNode(nodeId)

	if ok && err != nil {
		err = a.db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("nodes")).Delete([]byte{nodeId})
		})
	}

	return
}

func (a *ApplicationLayer) handleApplicationCommands() {
	for cmd := range a.serialApi.ControllerCommands() {
		switch cmd.CommandData[0] {

		case commandclass.CommandClassSecurity:
			a.handleSecurityCommand(cmd)

		default:
			if node, err := a.Node(cmd.SrcNodeId); err == nil {
				go node.receiveApplicationCommand(cmd)
			} else {
				fmt.Println("Received command for unknown node", cmd.SrcNodeId)
			}

		}

	}
}

func (a *ApplicationLayer) handleControllerUpdates() {
	for update := range a.serialApi.ControllerUpdates() {

		switch update.Status {

		case protocol.UpdateStateNodeInfoReceived,
			protocol.UpdateStateNodeInfoReqFailed:
			if node, ok := a.nodes[update.NodeId]; ok {
				node.receiveControllerUpdate(update)
			} else {
				fmt.Println("controller update:", spew.Sdump(update))
			}

		default:
			fmt.Println("controller update:", spew.Sdump(update))

		}

	}
}

func (a *ApplicationLayer) SendData(dstNode byte, payload []byte) error {
	_, err := a.serialApi.SendData(dstNode, payload)
	return err
}

// SendDataSecure encapsulates payload in a security encapsulation command and
// sends it to the destination node.
func (a *ApplicationLayer) SendDataSecure(dstNode byte, payload []byte) error {
	// This function wraps the private sendDataSecure because no external packages
	// should ever call this while in inclusion mode (and doing so would be incorrect)
	return a.sendDataSecure(dstNode, payload, false)
}

func (a *ApplicationLayer) requestNonceForNode(dstNode byte) (security.Nonce, error) {
	err := a.SendData(dstNode, commandclass.NewSecurityNonceGet())

	if err != nil {
		return nil, err
	}

	return a.securityLayer.WaitForExternalNonce(dstNode)
}

func (a *ApplicationLayer) getOrRequestNonceForNode(dstNode byte) (nonce security.Nonce, err error) {
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

func (a *ApplicationLayer) sendDataSecure(dstNode byte, payload []byte, inclusionMode bool) error {
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

	var securityByte byte = 0
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

func (a *ApplicationLayer) includeSecureNode(nodeId byte) error {
	a.secureInclusionStep[nodeId] = make(chan error)
	a.SendData(nodeId, commandclass.NewSecuritySchemeGet())

	defer close(a.secureInclusionStep[nodeId])
	defer delete(a.secureInclusionStep, nodeId)

	select {
	case err := <-a.secureInclusionStep[nodeId]:
		if err != nil {
			return err
		}
	case <-time.After(time.Second * 10):
		return errors.New("Secure inclusion timeout")
	}

	a.sendDataSecure(
		nodeId,
		commandclass.NewSecurityNetworkKeySet(a.networkKey),
		true,
	)

	select {
	case err := <-a.secureInclusionStep[nodeId]:
		return err
	case <-time.After(time.Second * 20):
		return errors.New("Secure inclusion timeout")
	}
}

func (a *ApplicationLayer) handleSecurityCommand(cmd serialapi.ApplicationCommand) {
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
			if ch, ok := a.secureInclusionStep[cmd.SrcNodeId]; ok {
				ch <- nil
			}
			return
		}

		if node, ok := a.nodes[cmd.SrcNodeId]; ok {
			cmd.CommandData = msg
			go node.receiveApplicationCommand(cmd)
		} else {
			fmt.Println("Received secure command for unknown node", cmd.SrcNodeId)
		}

	case commandclass.CommandSecurityNonceGet:
		nonce, err := a.securityLayer.GenerateInternalNonce()
		if err != nil {
			fmt.Println("error generating internal nonce", err)
		}

		reply := commandclass.NewSecurityNonceReport(nonce)
		a.SendData(cmd.SrcNodeId, reply)

	case commandclass.CommandSecurityNonceReport:
		nonceReport := commandclass.ParseSecurityNonceReport(cmd.CommandData)
		a.securityLayer.ReceiveNonce(cmd.SrcNodeId, nonceReport)

	case commandclass.CommandSecuritySchemeReport:
		if ch, ok := a.secureInclusionStep[cmd.SrcNodeId]; ok {
			ch <- nil
		}

	default:
		fmt.Println("Unexpected security command:", spew.Sdump(cmd))
	}
}

package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/security"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
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
	SecuritySequenceSequencedFlag   uint8 = 0x10
	SecuritySequenceSecondFrameFlag       = 0x20
	SecuritySequenceCounterMask           = 0x0f
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

	NodeList []uint8

	serialApi     serialapi.ISerialAPILayer
	securityLayer security.ISecurityLayer
	Nodes         map[uint8]*Node

	// maps node id to channel
	secureInclusionStep map[uint8]chan error
}

func NewApplicationLayer(serialApi serialapi.ISerialAPILayer) (app *ApplicationLayer, err error) {
	app = &ApplicationLayer{
		serialApi:     serialApi,
		securityLayer: security.NewSecurityLayer(),
		Nodes:         map[uint8]*Node{},

		secureInclusionStep: map[uint8]chan error{},
	}

	go app.handleApplicationCommands()
	go app.handleControllerUpdates()

	err = app.initialize()

	return
}

func (a *ApplicationLayer) initialize() error {
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
		a.Nodes[nodeId] = NewNode(a, nodeId)
		a.Nodes[nodeId].initialize()
	}

	return nil
}

func (a *ApplicationLayer) AddNode() error {
	newNodeInfo, err := a.serialApi.AddNode()
	if err != nil {
		return err
	}

	node := NewNode(a, newNodeInfo.Source)
	node.setFromAddNodeCallback(newNodeInfo)
	a.Nodes[node.NodeId] = node

	if !node.IsSecure() {
		return nil
	}

	fmt.Println("Starting secure inclusion")
	return a.includeSecureNode(node.NodeId)
}

func (a *ApplicationLayer) RemoveNode() error {
	_, err := a.serialApi.RemoveNode()
	return err
}

func (a *ApplicationLayer) RemoveFailedNode(nodeId byte) (bool, error) {
	return a.serialApi.RemoveFailedNode(nodeId)
}

func (a *ApplicationLayer) handleApplicationCommands() {
	for cmd := range a.serialApi.ControllerCommands() {
		fmt.Println("app command loop")

		switch cmd.CommandData[0] {
		case commandclass.CommandClassSecurity:
			a.handleSecurityCommand(cmd)
		default:
			// @todo emit command on channel
			fmt.Println("application command:", spew.Sdump(cmd))
		}

	}
}

func (a *ApplicationLayer) handleControllerUpdates() {
	for update := range a.serialApi.ControllerUpdates() {

		switch update.Status {

		case protocol.UpdateStateNodeInfoReceived,
			protocol.UpdateStateNodeInfoReqFailed:
			if node, ok := a.Nodes[update.NodeId]; ok {
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

func (a *ApplicationLayer) sendDataSecure(dstNode byte, payload []byte, inclusionMode bool) error {
	// Previously, this function would just split and prepare the payload based on
	// whether it should be split after figuring out whether to segment. For now,
	// we're just going to assume that we will never have to worry about segmenting.
	// It wasn't too hard to implement before, but since I couldn't find a real payload
	// big enough, it wasn't possible to verify the implementation, so I didn't port
	// it while refactoring (for simplicity's sake).

	// Get a nonce from the other node
	receiverNonce, err := a.securityLayer.GetExternalNonce(dstNode)
	if err != nil {
		fmt.Println("sending get nonce")
		a.SendData(dstNode, commandclass.NewSecurityNonceGet())
		receiverNonce, err = a.securityLayer.WaitForExternalNonce(dstNode)

		if err != nil {
			fmt.Println("error getting nonce", err)
			return err
		}
		fmt.Println("got nonce")
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
		commandclass.NewSecurityNetworkKeySet(security.NetworkKey), // @todo
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
			fmt.Println("secure inclusion: network key verify", spew.Sdump(msg))
			if ch, ok := a.secureInclusionStep[cmd.SrcNodeId]; ok {
				ch <- nil
			}
			return
		}

		// cmd.CommandData = msg
		// cc = cmd.CommandData[0]
		// @todo emit decrypted command on channel

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

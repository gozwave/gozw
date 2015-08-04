package application

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
)

type CommandClassSupport int

const (
	CommandClassNotSupported CommandClassSupport = iota
	CommandClassSupportedInsecure
	CommandClassSupportedSecure
)

type Node struct {
	NodeId byte

	Capability          byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte

	Failing bool

	SupportedCommandClasses        map[byte]bool
	SecureSupportedCommandClasses  map[byte]bool
	SecureControlledCommandClasses map[byte]bool

	CommandClassVersions map[byte]byte

	DoorLock   *DoorLock
	Thermostat *Thermostat

	ManufacturerID uint16
	ProductTypeID  uint16
	ProductID      uint16

	application          *ApplicationLayer
	receivedUpdate       chan bool
	receivedSecurityInfo chan bool
}

func NewNode(application *ApplicationLayer, nodeId byte) (*Node, error) {
	node := &Node{
		NodeId: nodeId,

		SupportedCommandClasses:        map[byte]bool{},
		SecureSupportedCommandClasses:  map[byte]bool{},
		SecureControlledCommandClasses: map[byte]bool{},

		CommandClassVersions: map[byte]byte{},

		application:          application,
		receivedUpdate:       make(chan bool),
		receivedSecurityInfo: make(chan bool),
	}

	err := node.loadFromDb()
	if err != nil {
		initErr := node.initialize()
		if initErr != nil {
			return nil, initErr
		}

		node.saveToDb()
	}

	if IsDoorLock(node) {
		if node.DoorLock != nil {
			node.DoorLock.initialize(node)
		} else {
			node.DoorLock = NewDoorLock(node)
		}
	}

	if IsThermostat(node) {
		if node.Thermostat != nil {
			node.Thermostat.initialize(node)
		} else {
			node.Thermostat = NewThermostat(node)
		}
	}

	return node, nil
}

func (n *Node) loadFromDb() error {
	var data []byte
	err := n.application.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		data = bucket.Get([]byte{n.NodeId})

		if len(data) == 0 {
			return errors.New("Node not found")
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = msgpack.Unmarshal(data, n)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) initialize() error {
	nodeInfo, err := n.application.serialApi.GetNodeProtocolInfo(n.NodeId)
	if err != nil {
		fmt.Println(err)
	} else {
		n.setFromNodeProtocolInfo(nodeInfo)
	}

	if n.NodeId == 1 {
		// self is never failing
		n.Failing = false
	} else {
		failing, err := n.application.serialApi.IsFailedNode(n.NodeId)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		n.Failing = failing
	}

	if IsDoorLock(n) {
		n.DoorLock = NewDoorLock(n)
	}

	return n.saveToDb()
}

func (n *Node) saveToDb() error {
	data, err := msgpack.Marshal(n)
	if err != nil {
		return err
	}

	return n.application.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("nodes"))
		return bucket.Put([]byte{n.NodeId}, data)
	})
}

func (n *Node) IsSecure() bool {
	_, found := n.SupportedCommandClasses[commandclass.CommandClassSecurity]
	return found
}

func (n *Node) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *Node) GetBasicDeviceClassName() string {
	return protocol.GetBasicDeviceTypeName(n.BasicDeviceClass)
}

func (n *Node) GetGenericDeviceClassName() string {
	return protocol.GetGenericDeviceTypeName(n.GenericDeviceClass)
}

func (n *Node) GetSpecificDeviceClassName() string {
	return protocol.GetSpecificDeviceTypeName(n.GenericDeviceClass, n.SpecificDeviceClass)
}

func (n *Node) GetDoorLock() (*DoorLock, error) {
	if !IsDoorLock(n) {
		return nil, errors.New("Node is not designated as a door lock")
	}

	if n.DoorLock == nil {
		n.DoorLock = NewDoorLock(n)
	}

	return n.DoorLock, nil
}

func (n *Node) GetThermostat() (*Thermostat, error) {
	if !IsThermostat(n) {
		return nil, errors.New("Node is not designated as a thermostat")
	}

	if n.Thermostat == nil {
		n.Thermostat = NewThermostat(n)
	}

	return n.Thermostat, nil
}

func (n *Node) SendCommand(commandClass byte, command byte, commandPayload ...byte) error {
	supportType := n.SupportsCommandClass(commandClass)

	switch supportType {
	case CommandClassSupportedSecure:
		return n.sendDataSecure(append([]byte{commandClass, command}, commandPayload...))
	case CommandClassSupportedInsecure:
		return n.sendData(append([]byte{commandClass, command}, commandPayload...))
	case CommandClassNotSupported:
		return errors.New("Command class not supported")
	default:
		return errors.New("Command class not supported")
	}
}

func (n *Node) SupportsCommandClass(commandClass byte) CommandClassSupport {
	if supported, ok := n.SupportedCommandClasses[commandClass]; ok && supported {
		return CommandClassSupportedInsecure
	}

	if supported, ok := n.SecureSupportedCommandClasses[commandClass]; ok && supported {
		return CommandClassSupportedSecure
	}

	return CommandClassNotSupported
}

func (n *Node) AddAssociation(groupId byte, nodeIds ...byte) error {
	// sort of an arbitrary limit for now, but I'm not sure what it should be
	if len(nodeIds) > 20 {
		return errors.New("Too many associated nodes")
	}

	fmt.Println("Associating")

	return n.SendCommand(
		commandclass.CommandClassAssociation,
		commandclass.AssociationSet,
		append([]byte{groupId}, nodeIds...)...,
	)
}

func (n *Node) RequestSupportedSecurityCommands() error {
	return n.sendDataSecure([]byte{
		commandclass.CommandClassSecurity,
		commandclass.CommandSecurityCommandsSupportedGet,
	})
}

func (n *Node) RequestNodeInformationFrame() error {
	return n.application.serialApi.RequestNodeInfo(n.NodeId)
}

func (n *Node) LoadCommandClassVersions() error {
	for cc, _ := range n.SupportedCommandClasses {
		time.Sleep(1 * time.Second)
		err := n.sendData([]byte{
			commandclass.CommandClassVersion,
			commandclass.CommandVersionCommandClassGet,
			cc,
		})

		if err != nil {
			return err
		}
	}

	for cc, _ := range n.SecureSupportedCommandClasses {
		err := n.sendDataSecure([]byte{
			commandclass.CommandClassVersion,
			commandclass.CommandVersionCommandClassGet,
			cc,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) sendData(payload []byte) error {
	return n.application.SendData(n.NodeId, payload)
}

func (n *Node) sendDataSecure(payload []byte) error {
	return n.application.SendDataSecure(n.NodeId, payload)
}

func (n *Node) receiveControllerUpdate(update serialapi.ControllerUpdate) {
	select {
	case n.receivedUpdate <- true:
	default:
	}

	n.setFromApplicationControllerUpdate(update)
	n.saveToDb()
}

// func (n *Node) updateSupportedSecureCommands() {
// if n.IsSecure() {
// 	n.manager.SendDataSecure(n.NodeId, []byte{
// 		commandclass.CommandClassSecurity,
// 		commandclass.CommandSecurityCommandsSupportedGet,
// 	})
// } else {
// 	n.receivedSecure <- true
// }
// }

// func (n *Node) sendNoOp() {
// 	n.manager.session.SendData(n.NodeId, []byte{
// 		commandclass.CommandClassNoOperation,
// 	})
// }

// func (n *Node) IsFailing() bool {
// 	result, err := n.manager.session.isNodeFailing(n.NodeId)
// 	if err != nil {
// 		fmt.Println("node.isFailing error:", err)
// 	}
//
// 	return result
// }

func (n *Node) setFromAddNodeCallback(nodeInfo *serialapi.AddRemoveNodeCallback) {
	n.NodeId = nodeInfo.Source
	n.BasicDeviceClass = nodeInfo.Basic
	n.GenericDeviceClass = nodeInfo.Generic
	n.SpecificDeviceClass = nodeInfo.Specific

	for _, cc := range nodeInfo.CommandClasses {
		n.SupportedCommandClasses[cc] = true
	}

	n.saveToDb()
}

func (n *Node) setFromApplicationControllerUpdate(nodeInfo serialapi.ControllerUpdate) {
	n.BasicDeviceClass = nodeInfo.Basic
	n.GenericDeviceClass = nodeInfo.Generic
	n.SpecificDeviceClass = nodeInfo.Specific

	for _, cc := range nodeInfo.CommandClasses {
		n.SupportedCommandClasses[cc] = true
	}

	n.saveToDb()
}

func (n *Node) setFromNodeProtocolInfo(nodeInfo *serialapi.NodeProtocolInfo) {
	n.Capability = nodeInfo.Capability
	n.BasicDeviceClass = nodeInfo.BasicDeviceClass
	n.GenericDeviceClass = nodeInfo.GenericDeviceClass
	n.SpecificDeviceClass = nodeInfo.SpecificDeviceClass

	n.saveToDb()
}

func (n *Node) receiveSecurityCommandsSupportedReport(cc *commandclass.SecurityCommandsSupportedReport) {
	for _, cc := range cc.SupportedCommandClasses {
		n.SecureSupportedCommandClasses[cc] = true
	}

	for _, cc := range cc.ControlledCommandClasses {
		n.SecureControlledCommandClasses[cc] = true
	}

	select {
	case n.receivedSecurityInfo <- true:
	default:
	}

	n.saveToDb()
}

func (n *Node) receiveApplicationCommand(cmd serialapi.ApplicationCommand) {
	switch cmd.CommandData[0] {
	case commandclass.CommandClassSecurity:
		switch cmd.CommandData[1] {
		case commandclass.CommandSecurityCommandsSupportedReport:
			fmt.Println("security commands supported report")

			n.receiveSecurityCommandsSupportedReport(
				commandclass.ParseSecurityCommandsSupportedReport(cmd.CommandData),
			)

			fmt.Println(n.GetSupportedSecureCommandClassStrings())
		}

	case commandclass.CommandClassAlarm:
		if IsDoorLock(n) {
			lock, err := n.GetDoorLock()
			if err != nil {
				fmt.Println(err)
				return
			}

			lock.handleAlarmCommandClass(cmd)
		} else {
			fmt.Println("Alarm command for non-lock")
			spew.Dump(cmd)
		}

	case commandclass.CommandClassUserCode:
		lock, err := n.GetDoorLock()
		if err != nil {
			fmt.Println(err)
			return
		}

		lock.handleUserCodeCommandClass(cmd)

	case commandclass.CommandClassDoorLock:
		lock, err := n.GetDoorLock()
		if err != nil {
			fmt.Println(err)
			return
		}

		lock.handleDoorLockCommandClass(cmd)

	case commandclass.CommandClassThermostatMode:
		thermostat, err := n.GetThermostat()
		if err != nil {
			fmt.Println(err)
			return
		}

		thermostat.handleThermostatModeCommandClass(cmd)

	case commandclass.CommandClassThermostatSetpoint:
		thermostat, err := n.GetThermostat()
		if err != nil {
			fmt.Println(err)
			return
		}

		thermostat.handleThermostatSetpointCommandClass(cmd)

	case commandclass.CommandClassVersion:

		if cmd.CommandData[1] == commandclass.CommandVersionCommandClassReport {
			version := commandclass.ParseVersionCommandClassReport(cmd.CommandData)
			n.CommandClassVersions[version.CommandClass] = version.Version
		}

		n.saveToDb()

	default:
		fmt.Printf("unhandled application command (%d): %s\n", n.NodeId, spew.Sdump(cmd))
	}
}

func (n *Node) String() string {
	str := fmt.Sprintf("Node %d: \n", n.NodeId)
	str += fmt.Sprintf("  Failing? %t\n", n.Failing)
	str += fmt.Sprintf("  Is listening? %t\n", n.IsListening())
	str += fmt.Sprintf("  Is secure? %t\n", n.IsSecure())
	str += fmt.Sprintf("  Basic device class: %s\n", n.GetBasicDeviceClassName())
	str += fmt.Sprintf("  Generic device class: %s\n", n.GetGenericDeviceClassName())
	str += fmt.Sprintf("  Specific device class: %s\n", n.GetSpecificDeviceClassName())
	str += fmt.Sprintf("  Supported command classes:\n")
	for _, cc := range n.GetSupportedCommandClassStrings() {
		str += fmt.Sprintf("    - %s\n", cc)
	}

	if len(n.SecureSupportedCommandClasses) > 0 {
		secureCommands := commandClassSetToStrings(n.SecureSupportedCommandClasses)
		str += fmt.Sprintf("  Supported command classes (secure):\n")
		for _, cc := range secureCommands {
			str += fmt.Sprintf("    - %s\n", cc)
		}
	}

	if len(n.SecureControlledCommandClasses) > 0 {
		secureCommands := commandClassSetToStrings(n.SecureControlledCommandClasses)
		str += fmt.Sprintf("  Controlled command classes (secure):\n")
		for _, cc := range secureCommands {
			str += fmt.Sprintf("    - %s\n", cc)
		}
	}

	return str
}

func (n *Node) GetSupportedCommandClassStrings() []string {
	strings := commandClassSetToStrings(n.SupportedCommandClasses)
	if len(strings) == 0 {
		return []string{
			"None (probably not loaded; need to request a NIF)",
		}
	}

	return strings
}

func (n *Node) GetSupportedSecureCommandClassStrings() []string {
	strings := commandClassSetToStrings(n.SecureSupportedCommandClasses)
	return strings
}

func commandClassSetToStrings(commandClasses map[byte]bool) []string {
	if len(commandClasses) == 0 {
		return []string{}
	}

	ccStrings := []string{}

	for cc, _ := range commandClasses {
		ccStrings = append(ccStrings, commandclass.GetCommandClassString(cc))
	}

	return ccStrings
}

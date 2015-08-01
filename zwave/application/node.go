package application

import (
	"errors"
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/davecgh/go-spew/spew"
	set "github.com/deckarep/golang-set"
)

type Node struct {
	NodeId byte

	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte

	Failing bool

	SupportedCommandClasses        set.Set
	SecureSupportedCommandClasses  set.Set
	SecureControlledCommandClasses set.Set

	application          *ApplicationLayer
	receivedUpdate       chan bool
	receivedSecurityInfo chan bool
}

func NewNode(application *ApplicationLayer, nodeId byte) *Node {
	return &Node{
		NodeId: nodeId,

		SupportedCommandClasses:        set.NewSet(),
		SecureSupportedCommandClasses:  set.NewSet(),
		SecureControlledCommandClasses: set.NewSet(),

		application:          application,
		receivedUpdate:       make(chan bool),
		receivedSecurityInfo: make(chan bool),
	}
}

func (n *Node) IsSecure() bool {
	return n.SupportedCommandClasses.Contains(byte(commandclass.CommandClassSecurity))
}

func (n *Node) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *Node) HasOptionalFunctions() bool {
	return n.Security&0x80 == 0x80
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

func (n *Node) AddAssociation(groupId byte, nodeIds ...byte) error {
	// sort of an arbitrary limit for now, but I'm not sure what it should be
	if len(nodeIds) > 20 {
		return errors.New("Too many associated nodes")
	}

	fmt.Println("Associating")

	payload := append([]byte{
		commandclass.CommandClassAssociation,
		commandclass.AssociationSet,
		groupId,
	}, nodeIds...)

	return n.sendDataSecure(payload)
}

func (n *Node) RequestSupportedSecurityCommands() error {
	return n.sendDataSecure([]byte{
		commandclass.CommandClassSecurity,
		commandclass.CommandSecurityCommandsSupportedGet,
	})
}

func (n *Node) LoadUserCode(userId byte) error {
	return n.sendDataSecure([]byte{
		commandclass.CommandClassUserCode,
		commandclass.CommandUserCodeGet,
		userId,
	})
}

func (n *Node) sendData(payload []byte) error {
	return n.application.SendData(n.NodeId, payload)
}

func (n *Node) sendDataSecure(payload []byte) error {
	return n.application.SendDataSecure(n.NodeId, payload)
}

func (n *Node) initialize() {
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
			return
		}

		n.Failing = failing
	}

}

func (n *Node) receiveControllerUpdate(update serialapi.ControllerUpdate) {
	select {
	case n.receivedUpdate <- true:
	default:
	}

	n.setFromApplicationControllerUpdate(update)
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
		n.SupportedCommandClasses.Add(cc)
	}
}

func (n *Node) setFromApplicationControllerUpdate(nodeInfo serialapi.ControllerUpdate) {
	n.BasicDeviceClass = nodeInfo.Basic
	n.GenericDeviceClass = nodeInfo.Generic
	n.SpecificDeviceClass = nodeInfo.Specific

	for _, cc := range nodeInfo.CommandClasses {
		n.SupportedCommandClasses.Add(cc)
	}
}

func (n *Node) setFromNodeProtocolInfo(nodeInfo *serialapi.NodeProtocolInfo) {
	n.Capability = nodeInfo.Capability
	n.Security = nodeInfo.Security
	n.BasicDeviceClass = nodeInfo.BasicDeviceClass
	n.GenericDeviceClass = nodeInfo.GenericDeviceClass
	n.SpecificDeviceClass = nodeInfo.SpecificDeviceClass
}

func (n *Node) receiveSecurityCommandsSupportedReport(cc *commandclass.SecurityCommandsSupportedReport) {
	for _, cc := range cc.SupportedCommandClasses {
		n.SecureSupportedCommandClasses.Add(cc)
	}

	for _, cc := range cc.ControlledCommandClasses {
		n.SecureControlledCommandClasses.Add(cc)
	}

	n.receivedSecurityInfo <- true
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
		}

	case commandclass.CommandClassAlarm:
		// This is special handling code that will probably only work with yale locks
		notif := commandclass.ParseAlarmReport(cmd.CommandData)
		switch notif.Type {
		case 0x70:
			if notif.Level == 0x00 {
				fmt.Println("Master code changed")
			} else {
				fmt.Println("User added", notif.Level)
				n.LoadUserCode(notif.Level)
			}
		case 0xA1:
			if notif.Level == 0x01 {
				fmt.Println("Keypad limit exceeded")
			} else {
				fmt.Println("Physical tampering")
			}
		case 0x16:
			fmt.Println("Manual unlock")
		case 0x19:
			fmt.Println("RF operate unlock")
		case 0x15:
			fmt.Println("Manual lock")
		case 0x18:
			fmt.Println("RF operate lock")
		case 0x12:
			fmt.Println("keypad lock by user", notif.Level)
			n.LoadUserCode(notif.Level)
		case 0x13:
			fmt.Println("keypad unlock by user", notif.Level)
			n.LoadUserCode(notif.Level)
		case 0x09:
			fmt.Println("deadbolt jammed")
		case 0xA9:
			fmt.Println("dead battery; lock inoperable")
		case 0xA8:
			fmt.Println("critical battery")
		case 0xA7:
			fmt.Println("low battery")
		case 0x1B:
			fmt.Println("auto re-lock syscle completed")
		case 0x71:
			fmt.Println("duplicate pin code error")
		case 0x82:
			fmt.Println("power restored")
		case 0x21:
			fmt.Println("user deleted", notif.Level)
		}

	case commandclass.CommandClassUserCode:
		fmt.Println("user code")
		code := commandclass.ParseUserCodeReport(cmd.CommandData)
		spew.Dump(code)

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

	if n.SecureSupportedCommandClasses.Cardinality() > 0 {
		secureCommands := commandClassSetToStrings(n.SecureSupportedCommandClasses)
		str += fmt.Sprintf("  Supported command classes (secure):\n")
		for _, cc := range secureCommands {
			str += fmt.Sprintf("    - %s\n", cc)
		}
	}

	if n.SecureControlledCommandClasses.Cardinality() > 0 {
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

func commandClassSetToStrings(commandClasses set.Set) []string {
	if commandClasses.Cardinality() == 0 {
		return []string{}
	}

	ccStrings := []string{}

	for cc := range commandClasses.Iter() {
		ccStrings = append(ccStrings, commandclass.GetCommandClassString(cc.(byte)))
	}

	return ccStrings
}

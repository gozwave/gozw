package application

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
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
}

func NewNode(nodeId byte) *Node {
	return &Node{
		NodeId: nodeId,

		SupportedCommandClasses:        set.NewSet(),
		SecureSupportedCommandClasses:  set.NewSet(),
		SecureControlledCommandClasses: set.NewSet(),
	}
}

func NewNodeFromAddNodeCallback(callback *serialapi.AddRemoveNodeCallback) *Node {
	newNode := NewNode(callback.Source)
	newNode.setFromAddNodeCallback(callback)
	return newNode
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

func (n *Node) setFromApplicationControllerUpdate(nodeInfo *serialapi.ControllerUpdate) {
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

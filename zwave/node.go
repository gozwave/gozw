package zwave

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/commandclass"
)

type Node struct {
	NodeId byte

	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte

	SupportedCommandClasses []byte

	manager *Manager
}

func NewNode(manager *Manager, nodeId byte) *Node {
	return &Node{
		NodeId: nodeId,

		SupportedCommandClasses: []byte{},

		manager: manager,
	}
}

func (n *Node) IsSecure() bool {
	for _, cc := range n.SupportedCommandClasses {
		if cc == commandclass.CommandClassSecurity {
			return true
		}
	}

	return false
}

func (n *Node) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *Node) HasOptionalFunctions() bool {
	return n.Security&0x80 == 0x80
}

func (n *Node) GetBasicDeviceClassName() string {
	return GetBasicTypeName(n.BasicDeviceClass)
}

func (n *Node) GetGenericDeviceClassName() string {
	return GetGenericTypeName(n.GenericDeviceClass)
}

func (n *Node) GetSpecificDeviceClassName() string {
	return GetSpecificTypeName(n.GenericDeviceClass, n.SpecificDeviceClass)
}

func (n *Node) setFromAddNodeCallback(nodeInfo *AddRemoveNodeCallback) {
	n.NodeId = nodeInfo.Source
	n.BasicDeviceClass = nodeInfo.Basic
	n.GenericDeviceClass = nodeInfo.Generic
	n.SpecificDeviceClass = nodeInfo.Specific
	n.SupportedCommandClasses = nodeInfo.CommandClasses
}

func (n *Node) setFromNodeProtocolInfo(nodeInfo *NodeProtocolInfoResponse) {
	n.Capability = nodeInfo.Capability
	n.Security = nodeInfo.Security
	n.BasicDeviceClass = nodeInfo.BasicDeviceClass
	n.GenericDeviceClass = nodeInfo.GenericDeviceClass
	n.SpecificDeviceClass = nodeInfo.SpecificDeviceClass
}

func (n *Node) String() string {
	str := fmt.Sprintf("Node %d: \n", n.NodeId)
	str += fmt.Sprintf("  Is listening? %t\n", n.IsListening())
	str += fmt.Sprintf("  Basic device class: %s\n", n.GetBasicDeviceClassName())
	str += fmt.Sprintf("  Generic device class: %s\n", n.GetGenericDeviceClassName())
	str += fmt.Sprintf("  Specific device class: %s\n", n.GetSpecificDeviceClassName())
	str += fmt.Sprintf("  Supported command classes:\n")
	for _, cc := range n.GetSupportedCommandClassStrings() {
		str += fmt.Sprintf("    - %s\n", cc)
	}
	return str
}

func (n *Node) GetSupportedCommandClassStrings() []string {
	if len(n.SupportedCommandClasses) == 0 {
		return []string{
			"None (probably not loaded; need to request a NIF)",
		}
	}

	ccStrings := make([]string, len(n.SupportedCommandClasses))

	for i, cc := range n.SupportedCommandClasses {
		ccStrings[i] = commandclass.GetCommandClassString(cc)
	}

	return ccStrings
}

package zwave

import "github.com/bjyoungblood/gozw/zwave/commandclass"

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
		NodeId:  nodeId,
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

/* -------------------------------------------------------------------------- */
/* Copyright 2002-2019, OpenNebula Project, OpenNebula Systems                */
/*                                                                            */
/* Licensed under the Apache License, Version 2.0 (the "License"); you may    */
/* not use this file except in compliance with the License. You may obtain    */
/* a copy of the License at                                                   */
/*                                                                            */
/* http://www.apache.org/licenses/LICENSE-2.0                                 */
/*                                                                            */
/* Unless required by applicable law or agreed to in writing, software        */
/* distributed under the License is distributed on an "AS IS" BASIS,          */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   */
/* See the License for the specific language governing permissions and        */
/* limitations under the License.                                             */
/*--------------------------------------------------------------------------- */

package goca

import "strings"

// For virtual network template values: http://docs.opennebula.org/5.8/operation/references/vnet_template.html#
type VirtualNetworkTemplateKeys string

// Physical network template keys
const (
	VNMadVNK           VirtualNetworkTemplateKeys = "VN_MAD"
	BridgeVNK          VirtualNetworkTemplateKeys = "BRIDGE"
	VlanIDVNK          VirtualNetworkTemplateKeys = "VLAN_ID"
	AutomaticVlanIDVNK VirtualNetworkTemplateKeys = "AUTOMATIC_VLAN_ID"
	PhyDevVNK          VirtualNetworkTemplateKeys = "PHYDEV"
)

// Quality of service template keys
const (
	InboundAvgBwVNK   VirtualNetworkTemplateKeys = "INBOUND_AVG_BW"
	InboundPeakBwVNK  VirtualNetworkTemplateKeys = "INBOUND_PEAK_BW"
	InboundPeakKbVNK  VirtualNetworkTemplateKeys = "INBOUND_PEAK_KB"
	OutboundAvgBwVNK  VirtualNetworkTemplateKeys = "OUTBOUND_AVG_BW"
	OutboundPeakBwVNK VirtualNetworkTemplateKeys = "OUTBOUND_PEAK_BW"
	OutboundPeakKbVNK VirtualNetworkTemplateKeys = "OUTBOUND_PEAK_KB"
)

// Contextualization template keys
const (
	NetworkMaskVNK      VirtualNetworkTemplateKeys = "NETWORK_MASK"
	NetworkAddressVNK   VirtualNetworkTemplateKeys = "NETWORK_ADDRESS"
	GatewayVNK          VirtualNetworkTemplateKeys = "GATEWAY"
	Gateway6VNK         VirtualNetworkTemplateKeys = "GATEWAY6"
	DNSVNK              VirtualNetworkTemplateKeys = "DNS"
	GuestMTUVNK         VirtualNetworkTemplateKeys = "GUEST_MTU"
	ContextForceIPV4VNK VirtualNetworkTemplateKeys = "CONTEXT_FORCE_IPV4"
	SearchDomainVNK     VirtualNetworkTemplateKeys = "SEARCH_DOMAIN"
	SecGroupsVNK        VirtualNetworkTemplateKeys = "SECURITY_GROUPS"
)

// Interface creation options template keys
const (
	ConfVNK          VirtualNetworkTemplateKeys = "CONF"
	BridgeConfVNK    VirtualNetworkTemplateKeys = "BRIDGE"
	OvsBridgeConfVNK VirtualNetworkTemplateKeys = "OVS_BRIDGE_CONF"
	IPLinkConfVNK    VirtualNetworkTemplateKeys = "IP_LINK_CONF"
)

// VirtualNetworkTemplate is a virtual network template
type VirtualNetworkTemplate struct {
	VNMad   string             `xml:"VN_MAD"`
	ARs     []AddressRange     `xml:"AR"`
	Dynamic dynamicTemplateAny `xml:",any"`
}

func (vn *VirtualNetworkTemplate) String() string {
	var s strings.Builder

	for _, ar := range vn.ARs {
		s.WriteString(ar.String() + "\n")
	}
	s.WriteString(vn.Dynamic.String())

	return s.String()
}

func NewVirtualNetworkTemplate() *VirtualNetworkTemplate {
	return &VirtualNetworkTemplate{}
}

// Get retrieve a value from vnet template
func (t *VirtualNetworkTemplate) Get(key VirtualNetworkTemplateKeys) (string, error) {
	return t.Dynamic.GetStr(string(key))
}

// Add adds a vnet template key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VirtualNetworkTemplate) Add(key VirtualNetworkTemplateKeys, value string) error {
	return t.Dynamic.AddPair(string(key), value)
}

// AddVNMad allow to add a VN MAD to the template
func (t *VirtualNetworkTemplate) AddVNMad(value string) {
	t.VNMad = value
}

// GetVNMad allow to retrieve VN MAD from template
func (t *VirtualNetworkTemplate) GetVNMad() []AddressRange {
	return t.ARs
}

// AddAR allow to add a AR to the template
func (t *VirtualNetworkTemplate) AddAR(n *AddressRange) {
	t.ARs = append(t.ARs, *n)
}

// GetAR allow to retrieve ARs from template
func (t *VirtualNetworkTemplate) GetAR() []AddressRange {
	return t.ARs
}

// Address Range template part

// AddressRangeKeys is here to help the user to keep track of XML tags defined in AddressRange
type AddressRangeKeys string

const (
	arVecVNK        string           = "AR"
	aridVNK         string           = "AR_ID"
	IPARK           AddressRangeKeys = "IP"
	SizeARK         AddressRangeKeys = "SIZE"
	TypeARK         AddressRangeKeys = "TYPE"
	MacARK          AddressRangeKeys = "MAC"
	GlobalPrefixARK AddressRangeKeys = "GLOBAL_PREFIX"
	UlaPrefixARK    AddressRangeKeys = "ULA_PREFIX"
	PrefixLengthARK AddressRangeKeys = "PREFIX_LENGTH"
)

// AddressRange is a structure allowing to parse AddressRange templates. Common to VM and VRouter.
type AddressRange struct {
	TemplateVector
}

// NewAddressRange returns a structure disk entity to build
func NewAddressRange() *AddressRange {
	return &AddressRange{
		TemplateVector{key: arVecVNK},
	}
}

// GetID allow to retrieve the AddressRange ID doing the conversion to an unsigned integer
func (n *AddressRange) ID() int {
	id, _ := n.GetInt(string(aridVNK))
	return id
}

// Get allow to retrieve AddressRange keys
func (n *AddressRange) Get(key AddressRangeKeys) (string, error) {
	return n.GetStr(string(key))
}

// GetID convert and returns a value as an ID. The key name generally ends with an "_ID"
func (n *AddressRange) GetID(key AddressRangeKeys) (int, error) {
	return n.GetInt(string(key))
}

// Add adds an AddressRange key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *AddressRange) Add(key AddressRangeKeys, value string) error {
	return n.AddPair(string(key), value)
}

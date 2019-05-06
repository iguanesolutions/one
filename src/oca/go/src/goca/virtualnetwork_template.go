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

// For virtual network template values: http://docs.opennebula.org/5.8/operation/references/vnet_template.html#
type VNetTemplateKeys string

// Physical network template part
const (
	NameVNK            VNetTemplateKeys = "NAME"
	DescriptionVNK     VNetTemplateKeys = "DESCRIPTION"
	VNMadVNK           VNetTemplateKeys = "VN_MAD"
	BridgeVNK          VNetTemplateKeys = "BRIDGE"
	VlanIDVNK          VNetTemplateKeys = "VLAN_ID"
	AutomaticVlanIDVNK VNetTemplateKeys = "AUTOMATIC_VLAN_ID"
	PhyDevVNK          VNetTemplateKeys = "PHYDEV"
)

// Quality of service template part
const (
	InboundAvgBwVNK   VNetTemplateKeys = "INBOUND_AVG_BW"
	InboundPeakBwVNK  VNetTemplateKeys = "INBOUND_PEAK_BW"
	InboundPeakKbVNK  VNetTemplateKeys = "INBOUND_PEAK_KB"
	OutboundAvgBwVNK  VNetTemplateKeys = "OUTBOUND_AVG_BW"
	OutboundPeakBwVNK VNetTemplateKeys = "OUTBOUND_PEAK_BW"
	OutboundPeakKbVNK VNetTemplateKeys = "OUTBOUND_PEAK_KB"
)

// Contextualization template part
const (
	NetworkMaskVNK      VNetTemplateKeys = "NETWORK_MASK"
	NetworkAddressVNK   VNetTemplateKeys = "NETWORK_ADDRESS"
	GatewayVNK          VNetTemplateKeys = "GATEWAY"
	Gateway6VNK         VNetTemplateKeys = "GATEWAY6"
	DNSVNK              VNetTemplateKeys = "DNS"
	GuestMTUVNK         VNetTemplateKeys = "GUEST_MTU"
	ContextForceIPV4VNK VNetTemplateKeys = "CONTEXT_FORCE_IPV4"
	SearchDomainVNK     VNetTemplateKeys = "SEARCH_DOMAIN"
	SecGroupsVNK        VNetTemplateKeys = "SECURITY_GROUPS"
)

// Interface creation options template part
const (
	ConfVNK          VNetTemplateKeys = "CONF"
	BridgeConfVNK    VNetTemplateKeys = "BRIDGE"
	OvsBridgeConfVNK VNetTemplateKeys = "OVS_BRIDGE_CONF"
	IPLinkConfVNK    VNetTemplateKeys = "IP_LINK_CONF"
)

// VNetTemplate is a virtual network template
type VNetTemplate struct {
	VNMad   string          `xml:"VN_MAD"`
	ARs     []AddressRange  `xml:"AR"`
	Dynamic DynamicTemplate `xml:",any"`
}

func NewVNetTemplate() *VNetTemplate {
	return &VNetTemplate{}
}

// Get retrieve a value from vnet template
func (t *VNetTemplate) Get(key VNetTemplateKeys) (string, error) {
	return t.Dynamic.getValue(string(key))
}

// Add adds a vnet template key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VNetTemplate) Add(key VNetTemplateKeys, value string) error {
	return t.Dynamic.AddPair(string(key), value)
}

// AddVNMad allow to add a VN MAD to the template
func (t *VNetTemplate) AddVNMad(value string) {
	t.VNMad = value
}

// GetVNMad allow to retrieve VN MAD from template
func (t *VNetTemplate) GetVNMad() []AddressRange {
	return t.ARs
}

// AddAR allow to add a AR to the template
func (t *VNetTemplate) AddAR(n *AddressRange) {
	t.ARs = append(t.ARs, *n)
}

// GetAR allow to retrieve ARs from template
func (t *VNetTemplate) GetAR() []AddressRange {
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
	DynamicTemplateVector
}

// NewAddressRange returns a structure disk entity to build
func NewAddressRange() *AddressRange {
	return &AddressRange{
		DynamicTemplateVector{key: arVecVNK},
	}
}

// Get allow to retrieve AddressRange keys
func (n *AddressRange) Get(key AddressRangeKeys) (string, error) {
	return n.getValue(string(key))
}

// GetID allow to retrieve the AddressRange ID doing the conversion to an unsigned integer
func (n *AddressRange) GetID() (uint, error) {
	return n.getID(string(aridVNK))
}

// Add adds an AddressRange key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *AddressRange) Add(key AddressRangeKeys, value string) error {
	return n.AddPair(string(key), value)
}

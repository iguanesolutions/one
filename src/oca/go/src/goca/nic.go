package goca

// NIC is a structure allowing to parse NIC templates. Common to VM and VRouter.
type NIC struct {
	DynamicTemplateVector
}

// NicKeys is here to help the user to keep track of XML tags defined in NIC
type NicKeys string

// Some keys are specific to VM some others to VRouter
// For VM values: https://docs.opennebula.org/5.8/operation/references/template.html#network-section
const (
	vecNK               string  = "NIC"
	idNK                NicKeys = "NIC_ID"
	BridgeNK            NicKeys = "BRIDGE"
	FilterNK            NicKeys = "FILTER"
	IPNK                NicKeys = "IP"
	MACNK               NicKeys = "MAC"
	NetworkNK           NicKeys = "NETWORK"
	NetworkMaskNK       NicKeys = "NETWORK_MASK"
	NetworkIDNK         NicKeys = "NETWORK_ID"
	NetworkUIDNK        NicKeys = "NETWORK_UID"
	NetworkUNameNK      NicKeys = "NETWORK_UNAME"
	NetworkAddressNK    NicKeys = "NETWORK_ADDRESS"
	PhyDevNK            NicKeys = "PHYDEV"
	SecGroupsNK         NicKeys = "SECURITY_GROUPS"
	TargetNK            NicKeys = "TARGET"
	VlanIDNK            NicKeys = "VLAN_ID"
	ScriptNK            NicKeys = "SCRIPT"
	ModelNK             NicKeys = "MODEL"
	InboundAvgBwNK      NicKeys = "INBOUND_AVG_BW"
	InboundPeakBwNK     NicKeys = "INBOUND_PEAK_BW"
	InboundPeakKNK      NicKeys = "INBOUND_PEAK_KB"
	OutboundAvgBwNK     NicKeys = "OUTBOUND_AVG_BW"
	OutboundPeakBwNK    NicKeys = "OUTBOUND_PEAK_BW"
	OutboundPeakKbNK    NicKeys = "OUTBOUND_PEAK_KB"
	NetworkModeNK       NicKeys = "NETWORK_MODE"
	SchedRequirementsNK NicKeys = "SCHED_REQUIREMENTS"
	SchedRankNK         NicKeys = "SCHED_RANK"
	NameNK              NicKeys = "NAME"
	ParentNK            NicKeys = "PARENT"
	ExternalNK          NicKeys = "EXTERNAL"
)

// NewNIC returns a structure disk entity to build
func NewNIC() *NIC {
	return &NIC{
		DynamicTemplateVector{key: vecNK},
	}
}

// Get is a getValue for all NIC keys
func (n *NIC) Get(key NicKeys) (string, error) {
	return n.getValue(string(key))
}

// GetID is a getValue for the NIC ID doing the conversion to an unsigned integer
func (n *NIC) GetID() (uint, error) {
	return n.getID(string(idNK))
}

// Add adds a NIC key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *NIC) Add(key NicKeys, value string) error {
	return n.AddPair(string(key), value)
}

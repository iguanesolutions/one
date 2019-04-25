package goca

// NIC is a structure allowing to parse NIC templates. Common to VM and VRouter.
type NIC struct {
	TemplateVector
}

// NICKeys is here to help the user to keep track of XML tags defined in NIC
type NICKeys string

// Some keys are specific to VM some others to VRouter
// For VM values: https://docs.opennebula.org/5.8/operation/references/template.html#network-section
const (
	vecNK               string  = "NIC"
	idNK                NICKeys = "NIC_ID"
	BridgeNK            NICKeys = "BRIDGE"
	FilterNK            NICKeys = "FILTER" // Define network filtering rule for the interface
	IPNK                NICKeys = "IP"
	MACNK               NICKeys = "MAC"
	NetworkNK           NICKeys = "NETWORK"
	NetworkMaskNK       NICKeys = "NETWORK_MASK"
	NetworkIDNK         NICKeys = "NETWORK_ID"
	NetworkUIDNK        NICKeys = "NETWORK_UID"
	NetworkUNameNK      NICKeys = "NETWORK_UNAME"
	NetworkAddressNK    NICKeys = "NETWORK_ADDRESS"
	SecGroupsNK         NICKeys = "SECURITY_GROUPS" // List of security group to be applied
	TargetNK            NICKeys = "TARGET"
	VlanIDNK            NICKeys = "VLAN_ID"
	ScriptNK            NICKeys = "SCRIPT"
	ModelNK             NICKeys = "MODEL"
	InboundAvgBwNK      NICKeys = "INBOUND_AVG_BW"
	InboundPeakBwNK     NICKeys = "INBOUND_PEAK_BW"
	InboundPeakKNK      NICKeys = "INBOUND_PEAK_KB"
	OutboundAvgBwNK     NICKeys = "OUTBOUND_AVG_BW"
	OutboundPeakBwNK    NICKeys = "OUTBOUND_PEAK_BW"
	OutboundPeakKbNK    NICKeys = "OUTBOUND_PEAK_KB"
	NetworkModeNK       NICKeys = "NETWORK_MODE"
	SchedRequirementsNK NICKeys = "SCHED_REQUIREMENTS"
	SchedRankNK         NICKeys = "SCHED_RANK"
	NameNK              NICKeys = "NAME"
	ParentNK            NICKeys = "PARENT"
	ExternalNK          NICKeys = "EXTERNAL"
)

// NewNIC returns a structure disk entity to build
func NewNIC() *NIC {
	return &NIC{
		TemplateVector{key: vecNK},
	}
}

// ID returns the NIC ID
func (n *NIC) ID() int {
	id, _ := n.GetInt(string(idNK))
	return id
}

// Get return the string value of a NIC key
func (n *NIC) Get(key NICKeys) (string, error) {
	return n.GetStr(string(key))
}

// GetID convert and returns a value as an ID. The key name generally ends with an "_ID"
func (n *NIC) GetID(key NICKeys) (int, error) {
	return n.GetInt(string(key))
}

// Add adds a NIC key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *NIC) Add(key NICKeys, value string) error {
	return n.AddPair(string(key), value)
}

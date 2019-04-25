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
	NicIDK          NicKeys = "NIC_ID"
	ArIDK           NicKeys = "AR_ID"
	BridgeK         NicKeys = "BRIDGE"
	BridgeTypeK     NicKeys = "BRIDGE_TYPE"
	ClusterIDK      NicKeys = "CLUSTER_ID"
	FilterK         NicKeys = "FILTER"
	FloatingIPK     NicKeys = "FLOATING_IP"
	GatewayK        NicKeys = "GATEWAY"
	IPK             NicKeys = "IP"
	MACK            NicKeys = "MAC"
	MTUK            NicKeys = "MTU"
	NetworkK        NicKeys = "NETWORK"
	NetworkMaskK    NicKeys = "NETWORK_MASK"
	NetworkIDK      NicKeys = "NETWORK_ID"
	NetworkUIDK     NicKeys = "NETWORK_UID"
	NetworkUNameK   NicKeys = "NETWORK_UNAME"
	NetworkAddressK NicKeys = "NETWORK_ADDRESS"
	PhyDevK         NicKeys = "PHYDEV"
	PublicIPK       NicKeys = "PULIC_IP"
	SecGroupsK      NicKeys = "SECURITY_GROUPS"
	TargetK         NicKeys = "TARGET"
	VlanIDK         NicKeys = "VLAN_ID"
	VNMadK          NicKeys = "VN_MAD"
)

// NewNIC returns a structure disk entity to build
func NewNIC() NIC {
	return NIC{
		DynamicTemplateVector{key: "NIC"},
	}
}

// Get is a getValue for all NIC keys
func (n *NIC) Get(key NicKeys) (string, error) {
	return n.getValue(string(key))
}

// GetID is a getValue for the NIC ID doing the conversion to an unsigned integer
func (n *NIC) GetID() (uint, error) {
	return n.getID(string(NicIDK))
}

// Add adds a NIC key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *NIC) Add(key NicKeys, value string) error {
	return n.AddPair(string(key), value)
}

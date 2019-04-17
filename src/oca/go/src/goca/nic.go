package goca

import "strconv"

// NIC is a structure allowing to parse/build NIC templates. Common between VM and VRouter
type NIC struct {
	DynTemplateVector
}

// NicKeys is here to help the user to keep track of XML tags defined in NIC
type NicKeys string

// Some keys are specific to VM some others to VRouter
const (
	NicIDK          NicKeys = "NIC_ID"
	ArIDK           NicKeys = "AR_ID"
	BridgeK         NicKeys = "BRIDGE"
	BridgeTypeK     NicKeys = "BRIDGE_TYPE"
	ClusterIDK      NicKeys = "CLUSTER_ID"
	FloatingIPK     NicKeys = "FLOATING_IP"
	GatewayK        NicKeys = "GATEWAY"
	IPK             NicKeys = "IP"
	MACK            NicKeys = "MAC"
	MTUK            NicKeys = "MTU"
	NetworkK        NicKeys = "NETWORK"
	NetworkIDK      NicKeys = "NETWORK_ID"
	NetworkMaskK    NicKeys = "NETWORK_MASK"
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
	return NIC{DynTemplateVector{key: "NIC"}}
}

// Get is a getter for all NIC keys
func (n *NIC) Get(key NicKeys) (string, error) {
	return n.getterVec(string(key))
}

// Add adds a NIC key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (n *NIC) Add(key NicKeys, value string) error {
	return n.AddPair(string(key), value)
}

// GetID is a getter for the NIC ID doing the conversion to an unsigned integer
func (n *NIC) GetID() (uint, error) {
	idStr, err := n.getterVec(string(NicIDK))
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(id), err
}

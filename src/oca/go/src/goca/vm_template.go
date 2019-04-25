package goca

import (
	"encoding/base64"
	"fmt"
)

// User template

// VMUserTemplate contains common and custom attributes
type VMUserTemplate struct {
	Error        string          `xml:"ERROR"`
	SchedMessage string          `xml:"SCHED_MESSAGE"`
	Dynamic      DynamicTemplate `xml:",any"`
}

// Get allow to get a pair by key
func (d *VMUserTemplate) Get(key string) (string, error) {
	return d.Dynamic.getValue(string(key))
}

// VM Template
// Available template parts and keys are listed here: https://docs.opennebula.org/5.8/operation/references/template.html
// Some specific part are not defined: vCenter, Public Cloud, Hypervisor, User Inputs

// Name of the vector in the template, they should match with tags in vmDynamicTemplate struct
const (
	vmContextK        string = "CONTEXT"
	vmCPUModelVecK    string = "CPU_MODEL"
	vmFeaturesVecK    string = "FEATURES"
	vmIOGraphicsVecK  string = "GRAPHICS"
	vmIOInputVecK     string = "INPUT"
	vmOSVecK          string = "OS"
	vmSchedActionVecK string = "SCHED_ACTION"
)

// Common to VM and Template entities
type vmDynamicTemplate struct {
	// TODO: improve segmentation (better autocompletion), but add code
	// CPUModel    VMCPUModel  `xml:"CPU_MODEL"`
	// Graphics    VMGraphics  `xml:"GRAPHICS"`
	// Input       VMInput     `xml:"INPUT"`
	// OS          VMOS        `xml:"OS"`

	Context VMContext `xml:"CONTEXT"`

	// Defined in disk.go and nic.go
	Disks []Disk `xml:"DISK"`
	NICs  []NIC  `xml:"NIC"`

	Dynamic DynamicTemplate `xml:",any"`
}

// VMTemplate is a structure allowing to parse VM templates.
// It's defined in a semi-static way to guide the user among the bunch of values
type VMTemplate struct {
	vmDynamicTemplate

	// Capacity fields are only in the parsed part because they may be useless in builder
	// due to instantiate method. They could be already defined in the OpenNebula Template entity.
	CPU    float64 `xml:"CPU"`
	Memory uint    `xml:"MEMORY"`
	VCPU   uint    `xml:"VCPU"` // If not defined: it's 1

	// These part are only parsed
	Snapshots          []VMSnapshot          `xml:"SNAPSHOT"`
	SecurityGroupRules []VMSecurityGroupRule `xml:"SECURITY_GROUP_RULE"`
}

// VMSecurityGroupRule
type VMSecurityGroupRule struct {
	SecurityGroupRule
	SecurityGroup string `xml:"SECURITY_GROUP_NAME"`
}

// VMDynamicTemplate is an extension of VMTemplate allowing to build VM templates
// Note: when building, there is no check enforced on the consitency of your template on GOCA side
type VMTemplateBuilder struct {
	vmDynamicTemplate
}

func (t *vmDynamicTemplate) String() string {
	s := ""
	endToken := "\n"

	// TODO: add all fields ?
	// We need to add manually the strings as the field where manually defined in structure

	s += t.Context.String() + endToken
	for _, disk := range t.Disks {
		s += disk.String() + endToken
	}
	for _, nic := range t.NICs {
		s += nic.String() + endToken
	}

	// TODO: check if it's correct. Check the len for the end token value
	s += t.Dynamic.String()

	return s
}

// NewVMTemplate returns a VMTemplate structure
func NewVMDynamicTemplate() VMTemplateBuilder {
	return VMTemplateBuilder{}
}

// Template parts

// Capacity template part

type VMCapacityKeys string

const (
	CPUK    VMCapacityKeys = "CPU"
	VCPUK   VMCapacityKeys = "VCPU"
	MemoryK VMCapacityKeys = "MEMORY"
)

//
func (t *VMTemplateBuilder) SetCapacity(CPU float64, VCPU, Memory uint) error {
	for _, key := range []VMCapacityKeys{CPUK, VCPUK, MemoryK} {
		if !t.Dynamic.Exists(string(key)) {
			continue
		}
		return fmt.Errorf("VMDynamicTemplate.SetCapacity: the key %s is already present in template", "")
	}
	t.Dynamic.AddPair("CPU", fmt.Sprint(CPU))
	t.Dynamic.AddPair("Memory", fmt.Sprint(Memory))
	t.Dynamic.AddPair("VCPU", fmt.Sprint(VCPU))

	return nil
}

// Disk template part

// AddDisk allow to add a disk to the template
func (t *VMTemplateBuilder) AddDisk(d Disk) {
	t.Disks = append(t.Disks, d)
}

// GetDisk allow to retrieve disks from template
func (t *VMTemplate) GetDisk() []Disk {
	return t.Disks
}

// NIC template part

// AddNIC allow to add a NIC to the template
func (t *VMTemplateBuilder) AddNIC(n NIC) {
	t.NICs = append(t.NICs, n)
}

// GetNIC allow to retrieve NICs from template
func (t *VMTemplate) GetNIC() []NIC {
	return t.NICs
}

// Show back template part

// VMShowbackKeys define keys for showback values
type VMShowbackKeys string

const (
	MemCostK  VMShowbackKeys = "MEMORY_COST"
	CPUCostK  VMShowbackKeys = "CPU_COST"
	DiskCostK VMShowbackKeys = "DISK_COST"
)

func (d *VMTemplateBuilder) SetShowback(key VMShowbackKeys, value string) error {
	for _, key := range []VMShowbackKeys{MemCostK, CPUCostK, DiskCostK} {
		if !d.Dynamic.Exists(string(key)) {
			continue
		}
		return fmt.Errorf("VMDynamicTemplate.SetShowbackAttr: the key %s is already present in template", "")
	}
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetShowback(key VMShowbackKeys) (string, error) {
	return d.Dynamic.getValue(string(key))
}

// OS template part

// VMOSKeys define keys for OS and boot values
type VMOSKeys string

const (
	ArchK       VMOSKeys = "ARCH"
	MachineK    VMOSKeys = "MACHINE"
	KernelK     VMOSKeys = "KERNEL"
	KernelDSK   VMOSKeys = "KERNEL_DS"
	InitrdK     VMOSKeys = "INITRD"
	InitrdDSK   VMOSKeys = "INITRD_DS"
	RootK       VMOSKeys = "ROOT"
	KernelCmdK  VMOSKeys = "KERNEL_CMD"
	BootloaderK VMOSKeys = "BOOTLOADER"
	BootK       VMOSKeys = "BOOT"
)

func (d *VMTemplateBuilder) AddOS(key VMShowbackKeys, value string) error {
	return d.Dynamic.addPairToVec(vmOSVecK, string(key), value)
}

func (d *VMTemplate) GetOS(key VMOSKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmOSVecK), string(key))
}

// CPU model part

type VMCPUModelKeys string

const (
	ModelK VMCPUModelKeys = "MODEL"
)

// There is only one key defined for the CPU_MODEL vector, so we just define a getValue and a setter for it
func (d *VMTemplateBuilder) AddCPUModel(value string) error {
	return d.Dynamic.addPairToVec(vmCPUModelVecK, string(ModelK), value)
}

func (d *VMTemplate) GetCPUModel(key VMCPUModelKeys) (string, error) {
	cpuMod, err := d.Dynamic.GetVector(string(vmCPUModelVecK))
	if err != nil {
		return "", fmt.Errorf("VMTemplate.GetCPUModel: vector %s: %s", vmCPUModelVecK, err)
	}
	return cpuMod.getValue(string(key))
}

// Features template part

type VMFeatureKeys string

const (
	PAEK              VMFeatureKeys = "PAE"
	ACPIK             VMFeatureKeys = "ACPI"
	APICK             VMFeatureKeys = "APIC"
	LocalTimeK        VMFeatureKeys = "LOCAL_TIME"
	GuestAgentK       VMFeatureKeys = "GUEST_AGENT"
	VirtIOScsiQueuesK VMFeatureKeys = "VIRTIO_SCSI_QUEUES"
)

func (d *VMTemplateBuilder) AddFeature(key VMFeatureKeys, value string) error {
	return d.Dynamic.addPairToVec(vmFeaturesVecK, string(key), value)
}

func (d *VMTemplate) GetFeature(key VMFeatureKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmFeaturesVecK), string(key))
}

// I/O devices template part

type VMIOGraphicsKeys string
type VMIOInputKeys string

const (
	InputTypeK VMIOInputKeys = "TYPE" // Values: mouse or tablet
	BusK       VMIOInputKeys = "BUS"  // Values: usb or ps2

	GraphicTypeK    VMIOGraphicsKeys = "TYPE" // Values: vnc, sdl, spice
	ListenK         VMIOGraphicsKeys = "LISTEN"
	PortK           VMIOGraphicsKeys = "PORT"
	PasswdK         VMIOGraphicsKeys = "PASSWD"
	KeymapK         VMIOGraphicsKeys = "KEYMAP"
	RandomPasswordK VMIOGraphicsKeys = "RANDOM_PASSWD"
)

func (d *VMTemplateBuilder) AddIOGraphic(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmIOGraphicsVecK, string(key), value)
}

func (d *VMTemplate) GetIOGraphic(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmIOGraphicsVecK), string(key))
}

func (d *VMTemplateBuilder) AddIOInput(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmIOInputVecK, string(key), value)
}

func (d *VMTemplate) GetIOInput(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmIOInputVecK), string(key))
}

// Context template part

type VMContext struct {
	DynamicTemplateVector
}

// VMContextKeys is here to help the user to keep track of XML tags defined in VM context
type VMContextKeys string
type VMContextB64Keys string

const (
	DNSK         VMContextKeys = "DNS"
	DNSHostNameK VMContextKeys = "DNS_HOSTNAME"
	EC2PubKeyK   VMContextKeys = "EC2_PUBLIC_KEY"
	FilesK       VMContextKeys = "FILES"
	FilesDSK     VMContextKeys = "FILES_DS"
	GatewayIface VMContextKeys = "GATEWAY_IFACE"
	NetworkCtxK  VMContextKeys = "NETWORK"
	InitScriptsK VMContextKeys = "INIT_SCRIPTS"
	SSHPubKeyK   VMContextKeys = "SSH_PUBLIC_KEY"
	TargetCtxK   VMContextKeys = "TARGET"
	TokenK       VMContextKeys = "TOKEN"
	UsernameK    VMContextKeys = "USERNAME"
	VariableK    VMContextKeys = "VARIABLE"
	SecureTTYK   VMContextKeys = "SECURETTY"
	SetHostnameK VMContextKeys = "SET_HOSTNAME"

	// Keys for Base64 values
	PasswordB64K    VMContextB64Keys = "PASSWORD_BASE64"
	StartScriptB64K VMContextB64Keys = "START_SCRIPT_BASE64"
	CryptedPassB64K VMContextB64Keys = "CRYPTED_PASSWORD_BASE64"

	// Note: ETHx_XXX values not mapped.
)

// Get is a getValue for all context keys
func (t *VMTemplate) GetCtx(key VMContextKeys) (string, error) {
	return t.Context.getValue(string(key))
}

// Add adds a context key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplateBuilder) AddCtx(key VMContextKeys, value string) error {
	return t.Context.AddPair(string(key), value)
}

// Add adds a context key with value. It will convert value to base64. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplateBuilder) AddB64Ctx(key VMContextB64Keys, value string) error {
	valueB64 := base64.StdEncoding.EncodeToString([]byte(value))
	return t.Context.AddPair(string(key), valueB64)
}

// Placement Template part

type VMPlacementKeys string

const (
	SchedRequirementsK   VMPlacementKeys = "SCHED_REQUIREMENTS"
	SchedRankK           VMPlacementKeys = "SCHED_RANK"
	SchedDSRequirementsK VMPlacementKeys = "SCHED_DS_REQUIREMENTS"
	SchedDSRankK         VMPlacementKeys = "SCHED_DS_RANK"
	UserPriorityK        VMPlacementKeys = "USER_PRIORITY"
)

func (d *VMTemplateBuilder) AddFeatureAttr(key VMFeatureKeys, value string) error {
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetFeatureAttr(key VMFeatureKeys) (string, error) {
	return d.Dynamic.getValue(string(key))
}

// Scheduled actions template part
type VMSchedActionKeys string

const (
	TimeK     VMSchedActionKeys = "TIME"
	RepeatK   VMSchedActionKeys = "REPEAT"
	DaysK     VMSchedActionKeys = "DAYS"
	ActionK   VMSchedActionKeys = "ACTION"
	EndTypeK  VMSchedActionKeys = "END_TYPE"
	EndValueK VMSchedActionKeys = "END_VALUE"
)

func (d *VMTemplateBuilder) AddSchedAction(key VMSchedActionKeys, value string) error {
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetSchedAction(key VMSchedActionKeys) (string, error) {
	return d.Dynamic.getValue(string(key))
}

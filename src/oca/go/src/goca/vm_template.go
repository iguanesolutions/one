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
	vmkContextVec     string = "CONTEXT"
	vmkCPUModelVec    string = "CPU_MODEL"
	vmkFeaturesVec    string = "FEATURES"
	vmkIOGraphicsVec  string = "GRAPHICS"
	vmkIOInputVec     string = "INPUT"
	vmkOSVec          string = "OS"
	vmkSchedActionVec string = "SCHED_ACTION"
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
	Disks        []Disk        `xml:"DISK"`
	NICs         []NIC         `xml:"NIC"`
	SchedActions []SchedAction `xml:"SCHED_ACTION"`

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
	CPUVMK    VMCapacityKeys = "CPU"
	VCPUVMK   VMCapacityKeys = "VCPU"
	MemoryVMK VMCapacityKeys = "MEMORY"
)

// SetCapacity adds once a capacity attribute
func (t *VMTemplateBuilder) SetCapacity(CPU float64, VCPU, Memory uint) error {
	for _, key := range []VMCapacityKeys{CPUVMK, VCPUVMK, MemoryVMK} {
		if !t.Dynamic.Exists(string(key)) {
			continue
		}
		return fmt.Errorf("VMTemplateBuilder.SetCapacity: the key %s is already present in template", "")
	}
	t.Dynamic.AddPair("CPU", fmt.Sprint(CPU))
	t.Dynamic.AddPair("Memory", fmt.Sprint(Memory))
	t.Dynamic.AddPair("VCPU", fmt.Sprint(VCPU))

	return nil
}

// Disk template part

// AddDisk allow to add a disk to the template
func (t *VMTemplateBuilder) AddDisk(d *Disk) {
	t.Disks = append(t.Disks, *d)
}

// GetDisk allow to retrieve disks from template
func (t *VMTemplate) GetDisk() []Disk {
	return t.Disks
}

// NIC template part

// AddNIC allow to add a NIC to the template
func (t *VMTemplateBuilder) AddNIC(n *NIC) {
	t.NICs = append(t.NICs, *n)
}

// GetNIC allow to retrieve NICs from template
func (t *VMTemplate) GetNIC() []NIC {
	return t.NICs
}

// Sched action get/set

// AddSchedAction returns a structure disk entity to build
func (t *VMTemplate) AddSchedAction(a *SchedAction) {
	t.SchedActions = append(t.SchedActions, *a)
}

// GetSchedActions allow to retrieve SchedActions from template
func (t *VMTemplate) GetSchedActions() []SchedAction {
	return t.SchedActions
}

// Show back template part

// VMShowbackKeys define keys for showback values
type VMShowbackKeys string

const (
	MemCostVMK  VMShowbackKeys = "MEMORY_COST"
	CPUCostVMK  VMShowbackKeys = "CPU_COST"
	DiskCostVMK VMShowbackKeys = "DISK_COST"
)

func (d *VMTemplateBuilder) SetShowback(key VMShowbackKeys, value string) error {
	for _, key := range []VMShowbackKeys{MemCostVMK, CPUCostVMK, DiskCostVMK} {
		if !d.Dynamic.Exists(string(key)) {
			continue
		}
		return fmt.Errorf("VMTemplateBuilder.SetShowbackAttr: the key %s is already present in template", "")
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
	ArchVMK       VMOSKeys = "ARCH"
	MachineVMK    VMOSKeys = "MACHINE"
	KernelVMK     VMOSKeys = "KERNEL"
	KernelDSVMK   VMOSKeys = "KERNEL_DS"
	InitrdVMK     VMOSKeys = "INITRD"
	InitrdDSVMK   VMOSKeys = "INITRD_DS"
	RootVMK       VMOSKeys = "ROOT"
	KernelCmdVMK  VMOSKeys = "KERNEL_CMD"
	BootloaderVMK VMOSKeys = "BOOTLOADER"
	BootVMK       VMOSKeys = "BOOT"
)

func (d *VMTemplateBuilder) AddOS(key VMShowbackKeys, value string) error {
	return d.Dynamic.addPairToVec(vmkOSVec, string(key), value)
}

func (d *VMTemplate) GetOS(key VMOSKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmkOSVec), string(key))
}

// CPU model part

type VMCPUModelKeys string

const (
	ModelVMK VMCPUModelKeys = "MODEL"
)

// There is only one key defined for the CPU_MODEL vector, so we just define a getValue and a setter for it
func (d *VMTemplateBuilder) AddCPUModel(value string) error {
	return d.Dynamic.addPairToVec(vmkCPUModelVec, string(ModelVMK), value)
}

func (d *VMTemplate) GetCPUModel(key VMCPUModelKeys) (string, error) {
	cpuMod, err := d.Dynamic.GetVector(string(vmkCPUModelVec))
	if err != nil {
		return "", fmt.Errorf("VMTemplate.GetCPUModel: vector %s: %s", vmkCPUModelVec, err)
	}
	return cpuMod.getValue(string(key))
}

// Features template part

type VMFeatureKeys string

const (
	PAEVMK              VMFeatureKeys = "PAE"
	ACPIVMK             VMFeatureKeys = "ACPI"
	APICVMK             VMFeatureKeys = "APIC"
	LocalTimeVMK        VMFeatureKeys = "LOCAL_TIME"
	GuestAgentVMK       VMFeatureKeys = "GUEST_AGENT"
	VirtIOScsiQueuesVMK VMFeatureKeys = "VIRTIO_SCSI_QUEUES"
)

func (d *VMTemplateBuilder) AddFeature(key VMFeatureKeys, value string) error {
	return d.Dynamic.addPairToVec(vmkFeaturesVec, string(key), value)
}

func (d *VMTemplate) GetFeature(key VMFeatureKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmkFeaturesVec), string(key))
}

// I/O devices template part

type VMIOGraphicsKeys string
type VMIOInputKeys string

const (
	InputTypeVMK VMIOInputKeys = "TYPE" // Values: mouse or tablet
	BusVMK       VMIOInputKeys = "BUS"  // Values: usb or ps2

	GraphicTypeVMK    VMIOGraphicsKeys = "TYPE" // Values: vnc, sdl, spice
	ListenVMK         VMIOGraphicsKeys = "LISTEN"
	PortVMK           VMIOGraphicsKeys = "PORT"
	PasswdVMK         VMIOGraphicsKeys = "PASSWD"
	KeymapVMK         VMIOGraphicsKeys = "KEYMAP"
	RandomPasswordVMK VMIOGraphicsKeys = "RANDOM_PASSWD"
)

func (d *VMTemplateBuilder) AddIOGraphic(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmkIOGraphicsVec, string(key), value)
}

func (d *VMTemplate) GetIOGraphic(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmkIOGraphicsVec), string(key))
}

func (d *VMTemplateBuilder) AddIOInput(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmkIOGraphicsVec, string(key), value)
}

func (d *VMTemplate) GetIOInput(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getValueFromVec(string(vmkIOGraphicsVec), string(key))
}

// Context template part

type VMContext struct {
	DynamicTemplateVector
}

// VMContextKeys is here to help the user to keep track of XML tags defined in VM context
type VMContextKeys string
type VMContextB64Keys string

const (
	DNSVMK          VMContextKeys = "DNS"
	DNSHostNameVMK  VMContextKeys = "DNS_HOSTNAME"
	EC2PubKeyVMK    VMContextKeys = "EC2_PUBLIC_KEY"
	FilesVMK        VMContextKeys = "FILES"
	FilesDSVMK      VMContextKeys = "FILES_DS"
	GatewayIfaceVMK VMContextKeys = "GATEWAY_IFACE"
	NetworkCtxVMK   VMContextKeys = "NETWORK"
	InitScriptsVMK  VMContextKeys = "INIT_SCRIPTS"
	SSHPubKeyVMK    VMContextKeys = "SSH_PUBLIC_KEY"
	TargetCtxVMK    VMContextKeys = "TARGET"
	TokenVMK        VMContextKeys = "TOKEN"
	UsernameVMK     VMContextKeys = "USERNAME"
	VariableVMK     VMContextKeys = "VARIABLE"
	SecureTTYVMK    VMContextKeys = "SECURETTY"
	SetHostnameVMK  VMContextKeys = "SET_HOSTNAME"

	// Keys for Base64 values
	PasswordB64VMK    VMContextB64Keys = "PASSWORD_BASE64"
	StartScriptB64VMK VMContextB64Keys = "START_SCRIPT_BASE64"
	CryptedPassB64VMK VMContextB64Keys = "CRYPTED_PASSWORD_BASE64"

	// NOTE: ETHx_XXX values not mapped.
)

// GetCtx retrieve a context key
func (t *VMTemplate) GetCtx(key VMContextKeys) (string, error) {
	return t.Dynamic.getValueFromVec(string(vmkContextVec), string(key))
}

// Add adds a context key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplateBuilder) AddCtx(key VMContextKeys, value string) error {
	return t.Dynamic.addPairToVec(vmkContextVec, string(key), value)
}

// Add adds a context key with value. It will convert value to base64. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplateBuilder) AddB64Ctx(key VMContextB64Keys, value string) error {
	valueB64 := base64.StdEncoding.EncodeToString([]byte(value))
	return t.Dynamic.addPairToVec(vmkContextVec, string(key), valueB64)
}

// Placement Template part

type VMPlacementKeys string

const (
	SchedRequirementsVMK   VMPlacementKeys = "SCHED_REQUIREMENTS"
	SchedRankVMK           VMPlacementKeys = "SCHED_RANK"
	SchedDSRequirementsVMK VMPlacementKeys = "SCHED_DS_REQUIREMENTS"
	SchedDSRankVMK         VMPlacementKeys = "SCHED_DS_RANK"
	UserPriorityVMK        VMPlacementKeys = "USER_PRIORITY"
)

// SetPlacement add once a placement attribute
func (d *VMTemplateBuilder) SetPlacement(key VMPlacementKeys, value string) error {
	for _, key := range []VMPlacementKeys{
		SchedRequirementsVMK,
		SchedRankVMK,
		SchedDSRequirementsVMK,
		SchedDSRankVMK,
		UserPriorityVMK,
	} {
		if !d.Dynamic.Exists(string(key)) {
			continue
		}
		return fmt.Errorf("VMTemplateBuilder.AddPlacement: the key %s is already present in template", "")
	}
	return d.Dynamic.AddPair(string(key), value)
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetPlacement(key VMPlacementKeys) (string, error) {
	return d.Dynamic.getValue(string(key))
}

// Scheduled actions template part
// VMSchedActionKeys is a type to describe scheduled action key
type VMSchedActionKeys string

const (
	TimeVMK     VMSchedActionKeys = "TIME"
	RepeatVMK   VMSchedActionKeys = "REPEAT"
	DaysVMK     VMSchedActionKeys = "DAYS"
	ActionVMK   VMSchedActionKeys = "ACTION"
	EndTypeVMK  VMSchedActionKeys = "END_TYPE"
	EndValueVMK VMSchedActionKeys = "END_VALUE"
)

// SchedAction is a scheduled action on VM
type SchedAction struct {
	DynamicTemplateVector
}

// NewSchedAction returns a structure disk entity to build
func (a *VMTemplate) NewSchedAction() *SchedAction {
	return &SchedAction{
		DynamicTemplateVector{key: vmkSchedActionVec},
	}
}

// Add adds a SchedAction key
func (t *SchedAction) Add(key VMSchedActionKeys, value string) error {
	return t.AddPair(string(key), value)
}

// Get retrieve a SchedAction key
func (t *SchedAction) Get(key VMSchedActionKeys) (string, error) {
	return t.getValue(string(key))
}

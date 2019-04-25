package goca

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// User template

// VMUserTemplate contains common and custom attributes
type VMUserTemplate struct {
	Error        string             `xml:"ERROR"`
	SchedMessage string             `xml:"SCHED_MESSAGE"`
	Dynamic      dynamicTemplateAny `xml:",any"`
}

// Get allows to get a pair by key
func (d *VMUserTemplate) Get(key string) (string, error) {
	return d.Dynamic.GetStr(string(key))
}

// VM Template
// Available template parts and keys are listed here: https://docs.opennebula.org/5.8/operation/references/template.html
// Some specific part are not defined: vCenter, Public Cloud, Hypervisor, User Inputs

// Name of the vector in the template, they should match with tags in VMTemplate struct
const (
	vmContextVecK     string = "CONTEXT"
	vmCPUModelVecK    string = "CPU_MODEL"
	vmFeaturesVecK    string = "FEATURES"
	vmIOGraphicsVecK  string = "GRAPHICS"
	vmIOInputVecK     string = "INPUT"
	vmOSVecK          string = "OS"
	vmSchedActionVecK string = "SCHED_ACTION"
)

// VMTemplate is a structure allowing to parse VM templates.
// It's defined in a semi-static way to guide the user among the bunch of values
type VMTemplate struct {
	VMCapacity

	Context      VMContext     `xml:"CONTEXT"`
	Disks        []Disk        `xml:"DISK"`
	NICs         []NIC         `xml:"NIC"`
	SchedActions []SchedAction `xml:"SCHED_ACTION"`

	Dynamic dynamicTemplateAny `xml:",any"`
}

func (t *VMTemplate) String() string {
	var s strings.Builder

	s.WriteString(t.VMCapacity.String())
	s.WriteString(t.Context.String() + "\n")
	for _, disk := range t.Disks {
		s.WriteString(disk.String() + "\n")
	}
	for _, nic := range t.NICs {
		s.WriteString(nic.String() + "\n")
	}
	for _, a := range t.SchedActions {
		s.WriteString(a.String() + "\n")
	}

	s.WriteString(t.Dynamic.String())

	return s.String()
}

// NewVMTemplate returns a VMTemplate structure
func NewVMTemplate() *VMTemplate {
	return &VMTemplate{}
}

// Template parts

// Capacity template part

type VMCapacity struct {
	CPU    float64 `xml:"CPU"`
	Memory uint    `xml:"MEMORY"`
	VCPU   uint    `xml:"VCPU"` // If not defined: it's 1
}

func (c *VMCapacity) String() string {
	return fmt.Sprintf("CPU=%f\nMemory=%d\nVCPU=%d", c.CPU, c.Memory, c.VCPU)
}

// NewVMCapacity returns a VM capacity structure. You may need this when resizing a VM
func NewVMCapacity(CPU float64, VCPU, Memory uint) *VMCapacity {
	return &VMCapacity{
		CPU:    CPU,
		Memory: Memory,
		VCPU:   VCPU,
	}
}

// SetCapacity set capacity attributes in VM template
func (t *VMTemplate) SetCapacity(CPU float64, VCPU, Memory uint) {
	t.CPU = CPU
	t.Memory = Memory
	t.VCPU = VCPU
}

// GetCapacity gets capacity attributes from a VM template
func (t *VMTemplate) GetCapacity() (CPU float64, VCPU uint, Memory uint) {
	CPU = t.CPU
	Memory = t.Memory
	if t.VCPU == 0 {
		VCPU = 1
	} else {
		VCPU = t.VCPU
	}
	return
}

// AddDisk allow to add a disk to the template
func (t *VMTemplate) AddDisk(d *Disk) {
	t.Disks = append(t.Disks, *d)
}

// AddNIC allow to add a NIC to the template
func (t *VMTemplate) AddNIC(n *NIC) {
	t.NICs = append(t.NICs, *n)
}

// AddSchedAction returns a structure disk entity to build
func (t *VMTemplate) AddSchedAction(a *SchedAction) {
	t.SchedActions = append(t.SchedActions, *a)
}

// Show back template part

// VMShowbackKeys define keys for showback values
type VMShowbackKeys string

const (
	MemCostVMK  VMShowbackKeys = "MEMORY_COST"
	CPUCostVMK  VMShowbackKeys = "CPU_COST"
	DiskCostVMK VMShowbackKeys = "DISK_COST"
)

func (d *VMTemplate) SetShowback(key VMShowbackKeys, value string) error {
	d.Dynamic.Del(string(key))
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetShowback(key VMShowbackKeys) (string, error) {
	return d.Dynamic.GetStr(string(key))
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

func (d *VMTemplate) AddOS(key VMShowbackKeys, value string) error {
	return d.Dynamic.addPairToVec(vmOSVecK, string(key), value)
}

func (d *VMTemplate) GetOS(key VMOSKeys) (string, error) {
	return d.Dynamic.getStrFromVec(string(vmOSVecK), string(key))
}

// CPU model part

type VMCPUModelKeys string

const (
	ModelVMK VMCPUModelKeys = "MODEL"
)

// There is only one key defined for the CPU_MODEL vector, so we just define a GetStr and a setter for it
func (d *VMTemplate) AddCPUModel(value string) error {
	return d.Dynamic.addPairToVec(vmCPUModelVecK, string(ModelVMK), value)
}

func (d *VMTemplate) GetCPUModel(key VMCPUModelKeys) (string, error) {
	cpuMod, err := d.Dynamic.GetVector(string(vmCPUModelVecK))
	if err != nil {
		return "", fmt.Errorf("VMTemplate.GetCPUModel: vector %s: %s", vmCPUModelVecK, err)
	}
	return cpuMod.GetStr(string(key))
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

func (d *VMTemplate) AddFeature(key VMFeatureKeys, value string) error {
	return d.Dynamic.addPairToVec(vmFeaturesVecK, string(key), value)
}

func (d *VMTemplate) GetFeature(key VMFeatureKeys) (string, error) {
	return d.Dynamic.getStrFromVec(string(vmFeaturesVecK), string(key))
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

func (d *VMTemplate) AddIOGraphic(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmIOGraphicsVecK, string(key), value)
}

func (d *VMTemplate) GetIOGraphic(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getStrFromVec(string(vmIOGraphicsVecK), string(key))
}

func (d *VMTemplate) AddIOInput(key VMIOGraphicsKeys, value string) error {
	return d.Dynamic.addPairToVec(vmIOGraphicsVecK, string(key), value)
}

func (d *VMTemplate) GetIOInput(key VMIOInputKeys) (string, error) {
	return d.Dynamic.getStrFromVec(string(vmIOGraphicsVecK), string(key))
}

// Context template part

type VMContext struct {
	TemplateVector
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

	// NOTE: ETHx_XXX values are not mapped.
)

// GetCtx retrieve a context key
func (t *VMTemplate) GetCtx(key VMContextKeys) (string, error) {
	return t.Context.GetStr(string(key))
}

// Add adds a context key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplate) AddCtx(key VMContextKeys, value string) error {
	return t.Context.AddPair(string(key), value)
}

// Add adds a context key with value. It will convert value to base64. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *VMTemplate) AddB64Ctx(key VMContextB64Keys, value string) error {
	valueB64 := base64.StdEncoding.EncodeToString([]byte(value))
	return t.Context.AddPair(string(key), valueB64)
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
func (d *VMTemplate) SetPlacement(key VMPlacementKeys, value string) error {
	d.Dynamic.Del(string(key))
	return d.Dynamic.AddPair(string(key), value)
}

func (d *VMTemplate) GetPlacement(key VMPlacementKeys) (string, error) {
	return d.Dynamic.GetStr(string(key))
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
	TemplateVector
}

// NewSchedAction returns a structure disk entity to build
func (a *VMTemplate) NewSchedAction() *SchedAction {
	return &SchedAction{
		TemplateVector{key: vmSchedActionVecK},
	}
}

// Add adds a SchedAction key
func (t *SchedAction) Add(key VMSchedActionKeys, value string) error {
	return t.AddPair(string(key), value)
}

// Get retrieve a SchedAction key
func (t *SchedAction) Get(key VMSchedActionKeys) (string, error) {
	return t.GetStr(string(key))
}

package goca

// Disk is a structure allowing to parse disk templates
type Disk struct {
	DynamicTemplateVector
}

// DiskKeys is here to help the user to keep track of XML tags defined in Disk
type DiskKeys string

// Some keys are specific to VM some others to VRouter
const (
	vecDK          string   = "DISK"
	idDK           DiskKeys = "DISK_ID"
	DatastoreDK    DiskKeys = "DATASTORE"
	DiskTypeDK     DiskKeys = "DISK_TYPE"
	Driver         DiskKeys = "DRIVER"
	ImageDK        DiskKeys = "IMAGE"
	ImageIDDK      DiskKeys = "IMAGE_ID"
	ImageUnameDK   DiskKeys = "IMAGE_UNAME"
	OriginalSizeDK DiskKeys = "ORIGINAL_SIZE"
	SizeDK         DiskKeys = "SIZE"
	TargetDiskDK   DiskKeys = "TARGET"
	TypeDK         DiskKeys = "TYPE"
)

// NewDisk returns a structure disk entity to build
func NewDisk() *Disk {
	return &Disk{
		DynamicTemplateVector{key: vecDK},
	}
}

// Get is a getValue for all Disk keys
func (d *Disk) Get(key DiskKeys) (string, error) {
	return d.getValue(string(key))
}

// GetID is a getValue for the Disk ID doing the conversion to an unsigned integer
func (d *Disk) GetID() (uint, error) {
	return d.getID(string(idDK))
}

// Add adds a Disk key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (d *Disk) Add(key DiskKeys, value string) error {
	return d.AddPair(string(key), value)
}

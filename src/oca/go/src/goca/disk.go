package goca

// Disk is a structure allowing to parse disk templates
type Disk struct {
	DynamicTemplateVector
}

// DiskKeys is here to help the user to keep track of XML tags defined in Disk
type DiskKeys string

// Some keys are specific to VM some others to VRouter
const (
	DatastoreK    DiskKeys = "DATASTORE"
	DiskIDK       DiskKeys = "DISK_ID"
	DiskTypeK     DiskKeys = "DISK_TYPE"
	DriverK       DiskKeys = "DRIVER"
	ImageK        DiskKeys = "IMAGE"
	ImageIDK      DiskKeys = "IMAGE_ID"
	ImageUnameK   DiskKeys = "IMAGE_UNAME"
	OriginalSizeK DiskKeys = "ORIGINAL_SIZE"
	SizeK         DiskKeys = "SIZE"
	TargetDiskK   DiskKeys = "TARGET"
	TypeK         DiskKeys = "TYPE"
)

// NewDisk returns a structure disk entity to build
func NewDisk() Disk {
	return Disk{
		DynamicTemplateVector{key: "DISK"},
	}
}

// Get is a getValue for all Disk keys
func (d *Disk) Get(key DiskKeys) (string, error) {
	return d.getValue(string(key))
}

// GetID is a getValue for the Disk ID doing the conversion to an unsigned integer
func (d *Disk) GetID() (uint, error) {
	return d.getID(string(DiskIDK))
}

// Add adds a Disk key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (d *Disk) Add(key DiskKeys, value string) error {
	return d.AddPair(string(key), value)
}

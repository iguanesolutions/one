package goca

import "strconv"

// Disk is a structure allowing to parse/build disk templates
type Disk struct {
	DynTemplateVector
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
	return Disk{DynTemplateVector{key: "DISK"}}
}

// Get is a getter for all Disk keys
func (d *Disk) Get(key DiskKeys) (string, error) {
	return d.getterVec(string(key))
}

// Add adds a Disk key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (d *Disk) Add(key DiskKeys, value string) error {
	return d.AddPair(string(key), value)
}

// GetID is a getter for the Disk ID doing the conversion to an unsigned integer
func (d *Disk) GetID() (uint, error) {
	idStr, err := d.getterVec(string(DiskIDK))
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(id), err
}

package goca

// Disk is a structure allowing to parse disk templates
type Disk struct {
	TemplateVector
}

// DiskKeys is here to help the user to keep track of XML tags defined in Disk
type DiskKeys string

// Some keys are specific to VM some others to VRouter
const (
	vecDK          string   = "DISK"
	idDK           DiskKeys = "DISK_ID"
	DatastoreDK    DiskKeys = "DATASTORE"
	DatastoreIDDK  DiskKeys = "DATASTORE_ID"
	DiskTypeDK     DiskKeys = "DISK_TYPE"
	Driver         DiskKeys = "DRIVER"
	ImageDK        DiskKeys = "IMAGE"
	ImageIDDK      DiskKeys = "IMAGE_ID"
	ImageUnameDK   DiskKeys = "IMAGE_UNAME"
	OriginalSizeDK DiskKeys = "ORIGINAL_SIZE"
	SizeDK         DiskKeys = "SIZE"
	TargetDiskDK   DiskKeys = "TARGET"
)

// NewDisk returns a structure disk entity to build
func NewDisk() *Disk {
	return &Disk{
		TemplateVector{key: vecDK},
	}
}

// ID returns the disk ID
func (d *Disk) ID() int {
	id, _ := d.GetInt(string(idDK))
	return id
}

// Get return the string value of a Disk key
func (d *Disk) Get(key DiskKeys) (string, error) {
	return d.GetStr(string(key))
}

// GetID convert and returns a value as an ID. The key name generally ends with an "_ID"
func (d *Disk) GetID(key DiskKeys) (int, error) {
	return d.GetInt(string(key))
}

// Add adds a Disk key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (d *Disk) Add(key DiskKeys, value string) error {
	return d.AddPair(string(key), value)
}

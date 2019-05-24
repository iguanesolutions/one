/* -------------------------------------------------------------------------- */
/* Copyright 2002-2019, OpenNebula Project, OpenNebula Systems                */
/*                                                                            */
/* Licensed under the Apache License, Version 2.0 (the "License"); you may    */
/* not use this file except in compliance with the License. You may obtain    */
/* a copy of the License at                                                   */
/*                                                                            */
/* http://www.apache.org/licenses/LICENSE-2.0                                 */
/*                                                                            */
/* Unless required by applicable law or agreed to in writing, software        */
/* distributed under the License is distributed on an "AS IS" BASIS,          */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   */
/* See the License for the specific language governing permissions and        */
/* limitations under the License.                                             */
/*--------------------------------------------------------------------------- */

package goca

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// ImagesController is a controller for Images
type ImagesController entitiesController

// ImageController is a controller for Image entities
type ImageController entityController

// ImageSnapshotController is a controller for an Image snapshot
type ImageSnapshotController subEntityController

// ImagePool represents an OpenNebula Image pool
type ImagePool struct {
	Images []Image `xml:"IMAGE"`
}

// Image represents an OpenNebula Image
type Image struct {
	ID              int           `xml:"ID"`
	UID             int           `xml:"UID"`
	GID             int           `xml:"GID"`
	UName           string        `xml:"UNAME"`
	GName           string        `xml:"GNAME"`
	Name            string        `xml:"NAME"`
	LockInfos       *Lock         `xml:"LOCK"`
	Permissions     *Permissions  `xml:"PERMISSIONS"`
	Type            int           `xml:"TYPE"`
	DiskType        int           `xml:"DISK_TYPE"`
	PersistentValue int           `xml:"PERSISTENT"`
	RegTime         int           `xml:"REGTIME"`
	Source          string        `xml:"SOURCE"`
	Path            string        `xml:"PATH"`
	FsType          string        `xml:"FSTYPE"`
	Size            int           `xml:"SIZE"`
	StateRaw        int           `xml:"STATE"`
	RunningVMs      int           `xml:"RUNNING_VMS"`
	CloningOps      int           `xml:"CLONING_OPS"`
	CloningID       int           `xml:"CLONING_ID"`
	TargetSnapshot  int           `xml:"TARGET_SNAPSHOT"`
	DatastoreID     int           `xml:"DATASTORE_ID"`
	Datastore       string        `xml:"DATASTORE"`
	VMsID           []int         `xml:"VMS>ID"`
	ClonesID        []int         `xml:"CLONES>ID"`
	AppClonesID     []int         `xml:"APP_CLONES>ID"`
	Snapshots       ImageSnapshot `xml:"SNAPSHOTS"`
	Template        ImageTemplate `xml:"TEMPLATE"`
}

// ImageState is the state of the Image
type ImageState int

const (
	// ImageInit image is being initialized
	ImageInit ImageState = iota

	// ImageReady image is ready to be used
	ImageReady

	// ImageUsed image is in use
	ImageUsed

	// ImageDisabled image is in disabled
	ImageDisabled

	// ImageLocked image is locked
	ImageLocked

	// ImageError image is in error state
	ImageError

	// ImageClone image is in clone state
	ImageClone

	// ImageDelete image is in delete state
	ImageDelete

	// ImageUsedPers image is in use and persistent
	ImageUsedPers

	// ImageLockUsed image is in locked state (non-persistent)
	ImageLockUsed

	// ImageLockUsedPers image is in locked state (persistent)
	ImageLockUsedPers
)

func (s ImageState) isValid() bool {
	if s >= ImageInit && s <= ImageLockUsedPers {
		return true
	}
	return false
}

// String returns the string version of the ImageState
func (s ImageState) String() string {
	return [...]string{
		"INIT",
		"READY",
		"USED",
		"DISABLED",
		"LOCKED",
		"ERROR",
		"CLONE",
		"DELETE",
		"USED_PERS",
		"LOCKED_USED",
		"LOCKED_USED_PERS",
	}[s]
}

// Images returns an Images controller
func (c *Controller) Images() *ImagesController {
	return &ImagesController{c}
}

// Image returns an Image controller
func (c *Controller) Image(id int) *ImageController {
	return &ImageController{c, id}
}

// Snapshot returns an Image snapshot controller
func (ic *ImageController) Snapshot(id int) *ImageSnapshotController {
	return &ImageSnapshotController{ic.c, ic.ID, id}
}

// ByName returns an Image ID from name
func (ic *ImagesController) ByName(name string, v *View) (int, error) {
	ids, err := ic.info(
		func(d *Image) (bool, error) {
			return d.Name == name, nil
		},
		v)
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, errors.New("resource not found")
	} else if len(ids) > 1 {
		return 0, errors.New("multiple resources with that name")
	}

	return ids[0], nil
}

// info is the base function to apply by attribute matching
func (ic *ImagesController) info(fn func(*Image) (bool, error), v *View) ([]int, error) {
	var ret []int

	pool, err := ic.Info(v)
	if err != nil {
		return ret, err
	}

	var ok bool
	for i := 0; i < len(pool.Images); i++ {
		ok, err = fn(&pool.Images[i])
		if !ok {
			continue
		}

		ret = append(ret, pool.Images[i].ID)
	}

	return ret, nil
}

// Info returns a new image pool. It accepts the scope of the query.
func (ic *ImagesController) Info(v *View) (*ImagePool, error) {

	response, err := ic.c.Client.Call("one.imagepool.info", v.who, v.id.start, v.id.end)
	if err != nil {
		return nil, err
	}

	pool := &ImagePool{}
	err = xml.Unmarshal([]byte(response.Body()), pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// Info connects to OpenNebula and fetches the information of the Image
func (ic *ImageController) Info() (*Image, error) {
	response, err := ic.c.Client.Call("one.image.info", ic.ID)
	if err != nil {
		return nil, err
	}
	image := &Image{}
	err = xml.Unmarshal([]byte(response.Body()), image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

// Create allocates a new image based on the template string provided. It
// returns the image ID.
func (ic *ImagesController) Create(name string, itype ImageTypesValues, dsid int, tpl *ImageTemplate) (int, error) {
	if tpl == nil {
		return 0, fmt.Errorf("Image Create: nil template arg")
	}
	tpl.SetName(name)
	tpl.SetType(itype)

	response, err := ic.c.Client.Call("one.image.allocate", tpl.String(), dsid)
	if err != nil {
		return 0, err
	}

	return response.BodyInt(), nil
}

// State looks up the state of the image and returns the ImageState
func (image *Image) State() (ImageState, error) {
	state := ImageState(image.StateRaw)
	if !state.isValid() {
		return -1, fmt.Errorf("Image State: this state value is not currently handled: %d\n", image.StateRaw)
	}
	return state, nil
}

// StateString returns the state in string format
func (image *Image) StateString() (string, error) {
	state := ImageState(image.StateRaw)
	if !state.isValid() {
		return "", fmt.Errorf("Image State: this state value is not currently handled: %d\n", image.StateRaw)
	}
	return state.String(), nil
}

// Clone clones an existing image. It returns the clone ID
func (ic *ImageController) Clone(cloneName string, dsid int) (int, error) {
	response, err := ic.c.Client.Call("one.image.clone", ic.ID, cloneName, dsid)
	if err != nil {
		return 0, err
	}

	return response.BodyInt(), nil
}

// Update replaces the image contents.
// * tpl: The new image contents.
// * uType: Update type: Replace: Replace the whole template.
//   Merge: Merge new template with the existing one.
func (ic *ImageController) Update(tpl *ImageTemplate, uType UpdateType) error {
	if tpl == nil {
		return fmt.Errorf("Image Update: nil template")
	}
	_, err := ic.c.Client.Call("one.image.update", ic.ID, tpl, uType)
	return err
}

// Chtype changes the type of the Image
func (ic *ImageController) Chtype(newType string) error {
	_, err := ic.c.Client.Call("one.image.chtype", ic.ID, newType)
	return err
}

// Chown changes the owner/group of the image. If uid or gid is -1 it will not
// change
func (ic *ImageController) Chown(uid, gid int) error {
	_, err := ic.c.Client.Call("one.image.chown", ic.ID, uid, gid)
	return err
}

// Chmod changes the permissions of the image. If any perm is -1 it will not
// change
func (ic *ImageController) Chmod(perm *Permissions) error {
	_, err := ic.c.Client.Call("one.image.chmod", perm.toArgs(ic.ID)...)
	return err
}

// Rename changes the name of the image
func (ic *ImageController) Rename(newName string) error {
	_, err := ic.c.Client.Call("one.image.rename", ic.ID, newName)
	return err
}

// Delete will delete a snapshot from the image
func (ic *ImageSnapshotController) Delete() error {
	_, err := ic.c.Client.Call("one.image.snapshotdelete", ic.entityID, ic.ID)
	return err
}

// Revert reverts image state to a previous snapshot
func (ic *ImageSnapshotController) Revert() error {
	_, err := ic.c.Client.Call("one.image.snapshotrevert", ic.entityID, ic.ID)
	return err
}

// Flatten flattens the snapshot image and discards others
func (ic *ImageSnapshotController) Flatten() error {
	_, err := ic.c.Client.Call("one.image.snapshotflatten", ic.entityID, ic.ID)
	return err
}

// Enable enables (or disables) the image
func (ic *ImageController) Enable(enable bool) error {
	_, err := ic.c.Client.Call("one.image.enable", ic.ID, enable)
	return err
}

// Persistent sets the image as persistent (or not)
func (ic *ImageController) Persistent(persistent bool) error {
	_, err := ic.c.Client.Call("one.image.persistent", ic.ID, persistent)
	return err
}

// Lock locks the image following lock level. See levels in locks.go.
func (ic *ImageController) Lock(level LockLevel) error {
	_, err := ic.c.Client.Call("one.image.lock", ic.ID, level)
	return err
}

// Unlock unlocks the image.
func (ic *ImageController) Unlock() error {
	_, err := ic.c.Client.Call("one.image.unlock", ic.ID)
	return err
}

// Delete will remove the image from OpenNebula, which will remove it from the
// backend.
func (ic *ImageController) Delete() error {
	_, err := ic.c.Client.Call("one.image.delete", ic.ID)
	return err
}

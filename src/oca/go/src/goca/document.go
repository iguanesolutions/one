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

// DocumentsController is a controller for documents entities
type DocumentsController struct {
	entitiesController
	dType int
}

// DocumentController is a controller for document entity
type DocumentController entityController

// DocumentPool represents an OpenNebula DocumentPool
type DocumentPool struct {
	Documents []Document `xml:"DOCUMENT"`
}

// Document represents an OpenNebula Document
type Document struct {
	ID          int             `xml:"ID"`
	UID         int             `xml:"UID"`
	GID         int             `xml:"GID"`
	UName       string          `xml:"UNAME"`
	GName       string          `xml:"GNAME"`
	Name        string          `xml:"NAME"`
	Type        string          `xml:"TYPE"`
	Permissions *Permissions    `xml:"PERMISSIONS"`
	LockInfos   *Lock           `xml:"LOCK"`
	Template    DynamicTemplate `xml:"TEMPLATE"`
}

// Documents returns a Documents controller
func (c *Controller) Documents(dType int) *DocumentsController {
	return &DocumentsController{entitiesController{c}, dType}
}

// Document returns a Document controller
func (c *Controller) Document(id int) *DocumentController {
	return &DocumentController{c, id}
}

// ByName returns an Image ID from name
func (dc *DocumentsController) ByName(name string, v *View) (int, error) {
	ids, err := dc.info(
		func(i *Document) (bool, error) {
			return i.Name == name, nil
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

/*
// ByState returns a list of image ID from state
func (c *DocumentsController) ByState(state ImageState, v *View) ([]int, error) {
	return c.info(
		func(i *Document) (bool, error) {
			state, err := i.State()
			if err != nil {
				return false, err
			}
			return ImageState(i.StateRaw) == state, nil
		},
		v)
}
*/

// ByPair returns an Image from a template pair
func (dc *DocumentsController) ByPair(p TemplatePair, v *View) ([]int, error) {
	return c.info(
		func(d *Document) (bool, error) {
			return d.Template.findPair(p), nil
		},
		v)
}

// info is the base function to apply by attribute matching
func (dc *DocumentsController) info(fn func(*Document) (bool, error), v *View) ([]int, error) {
	var ret []int

	pool, err := dc.Info(args...)
	if err != nil {
		return ret, err
	}

	var ok bool
	for i := 0; i < len(pool.Documents); i++ {
		ok, err = fn(&pool.Documents[i])
		if !ok {
			continue
		}

		ret = append(ret, pool.Documents[i].ID)
	}

	return ret, nil
}

// Info returns a document pool. A connection to OpenNebula is
// performed.
func (dc *DocumentsController) Info(v *View) (*DocumentPool, error) {

	v, err := NewView(args...)
	if err != nil {
		return nil, err
	}

	response, err := dc.c.Client.Call("one.documentpool.info", v.who, v.id.start, v.id.end, dc.dType)
	if err != nil {
		return nil, err
	}

	pool := &DocumentPool{}
	err = xml.Unmarshal([]byte(response.Body()), pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// Info retrieves information for the document.
func (dc *DocumentController) Info() (*Document, error) {
	response, err := dc.c.Client.Call("one.document.info", dc.ID)
	if err != nil {
		return nil, err
	}
	document := &Document{}
	err = xml.Unmarshal([]byte(response.Body()), document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// Create allocates a new document. It returns the new document ID.
func (dc *DocumentsController) Create(name string, tpl *DynamicTemplate) (int, error) {
	if tpl == nil {
		return 0, fmt.Errorf("Document Create: nil template arg")
	}
	tpl.SetName(name)

	response, err := dc.c.Client.Call("one.document.allocate", tpl.String(), dc.dType)
	if err != nil {
		return 0, err
	}

	return response.BodyInt(), nil
}

// Clone clones an existing document.
// * newName: Name for the new document.
func (dc *DocumentController) Clone(newName string) error {
	_, err := dc.c.Client.Call("one.document.clone", dc.ID, newName)
	return err
}

// Delete deletes the given document from the pool.
func (dc *DocumentController) Delete() error {
	_, err := dc.c.Client.Call("one.document.delete", dc.ID)
	return err
}

// Update replaces the document contents.
// * tpl: The new document contents.
// * uType: Update type: Replace: Replace the whole template.
//   Merge: Merge new template with the existing one.
func (dc *DocumentController) Update(tpl *DynamicTemplate, uType UpdateType) error {
	if tpl == nil {
		return fmt.Errorf("Document Update: empty template")
	}
	_, err := dc.c.Client.Call("one.document.update", dc.ID, tpl.String(), uType)
	return err
}

// Chmod changes the permission bits of a document.
func (dc *DocumentController) Chmod(perm *Permissions) error {
	_, err := dc.c.Client.Call("one.document.chmod", perm.toArgs(dc.ID)...)
	return err
}

// Chown changes the ownership of a document.
// * userID: The User ID of the new owner. If set to -1, it will not change.
// * groupID: The Group ID of the new group. If set to -1, it will not change.
func (dc *DocumentController) Chown(userID, groupID int) error {
	_, err := dc.c.Client.Call("one.document.chown", dc.ID, userID, groupID)
	return err
}

// Rename renames a document.
// * newName: The new name.
func (dc *DocumentController) Rename(newName string) error {
	_, err := dc.c.Client.Call("one.document.rename", dc.ID, newName)
	return err
}

// Lock locks the document following lock level. See levels in locks.go.
func (dc *DocumentController) Lock(level LockLevel) error {
	_, err := dc.c.Client.Call("one.document.lock", dc.ID, level)
	return err
}

// Unlock unlocks the document.
func (dc *DocumentController) Unlock() error {
	_, err := dc.c.Client.Call("one.document.unlock", dc.ID)
	return err
}

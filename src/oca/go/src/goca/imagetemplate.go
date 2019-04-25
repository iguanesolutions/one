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

type ImageTemplateKeys string

// Image template part
const (
	PersistentIK     ImageTemplateKeys = "PERSISTENT"
	PersistentTypeIK ImageTemplateKeys = "PERSISTENT_TYPE"
	SizeIK           ImageTemplateKeys = "SIZE"
	DevPrefixIK      ImageTemplateKeys = "DEV_PREFIX"
	TargetIK         ImageTemplateKeys = "TARGET"
	DriverIK         ImageTemplateKeys = "DRIVER"
	PathIK           ImageTemplateKeys = "PATH"
	SourceIK         ImageTemplateKeys = "SOURCE"
	DiskTypeIK       ImageTemplateKeys = "DISK_TYPE"
	ReadOnlyIK       ImageTemplateKeys = "READONLY"
	Md5IK            ImageTemplateKeys = "MD5"
	Sha1IK           ImageTemplateKeys = "SHA1"
)

type ImageTypesValues string

const (
	// Virtual Machine disks
	ImgDatablock ImageTypesValues = "DATABLOCK"
	ImgCDRom     ImageTypesValues = "CDROM"
	ImgOS        ImageTypesValues = "OS"

	// File types, can be registered only in File Datastores
	ImgKernel  ImageTypesValues = "KERNEL"
	ImgRamDisk ImageTypesValues = "RAMDISK"
	ImgContext ImageTypesValues = "CONTEXT"
)

// ImageTemplate is a dynamic part of the image entity
type ImageTemplate struct {
	DynamicTemplate
}

// NewImageTemplate returns a template
func NewImageTemplate() *ImageTemplate {
	return &ImageTemplate{
		DynamicTemplate{},
	}
}

// Get return the string value of a template image key
func (t *ImageTemplate) Get(key ImageTemplateKeys) (string, error) {
	return t.GetStr(string(key))
}

// Add adds a ImageTemplate key with value. NOT ALL KEYS SHOULD BE ADDED, see the documentation
func (t *ImageTemplate) Add(key ImageTemplateKeys, value string) error {
	return t.AddPair(string(key), value)
}

// SetType set an Image type
func (t *ImageTemplate) SetType(typ ImageTypesValues) error {
	t.Del(TypeK)
	return t.AddPair(TypeK, typ)
}

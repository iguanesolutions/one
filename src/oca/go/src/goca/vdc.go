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

// VdcsController is a controller for a pool of Vdcs
type VdcsController entitiesController

// VdcController is a controller for Vdc entities
type VdcController entityController

// VdcPool represents an OpenNebula VdcPool
type VdcPool struct {
	Vdcs []Vdc `xml:"VDC"`
}

// Vdc represents an OpenNebula Vdc
type Vdc struct {
	ID         int            `xml:"ID"`
	Name       string         `xml:"NAME"`
	GroupsID   []int          `xml:"GROUPS>ID"`
	Clusters   []vdcCluster   `xml:"CLUSTERS>CLUSTER"`
	Hosts      []vdcHost      `xml:"HOSTS>HOST"`
	Datastores []vdcDatastore `xml:"DATASTORES>DATASTORE"`
	VNets      []vdcVNet      `xml:"VNETS>VNET"`
	Template   vdcTemplate    `xml:"TEMPLATE"`
}

type vdcTemplate struct {
	DynamicTemplate
}

type vdcCluster struct {
	ZoneID    int `xml:"ZONE_ID"`
	ClusterID int `xml:"CLUSTER_ID"`
}

type vdcHost struct {
	ZoneID int `xml:"ZONE_ID"`
	HostID int `xml:"HOST_ID"`
}

type vdcDatastore struct {
	ZoneID      int `xml:"ZONE_ID"`
	DatastoreID int `xml:"DATASTORE_ID"`
}

type vdcVNet struct {
	ZoneID int `xml:"ZONE_ID"`
	VnetID int `xml:"VNET_ID"`
}

// Vdcs returns a Vdcs controller.
func (c *Controller) Vdcs() *VdcsController {
	return &VdcsController{c}
}

// Vdc returns a Vdc controller
func (c *Controller) Vdc(id int) *VdcController {
	return &VdcController{c, id}
}

// ByName returns a Vdc ID from name
func (c *VdcsController) ByName(name string) (int, error) {
	var id int

	vdcPool, err := c.Info()
	if err != nil {
		return 0, err
	}

	match := false
	for i := 0; i < len(vdcPool.Vdcs); i++ {
		if vdcPool.Vdcs[i].Name != name {
			continue
		}
		if match {
			return 0, errors.New("multiple resources with that name")
		}
		id = vdcPool.Vdcs[i].ID
		match = true
	}
	if !match {
		return 0, errors.New("resource not found")
	}

	return id, nil
}

// Info returns a vdc pool. A connection to OpenNebula is
// performed.
func (vc *VdcsController) Info() (*VdcPool, error) {
	response, err := vc.c.Client.Call("one.vdcpool.info")
	if err != nil {
		return nil, err
	}

	vdcPool := &VdcPool{}
	err = xml.Unmarshal([]byte(response.Body()), vdcPool)
	if err != nil {
		return nil, err
	}

	return vdcPool, nil
}

// Info retrieves information for the VDC.
func (vc *VdcController) Info() (*Vdc, error) {
	response, err := vc.c.Client.Call("one.vdc.info", vc.ID)
	if err != nil {
		return nil, err
	}
	vdc := &Vdc{}
	err = xml.Unmarshal([]byte(response.Body()), vdc)
	if err != nil {
		return nil, err
	}

	return vdc, nil
}

// Create allocates a new vdc. It returns the new vdc ID.
// * name:  The name of the Vdc
// * tpl:	A string containing the template of the VDC. Syntax can be the usual
//     attribute=value or XML.
// * clusterID: The cluster ID. If it is -1, this virtual network won’t be added
//     to any cluster
func (vc *VdcsController) Create(name string, clusterID int, tpl *DynamicTemplate) (int, error) {
	if tpl == nil {
		return 0, fmt.Errorf("Vdc Create: nil template arg")
	}
	tpl.SetName(name)

	response, err := vc.c.Client.Call("one.vdc.allocate", tpl.String(), clusterID)
	if err != nil {
		return 0, err
	}

	return response.BodyInt(), nil
}

// Delete deletes the given VDC from the pool.
func (vc *VdcController) Delete() error {
	_, err := vc.c.Client.Call("one.vdc.delete", vc.ID)
	return err
}

// Update replaces the cluster cluster contents.
// * tpl: The new cluster contents. Syntax can be the usual attribute=value or XML.
// * uType: Update type: Replace: Replace the whole template.
//   Merge: Merge new template with the existing one.
func (vc *VdcController) Update(tpl *DynamicTemplate, uType UpdateType) error {
	if tpl == nil {
		return fmt.Errorf("Vdc Update: empty template")
	}
	_, err := vc.c.Client.Call("one.vdc.update", vc.ID, tpl.String(), uType)
	return err
}

// Rename renames a VDC.
// * newName: The new name.
func (vc *VdcController) Rename(newName string) error {
	_, err := vc.c.Client.Call("one.vdc.rename", vc.ID, newName)
	return err
}

// AddGroup adds a group to the VDC
// * groupID: The group ID.
func (vc *VdcController) AddGroup(groupID int) error {
	_, err := vc.c.Client.Call("one.vdc.addgroup", vc.ID, int(groupID))
	return err
}

// DelGroup deletes a group from the VDC
// * groupID: The group ID.
func (vc *VdcController) DelGroup(groupID int) error {
	_, err := vc.c.Client.Call("one.vdc.delgroup", vc.ID, int(groupID))
	return err
}

// AddCluster adds a cluster to the VDC
// * zoneID: The Zone ID.
// * clusterID: The Cluster ID.
func (vc *VdcController) AddCluster(zoneID, clusterID int) error {
	_, err := vc.c.Client.Call("one.vdc.addcluster", vc.ID, int(zoneID), int(clusterID))
	return err
}

// DelCluster deletes a cluster from the VDC
// * zoneID: The Zone ID.
// * clusterID: The Cluster ID.
func (vc *VdcController) DelCluster(zoneID, clusterID int) error {
	_, err := vc.c.Client.Call("one.vdc.delcluster", vc.ID, int(zoneID), int(clusterID))
	return err
}

// AddHost adds a host to the VDC
// * zoneID: The Zone ID.
// * hostID: The Host ID.
func (vc *VdcController) AddHost(zoneID, hostID int) error {
	_, err := vc.c.Client.Call("one.vdc.addhost", vc.ID, int(zoneID), int(hostID))
	return err
}

// DelHost deletes a host from the VDC
// * zoneID: The Zone ID.
// * hostID: The Host ID.
func (vc *VdcController) DelHost(zoneID, hostID int) error {
	_, err := vc.c.Client.Call("one.vdc.delhost", vc.ID, int(zoneID), int(hostID))
	return err
}

// AddDatastore adds a datastore to the VDC
// * zoneID: The Zone ID.
// * dsID: The Datastore ID.
func (vc *VdcController) AddDatastore(zoneID, dsID int) error {
	_, err := vc.c.Client.Call("one.vdc.adddatastore", vc.ID, int(zoneID), int(dsID))
	return err
}

// DelDatastore deletes a datastore from the VDC
// * zoneID: The Zone ID.
// * dsID: The Datastore ID.
func (vc *VdcController) DelDatastore(zoneID, dsID int) error {
	_, err := vc.c.Client.Call("one.vdc.deldatastore", vc.ID, int(zoneID), int(dsID))
	return err
}

// AddVnet adds a vnet to the VDC
// * zoneID: The Zone ID.
// * vnetID: The Vnet ID.
func (vc *VdcController) AddVnet(zoneID, vnetID int) error {
	_, err := vc.c.Client.Call("one.vdc.addvnet", vc.ID, int(zoneID), int(vnetID))
	return err
}

// DelVnet deletes a vnet from the VDC
// * zoneID: The Zone ID.
// * vnetID: The Vnet ID.
func (vc *VdcController) DelVnet(zoneID, vnetID int) error {
	_, err := vc.c.Client.Call("one.vdc.delvnet", vc.ID, int(zoneID), int(vnetID))
	return err
}

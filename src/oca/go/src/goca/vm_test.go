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
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type VMSuite struct {
	templateID uint
	vmID       uint
	hostID     uint
}

var _ = Suite(&VMSuite{})

func (s *VMSuite) SetUpSuite(c *C) {
	// Create template
	tpl := NewDynamicTemplate()

	tpl.AddPair("NAME", GenName("VMSuite-template"))
	tpl.AddPair("CPU", 1)
	tpl.AddPair("MEMORY", "64")

	templateID, err := testCtrl.Templates().Create(tpl.String())
	c.Assert(err, IsNil)

	s.templateID = templateID

	s.hostID, _ = testCtrl.Hosts().Create("dummy-test", "dummy", "dummy", 0)

	tmpl := "TM_MAD=dummy\nDS_MAD=dummy"

	testCtrl.Datastore(1).Update(tmpl, 1)

	testCtrl.Datastore(0).Update(tmpl, 1)

}

func (s *VMSuite) SetUpTest(c *C) {
	vmID, err := testCtrl.Template(s.templateID).Instantiate("", true, nil, false)
	c.Assert(err, IsNil)
	s.vmID = vmID
}

func (s *VMSuite) TearDownTest(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.TerminateHard()
	if err != nil {
		err = vmC.RecoverDelete()
	}
	c.Assert(err, IsNil)
}

func (s *VMSuite) TearDownSuite(c *C) {
	testCtrl.Template(s.templateID).Delete()

	testCtrl.Host(s.hostID).Delete()
}

////////////////////////////////////////////////////////////////////////////////

func VMExpectState(c *C, vm *VM, state, lcmState string) func() bool {
	return func() bool {
		vm, err := testCtrl.VM(vm.ID).Info()
		if err != nil {
			return false
		}

		s, l, err := vm.StateString()
		if err != nil {
			return false
		}

		if lcmState != "" && l == lcmState {
			return true
		}

		if state != "" && s == state && lcmState == "" {
			return true
		}

		c.Logf("VM: %d. Expecting: %s/%s, Got: %s/%s", vm.ID, state, lcmState, s, l)

		return false
	}
}

////////////////////////////////////////////////////////////////////////////////

func (s *VMSuite) TestVMDeploy(c *C) {
	vmC := testCtrl.VM(s.vmID)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	err = vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)
}

func (s *VMSuite) TestVMHoldRelease(c *C) {
	vmC := testCtrl.VM(s.vmID)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "HOLD", "")), Equals, true)

	err = vmC.Release()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "PENDING", "")), Equals, true)
}

func (s *VMSuite) TestVMUpdate(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Update("A=B", 1)
	c.Assert(err, IsNil)

	vm, err := vmC.Info()
	c.Assert(err, IsNil)

	pair, err := vm.UserTemplate.Dynamic.GetPair("A")
	c.Assert(err, IsNil)
	c.Assert(pair.Value, Equals, "B")
}

// TODO: Hosts
// func (s *VMSuite) TestVMMigrate(c *C) {
// }
// func (s *VMSuite) TestVMLiveMigrate(c *C) {
// }

func (s *VMSuite) TestVMTerminate(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	err = vmC.Terminate()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "DONE", "")), Equals, true)
}

func (s *VMSuite) TestVMTerminateHard(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	vmC.TerminateHard()
	c.Assert(WaitResource(VMExpectState(c, vm, "DONE", "")), Equals, true)
}

func (s *VMSuite) TestVMStop(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	err = vmC.Stop()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "STOPPED", "")), Equals, true)
}

func (s *VMSuite) TestVMSuspend(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	err = vmC.Suspend()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "SUSPENDED", "")), Equals, true)
}

func (s *VMSuite) TestVMResume(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	err = vmC.Suspend()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "SUSPENDED", "")), Equals, true)

	err = vmC.Resume()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)
}

func (s *VMSuite) TestVMResize(c *C) {
	vmC := testCtrl.VM(s.vmID)
	err := vmC.Deploy(s.hostID, false, -1)
	c.Assert(err, IsNil)
	vm, err := vmC.Info()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	err = vmC.Poweroff()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "POWEROFF", "")), Equals, true)

	err = vmC.Resize("CPU=2.5\nMEMORY=512", false)
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "POWEROFF", "")), Equals, true)

	err = vmC.Resume()
	c.Assert(err, IsNil)
	c.Assert(WaitResource(VMExpectState(c, vm, "ACTIVE", "RUNNING")), Equals, true)

	vm, err = vmC.Info()
	c.Assert(err, IsNil)

	c.Assert(vm.Template.CPU, Equals, 2.5)

	c.Assert(vm.Template.Memory, Equals, 512)
}

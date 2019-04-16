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

import "fmt"

// WhoPool is a parameter to pool info methods allowing to so some filtering
type PoolWho int

const (
	// PoolWhoPrimaryGroup resources belonging to the user’s primary group.
	PoolWhoPrimaryGroup PoolWho = -4

	// PoolWhoMine to list resources that belong to the user that performs the
	// query.
	PoolWhoMine PoolWho = -3

	// PoolWhoAll to list all the resources seen by the user that performs the
	// query.
	PoolWhoAll PoolWho = -2

	// PoolWhoGroup to list all the resources that belong to the group that performs
	// the query.
	PoolWhoGroup PoolWho = -1
)

// Filterer is an interface for pool filters
type Filterer interface {
	SetUID(uid int) error
	SetVisibility(flag PoolWho)
	SetRange(start, end int)
	toArgs() []interface{}
}

// VMFilterer is an interface for VM pool filters
type VMFilterer interface {
	Filterer
	SetState(state VMState)
}

// Filter is a structure containing parameters for pool info methods
type Filter struct {
	who   int
	start int
	end   int
}

// ExtendedFilter allow extended filtering abilities (Only available on VMs for now)
type ExtendedFilter struct {
	Filter
	pair string // filter on KEY=VALUE
}

// VMFilter is filter for VM
type VMFilter struct {
	Filter
	state VMState
}

// VMExtendedFilter is extended filter for VM
type VMExtendedFilter struct {
	ExtendedFilter
	state VMState
}

// DocumentFilter is filter for VM
type DocumentFilter struct {
	Filter
	docType int
}

func NewFilter() Filter {
	return Filter{
		who:   int(PoolWhoMine),
		start: -1,
		end:   -1,
	}
}

func NewExtendedFilter() ExtendedFilter {
	return ExtendedFilter{
		Filter: NewFilter(),
	}
}

func NewVMFilter() VMFilter {
	return VMFilter{
		Filter: NewFilter(),
		state:  -1,
	}
}

func NewVMExtendedFilter() VMExtendedFilter {
	return VMExtendedFilter{
		ExtendedFilter: NewExtendedFilter(),
		state:          -1,
	}
}

func (f *Filter) toArgs() []interface{} {
	return []interface{}{f.who, f.start, f.end}
}

func (f *ExtendedFilter) toArgs() []interface{} {
	return []interface{}{f.who, f.start, f.end, f.pair}
}

func (f *VMFilter) toArgs() []interface{} {
	return []interface{}{f.who, f.start, f.end, f.state}
}

func (f *VMExtendedFilter) toArgs() []interface{} {
	return []interface{}{f.who, f.start, f.end, f.state, f.pair}
}

func (f *DocumentFilter) toArgs() []interface{} {
	return []interface{}{f.who, f.start, f.end, f.docType}
}

// SetUID enable filtering on a user ID filter.
// You have to choose between SetUID and SetVisibility methods.
func (f *Filter) SetUID(uid int) error {
	if uid < 0 {
		return fmt.Errorf("Filter.SetUID: parameters uid must be positive")
	}
	f.who = uid
	return nil
}

// SetVisibility enable filtering based on visibility from user/group used to connect to OpenNebula.
// You have to choose between SetUID and SetVisibility methods.
func (f *Filter) SetVisibility(flag PoolWho) {
	f.who = int(flag)
}

// TODO: see doc for pagination when values < -1
// SetRange enable filtering based on id range
func (f *Filter) SetRange(start, end int) {
	f.start = start
	f.end = end
}

// SetPair allow to filter based on a key, value pair
func (f *ExtendedFilter) SetPair(key, value string) {
	f.pair = fmt.Sprintf("%s=\"%s\"", key, value)
}

// SetState enable filtering based on id range
func (f *VMFilter) SetState(state VMState) {
	f.state = state
}

// SetState enable filtering based on id range
func (f *VMExtendedFilter) SetState(state VMState) {
	f.state = state
}

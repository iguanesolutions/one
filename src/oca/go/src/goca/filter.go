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
	"fmt"
	"strings"

	"github.com/OpenNebula/one/src/oca/go/src/goca/parameters"
)

// WhoPool is a parameter to pool info methods allowing to so some filtering
type WhoPool int

const (
	// PoolPrimaryGroup resources belonging to the user’s primary group.
	PoolPrimaryGroup WhoPool = -4

	// PoolMine to list resources that belong to the user that performs the
	// query.
	PoolMine WhoPool = -3

	// PoolAll to list all the resources seen by the user that performs the
	// query.
	PoolAll WhoPool = -2

	// PoolGroup to list all the resources that belong to the group that performs
	// the query.
	PoolGroup WhoPool = -1
)

type filter struct {
	args []int
	pair string
}

func newFilterDefault() *filter {
	return &filter{
		args: []int{parameters.PoolWhoMine, -1, -1},
	}
}

func newVMFilterDefault() *filter {
	return &filter{
		args: []int{parameters.PoolWhoMine, -1, -1, -1},
	}
}

func (f *filter) toArgs() []interface{} {
	return []interface{}{f.args[0], f.args[1], f.args[2]}
}

func (f *filter) toVMArgs() []interface{} {
	if len(f.pair) > 0 {
		return []interface{}{f.args[0], f.args[1], f.args[2], f.args[3], f.pair}
	}
	return []interface{}{f.args[0], f.args[1], f.args[2], f.args[3]}
}

type filterOption func(*filter)

func Who(w WhoPool) filterOption {
	return func(f *filter) {
		f.args[0] = int(w)
	}
}

func Range(start, end int) filterOption {
	return func(f *filter) {
		f.args[1] = start
		f.args[2] = end
	}
}

func State(s int) filterOption {
	return func(f *filter) {
		f.args[3] = -1
	}
}

func Pair(k string, v interface{}) filterOption {
	return func(f *filter) {
		f.pair = fmt.Sprintf("%s=%s", strings.ToUpper(k), fmt.Sprint(v))
	}
}

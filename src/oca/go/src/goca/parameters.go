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
	"errors"
	"fmt"
)

// UpdateType is a parameter to update methods indicating how to replace the template
type UpdateType int

const (
	// Replace to replace the whole template
	Replace UpdateType = 0

	// Merge to merge new template with existing one
	Merge UpdateType = 1
)

// PoolWho is a parameter to pool info methods allowing to so some filtering
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

type View struct {
	who int
	id  ViewRangeID
}

type ViewRangeID struct {
	start int
	end   int
}

// NewViewWho is a helper allowing you to only give the visibility, and optionally the range
func NewViewWho(who int, args ...int) (View, error) {
	var f View

	if f.who < int(PoolWhoPrimaryGroup) {
		return f, fmt.Errorf("bad Who value")
	}

	f.who = who

	switch len(args) {
	case 0:
		f.id = ViewRangeID{-1, -1}
	case 1:
		return f, errors.New("Wrong number of arguments: end of range is missing")
	case 2:
		f.id = ViewRangeID{args[0], args[1]}
	default:
		return f, errors.New("Wrong number of arguments: too much arguments")
	}

	return f, nil
}

// NewView is a function returning a View from a list of variable parameters
func NewView(args ...int) (View, error) {
	if len(args) > 1 {
		return NewViewWho(args[0], args...)
	}

	return View{
		int(PoolWhoMine),
		ViewRangeID{-1, -1},
	}, nil

}

// SetUID enable filtering on a user ID filter.
// You have to choose between SetUID and SetVisibility methods.
func (f *View) SetUID(uid int) error {
	if uid < 0 {
		return fmt.Errorf("View.SetUID: parameters uid must be positive")
	}
	f.who = uid
	return nil
}

// SetVisibility enable filtering based on visibility from user/group used to connect to OpenNebula.
// You have to choose between SetUID and SetVisibility methods.
func (f *View) SetVisibility(flag PoolWho) error {
	switch flag {
	case -4, -3, -2, -1:
	default:
		return fmt.Errorf("View.SetVisibility: bad parameter value, should be one of PoolWho values")
	}
	f.who = int(flag)

	return nil
}

func argsToInt(length int, args ...interface{}) ([]int, error) {

	min := length
	if len(args) < length {
		min = len(args)
	}

	argsI := make([]int, 0, min)
	for i := 0; i < min; i++ {
		argI, ok := args[i].(int)
		if !ok {
			return argsI, fmt.Errorf("wrong type, needs an int")
		}
		argsI = append(argsI, argI)
	}

	return argsI, nil
}

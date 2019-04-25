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
	"fmt"
	"io"
	"strconv"
	"strings"
)

// DynamicTemplate represents an OpenNebula syntax template
type DynamicTemplate struct {
	elements []DynamicTemplateElement
}

// DynamicTemplateElement is an interface that must implement the String
// function
type DynamicTemplateElement interface {
	String() string
	Key() string
}

// DynamicTemplatePair is a key / value pair
type DynamicTemplatePair struct {
	key   string
	Value string
}

// DynamicTemplateVector contains an array of keyvalue pairs
type DynamicTemplateVector struct {
	key   string
	pairs []DynamicTemplatePair
}

// Key return the pair key
func (t *DynamicTemplatePair) Key() string { return t.key }

// Key return the vector key
func (t *DynamicTemplateVector) Key() string { return t.key }

// String prints the DynamicTemplate in OpenNebula syntax
func (t *DynamicTemplate) String() string {
	s := ""
	endToken := "\n"

	for i, element := range t.elements {
		if i == len(t.elements)-1 {
			endToken = ""
		}
		s += element.String() + endToken
	}

	return s
}

// String prints a DynamicTemplatePair in OpenNebula syntax
func (t *DynamicTemplatePair) String() string {
	return fmt.Sprintf("%s=\"%s\"", t.key, t.Value)
}

func (t *DynamicTemplateVector) String() string {
	s := fmt.Sprintf("%s=[\n", strings.ToUpper(t.Key()))

	endToken := ",\n"
	for i, pair := range t.pairs {
		if i == len(t.pairs)-1 {
			endToken = ""
		}

		s += fmt.Sprintf("    %s%s", pair.String(), endToken)

	}
	s += " ]"

	return s
}

// Exists return true it the key exists
func (t *DynamicTemplate) Exists(key string) bool {
	for i := 0; i < len(t.elements); i++ {
		if t.elements[i].Key() != key {
			continue
		}
		return true
	}
	return false
}

// Exists return true it the key exists
func (t *DynamicTemplateVector) Exists(key string) bool {
	for i := 0; i < len(t.pairs); i++ {
		if t.pairs[i].key != key {
			continue
		}
		return true
	}
	return false
}

// GetPairs retrieve a list of pairs by key.
func (t *DynamicTemplate) GetPairs(key string) []*DynamicTemplatePair {
	pairs := make([]*DynamicTemplatePair, 0)
	for _, e := range t.elements {
		pair, ok := e.(*DynamicTemplatePair)
		if !ok || pair.key != key {
			continue
		}
		pairs = append(pairs, pair)
	}
	return pairs
}

// GetPairs retrieve a list of pairs by key.
func (t *DynamicTemplateVector) GetPairs(key string) []*DynamicTemplatePair {
	pairs := make([]*DynamicTemplatePair, 0)
	for i := 0; i < len(t.pairs); i++ {
		if t.pairs[i].key != key {
			continue
		}
		pairs = append(pairs, &t.pairs[i])
	}
	return pairs
}

// GetVectors retrieve a list of vectors by key.
func (t *DynamicTemplate) GetVectors(key string) []*DynamicTemplateVector {
	vectors := make([]*DynamicTemplateVector, 0)
	for _, e := range t.elements {
		vec, ok := e.(*DynamicTemplateVector)
		if !ok || vec.key != key {
			continue
		}
		vectors = append(vectors, vec)
	}
	return vectors
}

// Get retrieve a unique pair by key. Fail if not found or several instances
func (t *DynamicTemplate) GetPair(key string) (*DynamicTemplatePair, error) {
	pairs := t.GetPairs(string(key))
	if len(pairs) == 0 {
		return nil, fmt.Errorf("Get: tag %s not found", key)
	} else if len(pairs) > 1 {
		return nil, fmt.Errorf("Get: multiple entries with key %s", key)
	}
	return pairs[0], nil
}

// Get retrieve a unique pair by key. Fail if not found or several instances
func (t *DynamicTemplateVector) GetPair(key string) (*DynamicTemplatePair, error) {
	pairs := t.GetPairs(string(key))
	if len(pairs) == 0 {
		return nil, fmt.Errorf("Get: tag %s not found", key)
	} else if len(pairs) > 1 {
		return nil, fmt.Errorf("Get: multiple entries with key %s", key)
	}
	return pairs[0], nil
}

// GetVector retrieve a unique vector by key. Fail if not found or several instances
func (t *DynamicTemplate) GetVector(key string) (*DynamicTemplateVector, error) {
	vectors := t.GetVectors(string(key))
	if len(vectors) == 0 {
		return nil, fmt.Errorf("Get: tag %s not found", key)
	} else if len(vectors) > 1 {
		return nil, fmt.Errorf("Get: multiple entries with key %s", key)
	}
	return vectors[0], nil
}

func (t *DynamicTemplate) getValue(key string) (string, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return "", err
	}
	return pair.Value, nil
}

func (t *DynamicTemplateVector) getValue(key string) (string, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return "", err
	}
	return pair.Value, nil
}

func (t *DynamicTemplate) getValueFromVec(vecKey, key string) (string, error) {
	vector, err := t.GetVector(vecKey)
	if err != nil {
		return "", err
	}
	pair, err := vector.GetPair(key)
	if err != nil {
		return "", err
	}
	return pair.Value, nil
}

func (t *DynamicTemplate) getID(key string) (uint, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(pair.Value, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (t *DynamicTemplateVector) getID(key string) (uint, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(pair.Value, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// Dynamic template building

// NewDynamicTemplate returns a new DynamicTemplate object
func NewDynamicTemplate() *DynamicTemplate {
	return &DynamicTemplate{}
}

// AddVector creates a new vector in the template
func (t *DynamicTemplate) AddVector(key string) *DynamicTemplateVector {
	vector := &DynamicTemplateVector{key: key}

	t.elements = append(t.elements, vector)
	return vector
}

// AddPair adds a new pair to a DynamicTemplate objects
func (t *DynamicTemplate) AddPair(key string, v interface{}) error {
	var val string

	switch v := v.(type) {
	default:
		return fmt.Errorf("Unexpected type")
	case int, uint:
		val = fmt.Sprintf("%d", v)
	case string:
		val = v
	}

	pair := &DynamicTemplatePair{strings.ToUpper(key), val}
	t.elements = append(t.elements, pair)

	return nil
}

// AddPair adds a new pair to a DynamicTemplate
func (t *DynamicTemplateVector) AddPair(key string, v interface{}) error {
	var val string

	switch v := v.(type) {
	default:
		return fmt.Errorf("Unexpected type")
	case int, uint:
		val = fmt.Sprintf("%d", v)
	case string:
		val = v
	}

	pair := DynamicTemplatePair{strings.ToUpper(key), val}
	t.pairs = append(t.pairs, pair)

	return nil
}

func (t *DynamicTemplate) addPairToVec(vecKey, key string, value interface{}) error {
	var vector *DynamicTemplateVector

	vectors := t.GetVectors(string(key))
	switch len(vectors) {
	case 0:
		vector = t.AddVector(vecKey)
	case 1:
		vector = vectors[0]
	default:
		return fmt.Errorf("Get: multiple entries with key %s", key)
	}

	return vector.AddPair(string(key), value)
}

// Del remove an element from DynamicTemplate objects
func (t *DynamicTemplate) Del(key string) {
	for i := 0; i < len(t.elements); i++ {
		if t.elements[i].Key() != key {
			continue
		}
		t.elements = append(t.elements[:i], t.elements[i+1:]...)
	}
}

// Del remove a pair from DynamicTemplate
func (t *DynamicTemplateVector) Del(key string) {
	for i := 0; i < len(t.pairs); i++ {
		if t.pairs[i].Key() != key {
			continue
		}
		t.pairs = append(t.pairs[:i], t.pairs[i+1:]...)
	}
}

// Dynamic template parsing

// xmlTag contains the tag informations
type xmlTag struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

// UnmarshalXML parse dynamically templates
func (t *DynamicTemplate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var tplEl DynamicTemplateElement
	var tplVec *DynamicTemplateVector
	var e xmlTag

	for {

		// Look at next element
		token, err := d.Token()
		if err != nil {
			if err.Error() == io.EOF.Error() {
				break
			}
			return err
		}

		// Check type to handle the nesting
		switch token.(type) {
		case xml.StartElement:
			// It's a vector

			// Create it at first time
			if tplVec == nil {
				tplVec = &DynamicTemplateVector{
					key:   start.Name.Local,
					pairs: make([]DynamicTemplatePair, 0),
				}
			}

			// Decode pair and add it to the vec
			startEl, _ := token.(xml.StartElement)
			err := d.DecodeElement(&e, &startEl)
			if err != nil {
				return err
			}
			tplPair := DynamicTemplatePair{key: e.XMLName.Local, Value: e.Content}
			tplVec.pairs = append(tplVec.pairs, tplPair)

			tplEl = tplVec
		case xml.CharData:
			// It's a Pair

			val, _ := token.(xml.CharData)
			tplEl = &DynamicTemplatePair{key: start.Name.Local, Value: string(val)}
		}
	}
	t.elements = append(t.elements, tplEl)

	return nil
}

// UnmarshalXML parse dynamically templates vector
func (t *DynamicTemplateVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t.key = start.Name.Local
	t.pairs = make([]DynamicTemplatePair, 0)

	var e xmlTag
	for {
		// Look at next element
		token, err := d.Token()
		if err != nil {
			if err.Error() == io.EOF.Error() {
				break
			}
			return err
		}

		// Decode pair from start element
		startEl, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		err = d.DecodeElement(&e, &startEl)
		if err != nil {
			return err
		}
		t.pairs = append(t.pairs, DynamicTemplatePair{key: e.XMLName.Local, Value: e.Content})
	}

	return nil
}

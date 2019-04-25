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
	"sort"
	"strconv"
	"strings"
)

// DynamicTemplate represents an OpenNebula syntax template
type DynamicTemplate struct {
	elements []TemplateElement
}

// TemplateElement is an interface that must implement the String
// function
type TemplateElement interface {
	String() string
	Key() string
}

// TemplatePair is a key / value pair
type TemplatePair struct {
	key   string
	Value string
}

// DynamicTemplateVector contains an array of keyvalue pairs
type TemplateVector struct {
	key   string
	pairs []TemplatePair
}

type dynamicTemplateAny struct {
	DynamicTemplate
}

// Key return the pair key
func (t *TemplatePair) Key() string { return t.key }

// Key return the vector key
func (t *TemplateVector) Key() string { return t.key }

// String prints the DynamicTemplate in OpenNebula syntax
func (t *DynamicTemplate) String() string {
	var s strings.Builder
	endToken := "\n"

	for i, element := range t.elements {
		if i == len(t.elements)-1 {
			endToken = ""
		}
		s.WriteString(element.String() + endToken)
	}

	return s.String()
}

// String prints a DynamicTemplatePair in OpenNebula syntax
func (t *TemplatePair) String() string {
	return fmt.Sprintf("%s=\"%s\"", t.key, t.Value)
}

func (t *TemplateVector) String() string {
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

func (t *TemplateVector) getPairIdx(key string) (int, error) {

	// binary search on pairs
	idx := sort.Search(len(t.pairs), func(i int) bool {
		return key <= t.pairs[i].key
	})

	// key not found
	if idx == len(t.pairs) || t.pairs[idx].key != key {
		return -1, fmt.Errorf("key %s not found", key)
	}

	return idx, nil
}

// GetPair retrieve a pair by key
func (t *TemplateVector) GetPair(key string) (*TemplatePair, error) {

	idx, err := t.getPairIdx(key)
	if err != nil {
		return nil, err
	}

	return &t.pairs[idx], nil
}

// GetPair retrieve a pair by key
func (t *DynamicTemplate) GetPair(key string) (*TemplatePair, error) {

	// binary search on elements
	idx := sort.Search(len(t.elements), func(i int) bool {
		return key <= t.elements[i].Key()
	})

	// key not found
	if idx == len(t.elements) || t.elements[idx].Key() != key {
		return nil, fmt.Errorf("GetPair: key %s not found", key)
	}

	// if it's a pair, return it
	pair, ok := t.elements[idx].(*TemplatePair)
	if ok {
		return pair, nil
	}

	// look for a pair with this key
	for i := idx + 1; i < len(t.elements); i++ {

		pair, ok = t.elements[i].(*TemplatePair)
		if !ok {
			continue
		}

		if pair.key == key {
			return pair, nil
		}
	}

	return nil, fmt.Errorf("GetPair: key %s not found", key)
}

// GetVectors retrieve slice of vectors by key
func (t *DynamicTemplate) GetVectors(key string) []*TemplateVector {
	vectors := make([]*TemplateVector, 0, 1)

	idx := sort.Search(len(t.elements), func(i int) bool {
		return key <= t.elements[i].Key()
	})

	// key not found
	if idx == len(t.elements) || t.elements[idx].Key() != key {
		return nil
	}

	// append if it's a vector
	vec, ok := t.elements[idx].(*TemplateVector)
	if ok {
		vectors = append(vectors, vec)
	}

	// append other vector with same key
	for i := idx + 1; i < len(t.elements); i++ {

		vec, ok := t.elements[i].(*TemplateVector)
		if !ok {
			continue
		}

		if vec.key == key {
			vectors = append(vectors, vec)
		}
	}

	return vectors
}

// GetVector retrieve a vector by key
func (t *DynamicTemplate) GetVector(key string) (*TemplateVector, error) {

	// binary search on elements
	idx := sort.Search(len(t.elements), func(i int) bool {
		return key <= t.elements[i].Key()
	})

	// key not found
	if idx == len(t.elements) || t.elements[idx].Key() != key {
		return nil, fmt.Errorf("GetVector: key %s not found", key)
	}

	// if it's a vector, return it
	vec, ok := t.elements[idx].(*TemplateVector)
	if ok {
		return vec, nil
	}

	// look for a vector with this key
	for i := idx + 1; i < len(t.elements); i++ {

		vec, ok = t.elements[i].(*TemplateVector)
		if !ok {
			continue
		}

		if vec.key == key {
			return vec, nil
		}
	}

	return nil, fmt.Errorf("GetVector: key %s not found", key)
}

// GetStr allow to retrieve the value of a pair
func (t *DynamicTemplate) GetStr(key string) (string, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return "", err
	}
	return pair.Value, nil
}

// GetStr allow to retrieve the value of a pair
func (t *TemplateVector) GetStr(key string) (string, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return "", err
	}
	return pair.Value, nil
}

// GetInt get a pair, convert it's value to int and return it
func (t *DynamicTemplate) GetInt(key string) (int, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(pair.Value, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetInt get a pair, convert it's value to int and return it
func (t *TemplateVector) GetInt(key string) (int, error) {
	pair, err := t.GetPair(key)
	if err != nil {
		return -1, err
	}
	id, err := strconv.ParseInt(pair.Value, 10, 0)
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// helper to get a pair value from inside of a vector
func (t *DynamicTemplate) getStrFromVec(vecKey, key string) (string, error) {
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

// look if template match a whole pair
func (t *DynamicTemplate) matchPair(key, value string) bool {

	// keep vectors to look into it later
	tmpVecs := make([]*TemplateVector, 0, 1)

	idx := sort.Search(len(t.elements), func(i int) bool {
		vec, ok := t.elements[i].(*TemplateVector)
		if ok {
			tmpVecs = append(tmpVecs, vec)
		}

		return key <= t.elements[i].Key()
	})

	// a element was found with this key, look for a pair
	if idx < len(t.elements) {

		for i := idx; i < len(t.elements); i++ {

			pair, ok := t.elements[i].(*TemplatePair)
			if !ok {
				continue
			}

			if pair.key == key && pair.Value == value {
				return true
			}
		}
	}

	// look inside of each vectors for the pair
	for _, vec := range tmpVecs {

		idx, err := vec.getPairIdx(key)
		if err != nil {
			continue
		}

		for i := idx; i < len(vec.pairs); i++ {
			if vec.pairs[i].key == key && vec.pairs[i].Value == value {
				return true
			}
		}
	}

	return false
}

// Dynamic template building

// NewDynamicTemplate returns a new DynamicTemplate object
func NewDynamicTemplate() *DynamicTemplate {
	return &DynamicTemplate{}
}

// AddVector creates a new vector in the template
func (t *DynamicTemplate) AddVector(key string) *TemplateVector {
	vector := &TemplateVector{key: key}

	t.elements = append(t.elements, vector)
	return vector
}

// AddPair adds a new pair to a DynamicTemplate objects
func (t *DynamicTemplate) AddPair(key string, v interface{}) error {
	var val string

	switch v := v.(type) {
	default:
		return fmt.Errorf("AddPair: Unexpected type")
	case int, uint:
		val = fmt.Sprintf("%d", v)
	case string:
		val = v
	}

	pair := &TemplatePair{strings.ToUpper(key), val}
	t.elements = append(t.elements, pair)

	return nil
}

// AddPair adds a new pair to a DynamicTemplate
func (t *TemplateVector) AddPair(key string, v interface{}) error {
	var val string

	switch v := v.(type) {
	default:
		return fmt.Errorf("AddPair: Unexpected type")
	case int, uint:
		val = fmt.Sprintf("%d", v)
	case string:
		val = v
	}

	pair := TemplatePair{strings.ToUpper(key), val}
	t.pairs = append(t.pairs, pair)

	return nil
}

func (t *DynamicTemplate) addPairToVec(vecKey, key string, value interface{}) error {
	var vector *TemplateVector

	vectors := t.GetVectors(string(key))
	switch len(vectors) {
	case 0:
		vector = t.AddVector(vecKey)
	case 1:
		vector = vectors[0]
	default:
		return fmt.Errorf("Can't add pair to vector: multiple entries with key %s", key)
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
func (t *TemplateVector) Del(key string) {
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

// UnmarshalXML parse dynamically a template under the ",any" tag. It's used when we mix statically parsed Template with dynamic values.
func (t *dynamicTemplateAny) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return unmarshalTemplateElement(d, start, &t.DynamicTemplate)
}

// UnmarshalXML parse dynamically a template under the "TEMPLATE" tag
func (t *DynamicTemplate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {

		// Retrieve next token to start in unmarshalTemplateElement at the same position
		// than with dynamicTemplateAny.UnmarshalXML
		token, err := d.Token()
		if err != nil {
			if err.Error() == io.EOF.Error() {
				break
			}
			return err
		}

		// It's a start tag of a pair or a vector, we'll see this inside of unmarshalXMLTemplate
		switch token.(type) {
		case xml.StartElement:
			startTag, _ := token.(xml.StartElement)
			err := unmarshalTemplateElement(d, startTag, t)
			if err != nil {
				return err
			}
		case xml.EndElement:
			break
		}
	}

	return nil
}

// UnmarshalXML parse dynamically a template vector
func (t *TemplateVector) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	// In case we unmarshall a simple vector, then we need to initialize it
	if t.pairs == nil {
		t.key = start.Name.Local
		t.pairs = make([]TemplatePair, 0, 2)
	}

	var e xmlTag
	for {

		// Retrieve the next token
		token, err := d.Token()
		if err != nil {
			if err.Error() == io.EOF.Error() {
				return nil
			}
			return err
		}

		// add a pair on a start element
		startEl, ok := token.(xml.StartElement)
		if ok {
			err = d.DecodeElement(&e, &startEl)
			if err != nil {
				return err
			}
			t.AddPair(e.XMLName.Local, e.Content)
		}

		// End element are consumed implicitly

	}

	return nil
}

// unmarshal only one element, a pair, or a vector
func unmarshalTemplateElement(d *xml.Decoder, start xml.StartElement, t *DynamicTemplate) error {
	var vec *TemplateVector

	// Retrieve the next token
	token, err := d.Token()
	if err != nil {
		if err.Error() == io.EOF.Error() {
			return nil
		}
		return err
	}
	isVector := false

	switch token.(type) {
	case xml.StartElement:

		// It's a start tag: vector
		isVector = true

	case xml.CharData:

		// As we call Token method again, we must save the chardata buffer in case we need it later
		tokenBuf := xml.CopyToken(token)

		// need to look at a third token to distinguish between pair and vector
		tokenNext, err := d.Token()
		if err != nil {
			if err.Error() == io.EOF.Error() {
				return nil
			}
			return err
		}

		switch tokenNext.(type) {
		case xml.StartElement:

			isVector = true
			token = tokenNext

		case xml.EndElement:

			// Adds a pair

			val, ok := tokenBuf.(xml.CharData)
			if !ok {
				return fmt.Errorf("unmarshalTemplateElement UnmarshalXML: chardata element attended")
			}
			t.AddPair(start.Name.Local, string(val))

		}

	}

	if isVector {

		startVec, ok := token.(xml.StartElement)
		if !ok {
			return fmt.Errorf("unmarshalTemplateElement UnmarshalXML: start element attended")
		}

		vec = t.AddVector(start.Name.Local)

		// Add first element from token
		var e xmlTag
		err := d.DecodeElement(&e, &startVec)
		if err != nil {
			return err
		}
		vec.AddPair(e.XMLName.Local, e.Content)

		// unmarshal the rest of the vector
		err = vec.UnmarshalXML(d, startVec)
		if err != nil {
			return fmt.Errorf("unmarshalTemplateElement vector UnmarshalXML: %s", err)
		}

	}

	return nil
}

// Args needed by at a set of entities
const (
	DescriptionK string = "DESCRIPTION"
	NameK        string = "NAME"
	TypeK        string = "TYPE"
)

func (t *DynamicTemplate) SetName(desc string) error {
	t.Del(NameK)
	return t.AddPair(NameK, desc)
}

func (t *DynamicTemplate) SetDescription(desc string) error {
	t.Del(DescriptionK)
	return t.AddPair(DescriptionK, desc)
}

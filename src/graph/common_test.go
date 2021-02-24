/*
 * Copyright 2020 Paul Tatham <paul@nextmetaphor.io>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package graph

import (
	"fmt"
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func Test_getDefinitionCypherString(t *testing.T) {

	t.Run("ZeroFields", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{})
		assert.Equal(t, fmt.Sprintf(mergeCypherPrefix, "class1"), cypher)
	})

	t.Run("OneField", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{"name1": "value1"})
		assert.Equal(t, fmt.Sprintf(mergeCypherPrefix, "class1")+" SET "+fmt.Sprintf(mergeCypherField,
			"name1", "name1"), cypher)
	})

	t.Run("TwoFields", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{
			"name1": "value1", "name2": "value2"})

		// can't guarantee order here, so match with a regexp
		fieldRegExp := "n\\.name1=\\$name1|n\\.name2=\\$name2"
		cypherRegExp := regexp.QuoteMeta(fmt.Sprintf(mergeCypherPrefix, "class1")) +
			" SET (" + fieldRegExp + "),(" + fieldRegExp + ")"

		re, err := regexp.MatchString(cypherRegExp, cypher)
		assert.Nil(t, err)
		assert.True(t, re)

		assert.True(t, strings.Contains(cypher, "n.name1=$name1"))
		assert.True(t, strings.Contains(cypher, "n.name2=$name2"))
	})

	t.Run("ThreeFields", func(t *testing.T) {

		cypher := getDefinitionCypherString("class1", definition.Fields{
			"name1": "value1", "name2": "value2", "name3": "value3"})

		// can't guarantee order here, so match with a regexp
		fieldRegExp := "n\\.name1=\\$name1|n\\.name2=\\$name2|n\\.name3=\\$name3"
		cypherRegExp := regexp.QuoteMeta(fmt.Sprintf(mergeCypherPrefix, "class1")) +
			" SET (" + fieldRegExp + "),(" + fieldRegExp + "),(" + fieldRegExp + ")"

		re, err := regexp.MatchString(cypherRegExp, cypher)
		assert.Nil(t, err)
		assert.True(t, re)

		assert.True(t, strings.Contains(cypher, "n.name1=$name1"))
		assert.True(t, strings.Contains(cypher, "n.name2=$name2"))
		assert.True(t, strings.Contains(cypher, "n.name3=$name3"))
	})

}

func Test_getDefinitionID(t *testing.T) {
	const (
		pID = "parentID"
		dID = "definitionID"
	)

	trueVar := true
	falseVar := false

	t.Run("SpecTrueDfnNil", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: nil},
			definition.Specification{ PrependParentID: &trueVar},
		)

		assert.Equal(t, pID+dID, definitionID)
	})

	t.Run("SpecTrueDfnFalse", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &falseVar},
			definition.Specification{ PrependParentID: &trueVar},
		)

		assert.Equal(t, dID, definitionID)
	})

	t.Run("SpecTrueDfnTrue", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &trueVar},
			definition.Specification{ PrependParentID: &trueVar},
		)

		assert.Equal(t, pID+dID, definitionID)
	})

	t.Run("SpecFalseDfnNil", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: nil},
			definition.Specification{ PrependParentID: &falseVar},
		)

		assert.Equal(t, dID, definitionID)
	})

	t.Run("SpecFalseDfnFalse", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &falseVar},
			definition.Specification{ PrependParentID: &falseVar},
		)

		assert.Equal(t, dID, definitionID)
	})

	t.Run("SpecFalseDfnTrue", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &trueVar},
			definition.Specification{ PrependParentID: &falseVar},
		)

		assert.Equal(t, pID+dID, definitionID)
	})

	t.Run("SpecNilDfnNil", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: nil},
			definition.Specification{ PrependParentID: nil},
		)

		assert.Equal(t, dID, definitionID)
	})

	t.Run("SpecNilDfnFalse", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &falseVar},
			definition.Specification{ PrependParentID: nil},
		)

		assert.Equal(t, dID, definitionID)
	})

	t.Run("SpecNilDfnTrue", func(t *testing.T) {
		definitionID := getDefinitionID(
			pID,
			dID,
			definition.Definition{PrependParentID: &trueVar},
			definition.Specification{ PrependParentID: nil},
		)

		assert.Equal(t, pID+dID, definitionID)
	})

}
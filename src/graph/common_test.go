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

func Test_getEdgeCypherString(t *testing.T) {

	t.Run("FromOnly", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: true,
			RelationshipTo:   false,
			Fields:           nil,
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "<", ref.Relationship, "", ""), cypher)
	})

	t.Run("ToOnly", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: false,
			RelationshipTo:   true,
			Fields:           nil,
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "", ref.Relationship, "", ">"), cypher)
	})

	t.Run("NeitherFromTo", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: false,
			RelationshipTo:   false,
			Fields:           nil,
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "", ref.Relationship, "", ""), cypher)
	})

	t.Run("BothFromTo", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: true,
			RelationshipTo:   true,
			Fields:           nil,
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "<", ref.Relationship, "", ">"), cypher)
	})

	t.Run("BothFromToEmptyFields", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: true,
			RelationshipTo:   true,
			Fields:           definition.Fields{},
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "<", ref.Relationship, "", ">"), cypher)
	})

	t.Run("BothFromToSingleField", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: true,
			RelationshipTo:   true,
			Fields:           definition.Fields{"name1": "value1"},
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)
		assert.Equal(t, fmt.Sprintf(edgeCypher, "class2", "ID2", ref.Class, ref.ID, "<", ref.Relationship, "{name1: \"value1\"}", ">"), cypher)
	})

	t.Run("BothFromToMultipleFields", func(t *testing.T) {
		ref := definition.Reference{
			Class:            "class1",
			ID:               "ID1",
			Relationship:     "IS_A",
			RelationshipFrom: true,
			RelationshipTo:   true,
			Fields:           definition.Fields{"name1": "value1", "name2": "value2", "name3": "value3"},
		}

		cypher := getEdgeCypherString("class2", "ID2", ref)

		assert.Contains(t, cypher, "MATCH (n1:class2 {ID:\"ID2\"})")
		assert.Contains(t, cypher, "MATCH (n2:class1 {ID:\"ID1\"})")
		assert.Contains(t, cypher, "MERGE (n1)<-[:IS_A {")
		assert.Contains(t, cypher, "}]->(n2);")

		assert.Contains(t, cypher, "name1: \"value1\"")
		assert.Contains(t, cypher, "name2: \"value2\"")
		assert.Contains(t, cypher, "name3: \"value3\"")
	})

}

func Test_AddSpecification(t *testing.T) {
	t.Run("AddSpecification", func(t *testing.T) {
		bc := NewBatchConfig()
		spec := definition.Specification{
			Class: "Class1",
			Definitions: map[string]definition.Definition{
				"ID1": {
					Fields: definition.Fields{"name": "value1"},
					SubDefinitions: map[string]definition.Specification{
						"SubRel": {
							Class: "Class2",
							Definitions: map[string]definition.Definition{
								"ID2": {
									Fields: definition.Fields{"name": "value2"},
								},
							},
						},
					},
				},
			},
			References: []definition.Reference{
				{
					Class:        "ClassRef",
					ID:           "IDRef",
					Relationship: "REL",
				},
			},
		}
		bc.AddSpecification(spec, nil)

		assert.Len(t, bc.Nodes["Class1"], 1)
		assert.Equal(t, "ID1", bc.Nodes["Class1"][0]["ID"])
		assert.Equal(t, "value1", bc.Nodes["Class1"][0]["name"])

		assert.Len(t, bc.Nodes["Class2"], 1)
		assert.Equal(t, "ID2", bc.Nodes["Class2"][0]["ID"])
		assert.Equal(t, "value2", bc.Nodes["Class2"][0]["name"])

		// 1 from spec.References, 1 from parentReference (ID2 -> ID1)
		assert.Len(t, bc.Edges["Class1"], 1)
		assert.Equal(t, "ID1", bc.Edges["Class1"][0].ID1)
		assert.Equal(t, "IDRef", bc.Edges["Class1"][0].ID2)

		assert.Len(t, bc.Edges["Class2"], 1)
		assert.Equal(t, "ID2", bc.Edges["Class2"][0].ID1)
		assert.Equal(t, "ID1", bc.Edges["Class2"][0].ID2)
		assert.Equal(t, "SubRel", bc.Edges["Class2"][0].Relationship)
	})
}

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

package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getOrderClause(t *testing.T) {
	t.Run("SingleOrderClause", func(t *testing.T) {
		clause := getOrderClause(ClassFieldSelector{
			Class:        "class1",
			Fields:       nil,
			Relationship: "",
			OrderFields:  []string{"field1"},
		})

		assert.Equal(t, "class1.field1", clause)
	})

	t.Run("MultipleOrderClause", func(t *testing.T) {
		clause := getOrderClause(ClassFieldSelector{
			Class:        "class2",
			Fields:       nil,
			Relationship: "",
			OrderFields:  []string{"field1", "field2"},
		})

		assert.Equal(t, "class2.field1,class2.field2", clause)
	})
	t.Run("NilOrderClause", func(t *testing.T) {
		clause := getOrderClause(ClassFieldSelector{
			Class:        "class3",
			Fields:       nil,
			Relationship: "",
			OrderFields:  nil,
		})

		assert.Equal(t, "", clause)
	})
	t.Run("EmptyOrderClause", func(t *testing.T) {
		clause := getOrderClause(ClassFieldSelector{
			Class:        "class4",
			Fields:       nil,
			Relationship: "",
			OrderFields:  []string{},
		})

		assert.Equal(t, "", clause)
	})
}

func Test_getCypherForSelector(t *testing.T) {
	t.Run("ParentClass", func(t *testing.T) {
		cypher := getCypherForSelector("parent", "parentID",
			ClassFieldSelector{
				Class:        "child",
				Fields:       nil,
				Relationship: "inherits_from",
				OrderFields:  []string{"field1", "field2"},
			})

		assert.Equal(t, "match (child:child)-[:inherits_from]-(parent:parent {ID:\"parentID\"}) return child "+
			"order by child.field1,child.field2", cypher)
	})
	t.Run("NoParentClass", func(t *testing.T) {
		cypher := getCypherForSelector("", "",
			ClassFieldSelector{
				Class:        "myclass",
				Fields:       nil,
				Relationship: "",
				OrderFields:  []string{"field1", "field2"},
			})

		assert.Equal(t, "match (myclass:myclass) return myclass order by myclass.field1,myclass.field2", cypher)
	})
}

func Test_loadTemplateConf(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		tc, err := loadTemplateConf("_test/TemplateSection_invalid.yaml")
		assert.NotNil(t, err)
		assert.Nil(t, tc)
	})

	t.Run("Minimal_Valid", func(t *testing.T) {
		tc, err := loadTemplateConf("_test/TemplateSection_minimal_valid.yaml")
		assert.Nil(t, err)
		assert.Equal(t, &TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class1",
				Fields:       []string{"field1", "field2", "field3"},
				Relationship: "",
				OrderFields:  []string{"field1", "field3"},
			},
			AggregateClasses:  nil,
			CompositeSections: nil,
		}, tc)
	})

	t.Run("AggregateOnly_Valid", func(t *testing.T) {
		tc, err := loadTemplateConf("_test/TemplateSection_aggregate_only_valid.yaml")
		assert.Nil(t, err)
		assert.Equal(t, &TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class1",
				Fields:       []string{"field1", "field2", "field3"},
				Relationship: "",
				OrderFields:  []string{"field1", "field3"},
			},
			AggregateClasses: []ClassFieldSelector{
				{
					Class:        "departmentClass",
					Fields:       []string{"d_field1", "d_field2", "d_field3"},
					Relationship: "belongs_to",
					OrderFields:  []string{"d_field1", "d_field3"},
				},
				{
					Class:        "clubClass",
					Fields:       []string{"c_field1", "c_field2"},
					Relationship: "member_of",
					OrderFields:  []string{"c_field1"},
				},
			},
			CompositeSections: nil,
		}, tc)
	})

	t.Run("Aggregate_Composite_Valid", func(t *testing.T) {
		tc, err := loadTemplateConf("_test/TemplateSection_aggregate_composite_valid.yaml")
		assert.Nil(t, err)
		assert.Equal(t, &TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class1",
				Fields:       []string{"c1_field1", "c1_field2", "c1_field3"},
				Relationship: "",
				OrderFields:  []string{"c1_field1", "c1_field3"},
			},
			AggregateClasses: []ClassFieldSelector{
				{
					Class:        "departmentClass1",
					Fields:       []string{"d1_field1", "d1_field2", "d1_field3"},
					Relationship: "belongs_to_1",
					OrderFields:  []string{"d1_field1", "d1_field3"},
				},
				{
					Class:        "clubClass1",
					Fields:       []string{"cl1_field1", "cl1_field2"},
					Relationship: "member_of_1",
					OrderFields:  []string{"cl1_field1"},
				},
			},
			CompositeSections: []TemplateSection{{
				SectionClass: ClassFieldSelector{
					Class:        "class2",
					Fields:       []string{"c2_field1", "c2_field2", "c2_field3"},
					Relationship: "",
					OrderFields:  []string{"c2_field1", "c2_field3"},
				},
				AggregateClasses: []ClassFieldSelector{
					{
						Class:        "departmentClass2",
						Fields:       []string{"d2_field1", "d2_field2", "d2_field3"},
						Relationship: "belongs_to_2",
						OrderFields:  []string{"d2_field1", "d2_field3"},
					},
					{
						Class:        "clubClass2",
						Fields:       []string{"cl2_field1", "cl2_field2"},
						Relationship: "member_of_2",
						OrderFields:  []string{"cl2_field1"},
					},
				},
				CompositeSections: nil,
			}},
		}, tc)
	})
}

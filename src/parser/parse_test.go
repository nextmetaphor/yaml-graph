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
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_loadSpecification(t *testing.T) {
	t.Run("MissingSpecification", func(t *testing.T) {

		d := LoadDictionary("_test", "yaml")

		assert.Equal(t, d, Dictionary{
			"MyClass": {
				"Definition1_ID": {
					Fields: definition.Fields{
						"Name":        "Definition1_Name",
						"Description": "Definition1_Description",
					},
					References: []definition.Reference{
						{
							Class:        "MyFriendClass",
							ID:           "Friend",
							Relationship: "Friend",
						},
						{
							Class:        "MyEnemyClass",
							ID:           "Enemy",
							Relationship: "Enemy",
						},
					},
				},
				"Definition2_ID": {
					Fields: definition.Fields{
						"Name":        "Definition2_Name",
						"Description": "Definition2_Description",
					},
					References: []definition.Reference{
						{
							Class:        "MyClass",
							ID:           "Definition1_ID",
							Relationship: "LINKED_CLASS",
						},
						{
							Class:        "MyFriendClass",
							ID:           "Friend",
							Relationship: "Friend",
						},
						{
							Class:        "MyEnemyClass",
							ID:           "Enemy",
							Relationship: "Enemy",
						},
					},
				},
			},
			"ChildClass": {
				"ChildClass1": {
					Fields: definition.Fields{
						"Name":        "ChildClass1",
						"Description": "ChildClassDescription1",
					},
					References: []definition.Reference{
						{
							Class:        "MyClass",
							ID:           "Definition1_ID",
							Relationship: "child_of",
						},
					},
				},
			},
			"Animal": {
				"Definition3_ID": {
					Fields: definition.Fields{
						"Name":        "Definition3_Name",
						"Description": "Definition3_Description",
					},
					References: []definition.Reference{
						{
							Class:        "Parent",
							ID:           "id1",
							Relationship: "ParentRelationship",
						},
					},
				},
				"Definition4_ID": {
					Fields: definition.Fields{
						"Name":        "Definition4_Name",
						"Description": "Definition4_Description",
					},
					References: []definition.Reference{
						{
							Class:        "AnotherClass",
							ID:           "DefinitionX_ID",
							Relationship: "Tenuous Link",
						},
						{
							Class:        "Parent",
							ID:           "id1",
							Relationship: "ParentRelationship",
						},
					},
				},
			},
			"SubChildClass": {
				"SubChildClass1": {
					Fields: definition.Fields{
						"Name":        "SubChildClass1",
						"Description": "SubChildClassDescription1",
					},
					References: []definition.Reference{
						{
							Class:        "Animal",
							ID:           "Definition3_ID",
							Relationship: "subclassed_by",
						},
					},
				},
			},
			"SubSubChildClass": {
				"SubSubChildClass1": {
					Fields: definition.Fields{
						"Name":        "SubSubChildClass1",
						"Description": "SubSubChildClassDescription1",
					},
					References: []definition.Reference{
						{
							Class:        "SubChildClass",
							ID:           "SubChildClass1",
							Relationship: "further_subclassed_by",
						},
					},
				},
			},
		})
	})
}

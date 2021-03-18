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

package cmd

import (
	"github.com/nextmetaphor/yaml-graph/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mergeDefinitionFormat(t *testing.T) {
	class1Spec := parser.ClassDefinitionFormat{
		Description: "class1 description",
		MandatoryFields: map[string]parser.ClassField{
			"name":        {Description: "name field"},
			"description": {},
		},
		OptionalFields: map[string]parser.ClassField{
			"best friend": {Description: "name of best friend"},
		},
	}
	class2Spec := parser.ClassDefinitionFormat{
		Description: "class2 description",
		MandatoryFields: map[string]parser.ClassField{
			"name field":        {Description: "class 2 name field"},
			"description field": {},
		},
		OptionalFields: map[string]parser.ClassField{
			"best friend field": {Description: "class 2 name of best friend"},
		},
	}

	t.Run("EmptyNewFormat", func(t *testing.T) {
		currentFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class1": &class1Spec,
			"class2": &class2Spec,
		}}

		newFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{}}

		err := mergeDefinitionFormat(&currentFormat, &newFormat)

		assert.Nil(t, err)
		assert.Equal(t,
			parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
				"class1": &class1Spec,
				"class2": &class2Spec,
			}}, currentFormat)
	})

	t.Run("SingleNewFormat", func(t *testing.T) {
		currentFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class1": &class1Spec,
		}}

		newFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class2": &class2Spec,
		}}

		err := mergeDefinitionFormat(&currentFormat, &newFormat)

		assert.Nil(t, err)
		assert.Equal(t,
			parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
				"class1": &class1Spec,
				"class2": &class2Spec,
			}}, currentFormat)
	})

	t.Run("MultipleNewFormat", func(t *testing.T) {
		currentFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{}}

		newFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class1": &class1Spec,
			"class2": &class2Spec,
		}}

		err := mergeDefinitionFormat(&currentFormat, &newFormat)

		assert.Nil(t, err)
		assert.Equal(t,
			parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
				"class1": &class1Spec,
				"class2": &class2Spec,
			}}, currentFormat)
	})

	t.Run("ClashingDefinitions", func(t *testing.T) {
		currentFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class1": &class1Spec,
		}}

		newFormat := parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"class1": &class1Spec,
			"class2": &class2Spec,
		}}

		err := mergeDefinitionFormat(&currentFormat, &newFormat)

		assert.NotNil(t, err)
	})

}

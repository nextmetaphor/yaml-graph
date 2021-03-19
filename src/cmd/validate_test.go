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

func Test_loadDefinitionFormatConf(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		_, err := loadDefinitionFormatConf("_test/validate/format-invalid.yml")

		assert.NotNil(t, err)
	})

	t.Run("Empty", func(t *testing.T) {
		fmt, err := loadDefinitionFormatConf("_test/validate/format-no-entries.yml")

		assert.Nil(t, err)
		assert.Equal(t, &parser.DefinitionFormat{ClassFormat: nil}, fmt)
	})

	t.Run("Single", func(t *testing.T) {
		fmt, err := loadDefinitionFormatConf("_test/validate/format-single-entry.yml")

		assert.Nil(t, err)
		assert.Equal(t, &parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"Attribute": {
				Description: "",
				MandatoryFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: ""},
				},
				OptionalFields: nil,
			},
		}}, fmt)
	})

	t.Run("Two", func(t *testing.T) {
		fmt, err := loadDefinitionFormatConf("_test/validate/format-two-entries.yml")

		assert.Nil(t, err)
		assert.Equal(t, &parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"Capability": {
				Description: "",
				MandatoryFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: ""},
				},
				OptionalFields: nil,
			},
			"Category": {
				Description: "",
				MandatoryFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: ""},
				},
				OptionalFields: nil,
			},
		}}, fmt)
	})

	t.Run("Three", func(t *testing.T) {
		fmt, err := loadDefinitionFormatConf("_test/validate/format-three-entries.yml")

		assert.Nil(t, err)
		assert.Equal(t, &parser.DefinitionFormat{ClassFormat: map[string]*parser.ClassDefinitionFormat{
			"Provider": {
				Description: "",
				MandatoryFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: ""},
				},
				OptionalFields: map[string]parser.ClassField{
					"Test1": {Description: ""},
				},
			},
			"Service": {
				Description: "",
				MandatoryFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: ""},
					"Link":        {Description: ""},
				},
				OptionalFields: map[string]parser.ClassField{
					"Test1": {Description: ""},
					"Test2": {Description: ""},
				},
			},
			"Tenancy": {
				Description:     "",
				MandatoryFields: nil,
				OptionalFields: map[string]parser.ClassField{
					"Name":        {Description: ""},
					"Description": {Description: "here is the description of the description"},
				},
			},
		}}, fmt)
	})

}

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

package definition

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_loadSpecificationFromFile(t *testing.T) {
	t.Run("MissingSpecification", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/EmptyDefinition_zzz.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("EmptySpecification", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/EmptyDefinition.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("InvalidSpecification", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/InvalidYAMLStructure.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/InvalidYAML.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("CompleteSpecification_Flow", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/CompleteDefinition_Flow.yaml")

		assert.Nil(t, err)
		testString := "MyClass"

		assert.Equal(t, &Specification{
			Class: testString,
			References: []Reference{
				{Class: "MyFriendClass", ID: "Friend", Relationship: "Friend"},
				{Class: "MyEnemyClass", ID: "Enemy", Relationship: "Enemy"},
			},
			Definitions: map[string]Definition{
				"Definition1_ID": {
					Fields:     map[string]interface{}{"Name": "Definition1_Name", "Description": "Definition1_Description"},
					References: nil,
				},
				"Definition2_ID": {
					Fields: map[string]interface{}{"Name": "Definition2_Name", "Description": "Definition2_Description", "ImgSrc": "data:image;base64,c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ=="},
					References: []Reference{
						{
							Class:        "MyClass",
							ID:           "Definition1_ID",
							Relationship: "LINKED_CLASS",
						},
					},
					FileFields: FileFields{"ImgSrc": FileDefinition{Path: "simple-file.txt", Prefix: "data:image;base64,"}},
				},
			},
		}, spec)
	})

	t.Run("CompleteSpecification_NonFlow", func(t *testing.T) {
		spec, err := LoadSpecificationFromFile("./_test/Structured/CompleteDefinition_NonFlow.yaml")

		assert.Nil(t, err)
		testString := "MyClass"

		assert.Equal(t, &Specification{
			Class: testString,
			References: []Reference{
				{Class: "MyFriendClass", ID: "Friend", Relationship: "Friend"},
				{Class: "MyEnemyClass", ID: "Enemy", Relationship: "Enemy"},
			},
			Definitions: map[string]Definition{
				"Definition1_ID": {
					Fields:     map[string]interface{}{"Name": "Definition1_Name", "Description": "Definition1_Description"},
					References: nil,
					SubDefinitions: map[string]Specification{"child_of": {
						Class:      "ChildClass",
						References: nil,
						Definitions: map[string]Definition{"ChildClass1": {
							Fields:         map[string]interface{}{"Name": "ChildClass1", "Description": "ChildClassDescription1", "ImgSrc": "data:image;base64,c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ=="},
							FileFields:     FileFields{"ImgSrc": FileDefinition{Path: "simple-file.txt", Prefix: "data:image;base64,"}},
							References:     nil,
							SubDefinitions: nil,
						}},
					}},
				},
				"Definition2_ID": {
					Fields: map[string]interface{}{"Name": "Definition2_Name", "Description": "Definition2_Description"},
					References: []Reference{
						{
							Class:        "MyClass",
							ID:           "Definition1_ID",
							Relationship: "LINKED_CLASS",
						},
					},
				},
			},
		}, spec)
	})
}

func Test_processFiles(t *testing.T) {
	filesProcessed := map[string]string{}

	processFileFunc := func(filePath string, fileInfo os.FileInfo) (err error) {
		filesProcessed[fileInfo.Name()] = filePath

		return nil
	}

	t.Run("MissingRootDir", func(t *testing.T) {
		err := ProcessFiles("./_test/NotThere", "yaml", processFileFunc)

		assert.NotNil(t, err)
	})

	t.Run("ValidRootDir", func(t *testing.T) {
		err := ProcessFiles("./_test/ProcessFiles", "yaml", processFileFunc)

		assert.Nil(t, err)
		assert.Equal(t, map[string]string{
			"1.2.1.yaml": "_test/ProcessFiles/1/1.2/1.2.1/1.2.1.yaml",
			"1.2.yaml":   "_test/ProcessFiles/1/1.2/1.2.yaml",
			"2.yaml":     "_test/ProcessFiles/2/2.yaml",
		}, filesProcessed)
	})
}

func Test_getFileFieldAsBase64(t *testing.T) {
	t.Run("MissingFile", func(t *testing.T) {
		str, err := getFileFieldAsBase64("./_test/", FileDefinition{Path: "NotThere"})
		assert.Nil(t, str)
		assert.NotNil(t, err)
	})

	t.Run("ValidFile", func(t *testing.T) {
		str, err := getFileFieldAsBase64("./_test/Base64/", FileDefinition{Path: "simple-file.txt"})
		assert.Nil(t, err)

		simpleFileBase64Str := "c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ=="
		assert.Equal(t, str, &simpleFileBase64Str)
	})

	t.Run("ValidFileWithPrefix", func(t *testing.T) {
		str, err := getFileFieldAsBase64("./_test/Base64/", FileDefinition{Path: "simple-file.txt", Prefix: "data:image;base64,"})
		assert.Nil(t, err)

		simpleFileBase64Str := "data:image;base64,c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ=="
		assert.Equal(t, str, &simpleFileBase64Str)
	})
}

func Test_getFileFields(t *testing.T) {
	dfn := Definition{
		Fields:     Fields{"Name": "Definition1_Name", "Description": "Definition1_Description"},
		FileFields: FileFields{"ImgSrc": {Path: "simple-file.txt"}},
		References: nil,
		SubDefinitions: map[string]Specification{"child_of": {
			Class:      "ChildClass",
			References: nil,
			Definitions: map[string]Definition{"ChildClass1": {
				Fields:     map[string]interface{}{"Name": "ChildClass1", "Description": "ChildClassDescription1"},
				FileFields: FileFields{"ImgSrc": FileDefinition{Path: "simple-file.txt", Prefix: "data:image;base64,"}},
			}},
		}},
	}

	t.Run("ValidFile", func(t *testing.T) {
		getFileFields("./_test/Base64/", &dfn)
		assert.Equal(t, Fields{"Name": "Definition1_Name", "ImgSrc": "c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ==", "Description": "Definition1_Description"}, dfn.Fields)
		assert.Equal(t, Fields{"Name": "ChildClass1", "ImgSrc": "data:image;base64,c2ltcGxlIGZpbGUgdG8gYjY0IGVuY29kZQ==", "Description": "ChildClassDescription1"}, dfn.SubDefinitions["child_of"].Definitions["ChildClass1"].Fields)
	})
}

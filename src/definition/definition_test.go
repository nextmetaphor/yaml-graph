package definition

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_loadSpecificationFromFile(t *testing.T) {
	t.Run("MissingSpecification", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/EmptyDefinition_zzz.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("EmptySpecification", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/EmptyDefinition.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("InvalidSpecification", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/InvalidYAMLStructure.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/InvalidYAML.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, spec)
	})

	t.Run("CompleteSpecification_Flow", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/CompleteDefinition_NonFlow.yaml")

		assert.Nil(t, err)
		testString := "MyClass"

		assert.Equal(t, &Specification{
			Class: &testString,
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

	t.Run("CompleteSpecification_NonFlow", func(t *testing.T) {
		spec, err := loadSpecificationFromFile("./_test/Structured/CompleteDefinition_NonFlow.yaml")

		assert.Nil(t, err)
		testString := "MyClass"

		assert.Equal(t, &Specification{
			Class: &testString,
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

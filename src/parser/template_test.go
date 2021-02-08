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
	"bufio"
	"bytes"
	"fmt"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_getOrderClause(t *testing.T) {
	t.Run("SingleOrderClause", func(t *testing.T) {
		clause := getOrderClause(
			TemplateSection{
				SectionClass: ClassFieldSelector{
					Class:        "class1",
					Fields:       nil,
					Relationship: "",
					OrderFields:  []string{"field1"},
				},
				AggregateClasses:  nil,
				CompositeSections: nil,
			})

		assert.Equal(t, "class1.field1", clause)
	})

	t.Run("MultipleOrderClause", func(t *testing.T) {
		clause := getOrderClause(TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class2",
				Fields:       nil,
				Relationship: "",
				OrderFields:  []string{"field1", "field2"},
			},
			AggregateClasses:  nil,
			CompositeSections: nil,
		})

		assert.Equal(t, "class2.field1,class2.field2", clause)
	})
	t.Run("NilOrderClause", func(t *testing.T) {
		clause := getOrderClause(TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class3",
				Fields:       nil,
				Relationship: "",
				OrderFields:  nil,
			},
			AggregateClasses:  nil,
			CompositeSections: nil,
		})

		assert.Equal(t, "", clause)
	})
	t.Run("EmptyOrderClause", func(t *testing.T) {
		clause := getOrderClause(TemplateSection{
			SectionClass: ClassFieldSelector{
				Class:        "class4",
				Fields:       nil,
				Relationship: "",
				OrderFields:  []string{},
			},
			AggregateClasses:  nil,
			CompositeSections: nil,
		})

		assert.Equal(t, "", clause)
	})
}

func Test_getCypherForSelector(t *testing.T) {
	t.Run("ParentClass", func(t *testing.T) {
		cypher := getCypherForSection("parent", "parentID",
			TemplateSection{
				SectionClass: ClassFieldSelector{
					Class:        "child",
					Fields:       nil,
					Relationship: "inherits_from",
					OrderFields:  []string{"field1", "field2"},
				},
				AggregateClasses:  nil,
				CompositeSections: nil,
			})

		assert.Equal(t, "match (child:child)-[:inherits_from]-(parent:parent {ID:\"parentID\"}) return child "+
			"order by child.field1,child.field2", cypher)
	})
	t.Run("NoParentClass", func(t *testing.T) {
		cypher := getCypherForSection("", "",
			TemplateSection{
				SectionClass: ClassFieldSelector{
					Class:        "myclass",
					Fields:       nil,
					Relationship: "",
					OrderFields:  []string{"field1", "field2"},
				},
				AggregateClasses:  nil,
				CompositeSections: nil,
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

func Test_recurseTemplateSection(t *testing.T) {
	t.Run("MinimalConfiguration", func(t *testing.T) {

		// first load the template configuration
		tc, err := loadTemplateConf("_test/template/TemplateSection_minimal_valid.yaml")
		assert.Nil(t, err)

		// then connect to the graph database
		driver, session, err := graph.Init("bolt://localhost:7687", "username", "password")
		assert.Nil(t, err)
		assert.NotNil(t, tc)

		defer driver.Close()
		defer session.Close()

		definitions, err := recurseTemplateSection(session, *tc, nil, nil)
		assert.Nil(t, err)
		assert.Equal(t, []SectionDefinition{
			{
				Class: "Provider",
				ID:    "aws",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "AWS",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Amazon Web Services",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Provider",
				ID:    "alibaba",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Alibaba",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Alibaba Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Provider",
				ID:    "azure",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Azure",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Microsoft Azure",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Provider",
				ID:    "gcp",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "GCP",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Google Cloud Platform",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Provider",
				ID:    "ibm",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "IBM",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "IBM",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Provider",
				ID:    "oracle",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Oracle Cloud",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Oracle Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			}},
			definitions)
	})
	t.Run("MinimalAggregateConfiguration", func(t *testing.T) {

		// first load the template configuration
		tc, err := loadTemplateConf("_test/template/TemplateSection_minimal_aggregate_valid.yaml")
		assert.Nil(t, err)

		// then connect to the graph database
		driver, session, err := graph.Init("bolt://localhost:7687", "username", "password")
		assert.Nil(t, err)
		assert.NotNil(t, tc)

		defer driver.Close()
		defer session.Close()

		definitions, err := recurseTemplateSection(session, *tc, nil, nil)
		assert.Nil(t, err)

		fmt.Println(definitions)
		assert.Equal(t, []SectionDefinition{
			{
				Class: "Service",
				ID:    "app-service",
				Fields: map[string]string{

					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "App Service",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "application-gateway",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Application Gateway",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "archive-storage",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Archive Storage",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "azure-function",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Azure Function",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "azure-kubernetes-service",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Azure Kubernetes Service (AKS)",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "blob-storage",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Blob Storage",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "container-instance",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Container Instance",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "content-delivery-network",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Content Delivery Network",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "disk-storage",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Disk Storage",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "file-storage",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "File Storage",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "load-balancer",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Load Balancer",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "vpn-gateway",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "VPN Gateway",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "virtual-machine",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Machine",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "virtual-machine-scale-set",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Machine Scale Set",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
			{
				Class: "Service",
				ID:    "virtual-network",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Network",
					fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{},
			},
		},

			definitions)
	})

	t.Run("MinimalCompositeConfiguration", func(t *testing.T) {

		// first load the template configuration
		tc, err := loadTemplateConf("_test/template/TemplateSection_minimal_composite_valid.yaml")
		assert.Nil(t, err)

		// then connect to the graph database
		driver, session, err := graph.Init("bolt://localhost:7687", "username", "password")
		assert.Nil(t, err)
		assert.NotNil(t, tc)

		defer driver.Close()
		defer session.Close()

		definitions, err := recurseTemplateSection(session, *tc, nil, nil)
		assert.Nil(t, err)

		fmt.Println(definitions)
		assert.Equal(t, []SectionDefinition{
			{
				Class: "Provider",
				ID:    "aws",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "AWS",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Amazon Web Services",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "alibaba",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Alibaba",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Alibaba Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "azure",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Azure",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Microsoft Azure",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": {
						{
							Class: "Service",
							ID:    "app-service",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "App Service",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "application-gateway",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Application Gateway",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "archive-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Archive Storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "azure-function",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Azure Function",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "azure-kubernetes-service",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Azure Kubernetes Service (AKS)",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "blob-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Blob Storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "container-instance",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Container Instance",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "content-delivery-network",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Content Delivery Network",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "disk-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Disk Storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "file-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "File Storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "load-balancer",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Load Balancer",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "vpn-gateway",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "VPN Gateway",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-machine",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Virtual Machine",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-machine-scale-set",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Virtual Machine Scale Set",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-network",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"): "Virtual Network",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
					},
				},
			},
			{
				Class: "Provider",
				ID:    "gcp",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "GCP",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Google Cloud Platform",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "ibm",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "IBM",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "IBM",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "oracle",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Oracle Cloud",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Oracle Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
		},
			definitions)
	})
	t.Run("MinimalCompositeAggregateConfiguration", func(t *testing.T) {

		// first load the template configuration
		tc, err := loadTemplateConf("_test/template/TemplateSection_minimal_composite_aggregate_valid.yaml")
		assert.Nil(t, err)

		// then connect to the graph database
		driver, session, err := graph.Init("bolt://localhost:7687", "username", "password")
		assert.Nil(t, err)
		assert.NotNil(t, tc)

		defer driver.Close()
		defer session.Close()

		definitions, err := recurseTemplateSection(session, *tc, nil, nil)
		assert.Nil(t, err)

		fmt.Println(definitions)
		assert.Equal(t, []SectionDefinition{
			{
				Class: "Provider",
				ID:    "aws",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "AWS",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Amazon Web Services",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "alibaba",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Alibaba",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Alibaba Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "azure",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Azure",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Microsoft Azure",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": {
						{
							Class: "Service",
							ID:    "app-service",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "App Service",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "application-gateway",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Application Gateway",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "archive-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Archive Storage",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "azure-function",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Azure Function",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "azure-kubernetes-service",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Azure Kubernetes Service (AKS)",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "blob-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Blob Storage",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "container-instance",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Container Instance",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "content-delivery-network",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Content Delivery Network",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "disk-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Disk Storage",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "file-storage",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "File Storage",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "storage",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "load-balancer",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Load Balancer",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "vpn-gateway",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "VPN Gateway",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-machine",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Machine",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-machine-scale-set",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Machine Scale Set",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "compute",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
						{
							Class: "Service",
							ID:    "virtual-network",
							Fields: map[string]string{
								fmt.Sprintf(classFieldIdentifier, "Service", "Name"):  "Virtual Network",
								fmt.Sprintf(classFieldIdentifier, "Category", "Name"): "networking",
							},
							CompositeSectionDefinitions: map[string][]SectionDefinition{},
						},
					},
				},
			},
			{
				Class: "Provider",
				ID:    "gcp",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "GCP",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Google Cloud Platform",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "ibm",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "IBM",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "IBM",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
			{
				Class: "Provider",
				ID:    "oracle",
				Fields: map[string]string{
					fmt.Sprintf(classFieldIdentifier, "Provider", "Name"):        "Oracle Cloud",
					fmt.Sprintf(classFieldIdentifier, "Provider", "Description"): "Oracle Cloud",
				},
				CompositeSectionDefinitions: map[string][]SectionDefinition{
					"PROVIDED_BY": nil,
				},
			},
		},
			definitions)
	})
}

func Test_parseTemplate(t *testing.T) {
	t.Run("ParseCompositeAggregateTemplate", func(t *testing.T) {
		var writer bytes.Buffer
		bufferWriter := bufio.NewWriter(&writer)

		err := parseTemplate("bolt://localhost:7687", "username", "password", "_test/template/TemplateSection_minimal_composite_aggregate_valid.yaml", "_test/template/output-template.gotmpl", bufferWriter)
		assert.Nil(t, err)

		expectedBytes, err := ioutil.ReadFile("_test/template/output-template.result")
		assert.Nil(t, err)

		bufferWriter.Flush()
		assert.Equal(t, expectedBytes, writer.Bytes())})
}

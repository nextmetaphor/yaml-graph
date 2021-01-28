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
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	rootJSONCypher  = "match (n:%s) return n"
	childJSONCypher = "match (n:%s)-[:%s]-(p:%s {ID:\"%s\"}) return n"

	logErrorExecutingJSONCypher                       = "error executing cypher"
	logErrorCouldNotOpenJSONConfiguration             = "could not open JSON configuration [%s]"
	logErrorCouldNotUnmarshalJSONConfiguration        = "could not unmarshal JSON configuration [%s]"
	logDebugSuccessfullyUnmarshalledJSONConfiguration = "successfully unmarshalled JSON configuration [%s]"
)

type (
	// JSONLevel TODO
	JSONLevel struct {
		// Class indicates the class of definition to use for the node
		Class string `yaml:"Class"`
		// NameField indicates which field to use as the "name" element
		NameField string `yaml:"NameField"`
		// Colour field is used to add a "colour" element
		Colour string `yaml:"Colour"`
		// Size field is used to add a "colour" element
		Size int `yaml:"Size"`

		// DetailFields indicates which fields which will be extracted
		DetailFields []string `yaml:"DetailFields"`

		ParentRelationship string      `yaml:"ParentRelationship"`
		ChildLevel         []JSONLevel `yaml:"ChildLevel"`
	}
)

type (
	jsonNode struct {
		Class        string            `json:"class"`
		Name         string            `json:"name"`
		Colour       string            `json:"colour"`
		Size         int               `json:"size"`
		DetailFields map[string]string `json:"detail-fields"`
		Children     []*jsonNode       `json:"children"`
	}
)

var (
	jsonCmd = &cobra.Command{
		Use:   commandJSONUse,
		Short: commandJSONUseShort,
		Run:   jsonFunc,
	}

	jsonDefinition string
)

func init() {
	rootCmd.AddCommand(jsonCmd)

	jsonCmd.PersistentFlags().StringVarP(&jsonDefinition, flagJSONDefinitionName, flagJSONDefinitionShorthand,
		"", flagJSONDefinitionUsage)
	jsonCmd.MarkPersistentFlagRequired(flagJSONDefinitionName)
}

func loadJSONConf(cfgPath string) (ms *JSONLevel, err error) {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenJSONConfiguration, cfgPath)
		return nil, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &ms)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotUnmarshalJSONConfiguration, cfgPath)
		return nil, err
	}

	log.Debug().Msgf(logDebugSuccessfullyUnmarshalledJSONConfiguration, cfgPath)
	return ms, nil
}

func jsonFunc(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	// first load the jsonFunc configuration
	jsonLevel, err := loadJSONConf(jsonDefinition)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeJSONCmdFailed)
	}

	// then connect to the graph database
	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeJSONCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	// now recurse through the sections
	rootNode := new(jsonNode)
	rootNode.Colour = "#4dc2ca"

	j, _ := recurseLevel(session, *jsonLevel, nil, nil)
	rootNode.Children = append(rootNode.Children, j...)

	jb, e := json.Marshal(rootNode)
	if e == nil {
		fmt.Print(string(jb))
	} else {
		fmt.Print(e)
	}
}

func recurseLevel(session neo4j.Session, level JSONLevel, parentClass, parentID *string) ([]*jsonNode, error) {
	var res neo4j.Result
	var err error
	var nodes []*jsonNode

	if (parentClass == nil) || (parentID == nil) {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(rootJSONCypher, level.Class), nil)
	} else {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(childJSONCypher, level.Class, level.ParentRelationship, *parentClass,
			*parentID), nil)
	}

	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingJSONCypher)
		return nil, err
	}

	for res.Next() {
		record := res.Record()
		for _, kv := range record.Values() {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				jNode := new(jsonNode)
				nodes = append(nodes, jNode)
				jNode.Class = level.Class
				if node.Props()[level.NameField] != nil {
					jNode.Name = node.Props()[level.NameField].(string)
				}
				jNode.Colour = level.Colour
				jNode.Size = level.Size
				jNode.DetailFields = map[string]string{}
				jNode.Children = []*jsonNode{}

				for _, detailField := range level.DetailFields {
					if node.Props()[detailField] != nil {
						jNode.DetailFields[detailField] = node.Props()[detailField].(string)
					}
				}

				// TODO recursion, really?
				for _, childLevel := range level.ChildLevel {
					nodeID := node.Props()["ID"].(string)
					childNodes, err := recurseLevel(session, childLevel, &(level.Class), &nodeID)
					if err != nil {
						log.Error().Err(err).Msgf(logErrorExecutingJSONCypher)
						return nil, err
					}
					jNode.Children = append(jNode.Children, childNodes...)
				}
			}
		}
	}

	return nodes, err
}

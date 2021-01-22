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
	classSection = "\"class\": \"%s\","
	nameSection  = "\"name\": \"%s\""

	childrenLevelPrefix = ",\"children\": ["
	childrenLevelSuffix = "]"

	rootJSONCypher  = "match (n:%s) return n"
	childJSONCypher = "match (n:%s)-[:%s]-(p:%s {ID:\"%s\"}) return n"

	//markdownSection     = "%s%s%s"       //prefix section suffix
	//markdownDetailField = "%s%s%s%s%s%s" //field key (prefix value suffix) + field value (prefix value suffix)

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

		NameField string `yaml:"NameField"`

		// DetailFields indicates which fields which will be extracted
		DetailFields []string `yaml:"DetailFields"`

		ParentRelationship string      `yaml:"ParentRelationship"`
		ChildLevel         []JSONLevel `yaml:"ChildLevel"`
	}
)

var (
	jsonCmd = &cobra.Command{
		Use:   commandJSONUse,
		Short: commandJSONUseShort,
		Run:   json,
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

func json(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	// first load the json configuration
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
	fmt.Print("{" + fmt.Sprintf(nameSection, "root") + childrenLevelPrefix)
	recurseLevel(session, *jsonLevel, nil, nil)
	fmt.Print(childrenLevelSuffix + "}")
}

func recurseLevel(session neo4j.Session, level JSONLevel, parentClass, parentID *string) error {
	var res neo4j.Result
	var err error

	if (parentClass == nil) || (parentID == nil) {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(rootJSONCypher, level.Class), nil)
	} else {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(childJSONCypher, level.Class, level.ParentRelationship, *parentClass,
			*parentID), nil)
	}

	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingJSONCypher)
		return err
	}

	firstResult := true
	for res.Next() {
		record := res.Record()
		for _, kv := range record.Values() {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				if firstResult {
					fmt.Print("{")
				} else {
					fmt.Print(",{")
				}
				firstResult = false

				fmt.Print(fmt.Sprintf(classSection, level.Class))
				fmt.Print(fmt.Sprintf(nameSection, node.Props()[level.NameField]))

				// TODO recursion, really?
				for _, childLevel := range level.ChildLevel {
					nodeID := node.Props()["ID"].(string)
					fmt.Print(childrenLevelPrefix)
					recurseLevel(session, childLevel, &(level.Class), &nodeID)
					fmt.Print(childrenLevelSuffix)
				}

				fmt.Print("}")
			}
		}
	}

	return err
}

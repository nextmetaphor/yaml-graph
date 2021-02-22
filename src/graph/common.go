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
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/rs/zerolog/log"
)

const (
	deleteAllCypher   = `MATCH (n) DETACH DELETE(n);`
	mergeCypherPrefix = `MERGE (n:%s {ID:$ID})`
	mergeCypherField  = `n.%s=$%s`

	edgeCypher = `
		MATCH (n1:%s {ID:"%s"})
		MATCH (n2:%s {ID:"%s"})
		MERGE (n1)-[:%s]->(n2);`

	reverseEdgeCypher = `
		MATCH (n1:%s {ID:"%s"})
		MATCH (n2:%s {ID:"%s"})
		MERGE (n1)-[:%s]->(n2);`

	logErrorCannotConnectToGraphDatabase = "cannot connect to graph database"
	logErrorCannotCreateGraphSession     = "cannot create graph session"
	logDebugCypherDetails                = "about to execute cypher [%s] with [%s]"
	logErrorCannotRunCypher              = "error when running cypher"
	logWarnCypherReturnedError           = "cypher returned error"
)

// DeleteAll TODO
func DeleteAll(session neo4j.Session) (neo4j.Result, error) {
	return ExecuteCypher(session, deleteAllCypher, nil)
}

// Init TODO
func Init(dbURL, username, password string) (driver neo4j.Driver, session neo4j.Session, err error) {
	configForNeo4j40 := func(conf *neo4j.Config) {
		conf.Encrypted = false
	}

	driver, err = neo4j.NewDriver(dbURL, neo4j.BasicAuth(username, password, ""), configForNeo4j40)
	if err != nil {
		log.Error().Err(err).Msg(logErrorCannotConnectToGraphDatabase)
		return
	}

	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Error().Err(err).Msg(logErrorCannotCreateGraphSession)
		return
	}

	return
}

// ExecuteCypher TODO
func ExecuteCypher(session neo4j.Session, cypher string, param map[string]interface{}) (res neo4j.Result, err error) {
	log.Debug().Msgf(logDebugCypherDetails, cypher, param)

	if res, err = session.Run(cypher, param); err != nil {
		log.Error().Err(err).Msg(logErrorCannotRunCypher)

		return nil, err
	}

	return
}

func getDefinitionCypherString(class string, fields definition.Fields) (cypher string) {
	cypher = fmt.Sprintf(mergeCypherPrefix, class)

	firstValue := true
	for fieldName := range fields {
		if !firstValue {
			cypher = cypher + ","
		} else {
			cypher = cypher + " SET "
			firstValue = false
		}

		cypher = cypher + fmt.Sprintf(mergeCypherField, fieldName, fieldName)
	}

	return cypher
}

func getEdgeCypherString(class, ID string, refs definition.Reference) string {
	if refs.FromNotTo {
		return fmt.Sprintf(edgeCypher, class, ID, refs.Class, refs.ID, refs.Relationship)
	}
	return fmt.Sprintf(reverseEdgeCypher, class, ID, refs.Class, refs.ID, refs.Relationship)
}

// CreateSpecification TODO
func CreateSpecification(session neo4j.Session, spec definition.Specification) {
	class := spec.Class
	for definitionID := range spec.Definitions {
		definitionCypher := getDefinitionCypherString(class, spec.Definitions[definitionID].Fields)

		if spec.Definitions[definitionID].Fields == nil {
			ExecuteCypher(session, definitionCypher, map[string]interface{}{"ID": definitionID})
		} else {
			spec.Definitions[definitionID].Fields["ID"] = definitionID
			ExecuteCypher(session, definitionCypher, spec.Definitions[definitionID].Fields)
		}

		// now recurse through the subdefinitions
		// TODO - do we really want to use recursion for this?
		if spec.Definitions[definitionID].SubDefinitions != nil {
			for subdefinitionID := range spec.Definitions[definitionID].SubDefinitions {
				CreateSpecification(session, spec.Definitions[definitionID].SubDefinitions[subdefinitionID])
			}
		}
	}
}

// CreateSpecificationEdge TODO
func CreateSpecificationEdge(session neo4j.Session, spec definition.Specification, parentReference *definition.Reference) {
	class := spec.Class

	for definitionID := range spec.Definitions {
		for _, ref := range spec.References {
			// create the Specification-scoped references
			referenceCypher := getEdgeCypherString(class, definitionID, ref)
			ExecuteCypher(session, referenceCypher, nil)
		}

		for _, ref := range spec.Definitions[definitionID].References {
			// create any Definition-scoped references
			referenceCypher := getEdgeCypherString(class, definitionID, ref)
			ExecuteCypher(session, referenceCypher, nil)
		}

		// if we have a parentReference then create an appropriate relationship
		if parentReference != nil {
			referenceCypher := getEdgeCypherString(class, definitionID, *parentReference)
			ExecuteCypher(session, referenceCypher, nil)
		}

		// recurse through the subdefinitions
		// TODO - do we really want to use recursion for this?
		if spec.Definitions[definitionID].SubDefinitions != nil {
			for subdefRelationship := range spec.Definitions[definitionID].SubDefinitions {
				CreateSpecificationEdge(session, spec.Definitions[definitionID].SubDefinitions[subdefRelationship],
					&definition.Reference{
						Class:        class,
						ID:           definitionID,
						Relationship: subdefRelationship,
					})
			}
		}
	}
}

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

	logErrorCannotConnectToGraphDatabase = "cannot connect to graph database"
	logErrorCannotCreateGraphSession     = "cannot create graph session"
	logDebugCypherDetails                = "about to execute cypher [%s] with [%s]"
	logErrorCannotRunCypher              = "error when running cypher"
	logWarnCypherReturnedError           = "cypher returned error"
)

// DeleteAll TODO
func DeleteAll(session neo4j.Session) error {
	return executeCypher(session, deleteAllCypher, nil)
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

func executeCypher(session neo4j.Session, cypher string, param map[string]interface{}) error {
	var res neo4j.Result
	var err error

	log.Debug().Msgf(logDebugCypherDetails, cypher, param["ID"])

	if res, err = session.Run(cypher, param); err != nil {
		log.Error().Err(err).Msg(logErrorCannotRunCypher)

		return err
	}

	if res.Err() != nil {
		log.Warn().Err(res.Err()).Msg(logWarnCypherReturnedError)
	}
	return res.Err()
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
	return fmt.Sprintf(edgeCypher, class, ID, refs.Class, refs.ID, refs.Relationship)
}

// CreateSpecification TODO
func CreateSpecification(session neo4j.Session, spec definition.Specification) {
	class := spec.Class
	for definitionID := range spec.Definitions {
		definitionCypher := getDefinitionCypherString(class, spec.Definitions[definitionID].Fields)
		spec.Definitions[definitionID].Fields["ID"] = definitionID
		executeCypher(session, definitionCypher, spec.Definitions[definitionID].Fields)
	}
}

// CreateSpecificationEdge TODO
func CreateSpecificationEdge(session neo4j.Session, spec definition.Specification) {
	class := spec.Class
	for _, ref := range spec.References {
		for definitionID := range spec.Definitions {
			referenceCypher := getEdgeCypherString(class, definitionID, ref)
			executeCypher(session, referenceCypher, nil)
		}
	}
}

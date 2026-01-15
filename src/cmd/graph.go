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

	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	nodeHeaderString = "{\"nodes\": ["
	linkHeaderString = "],\"links\": ["
	linkFooterString = "]}"
)

var (
	graphCmd = &cobra.Command{
		Use:   commandGraphUse,
		Short: commandGraphUseShort,
		Run:   graphFunc,
	}
)

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.Flags().StringSliceVarP(&sourceDir, flagSourceName, flagSourceShorthand, []string{flagSourceDefault},
		flagSourceUsage)
	graphCmd.Flags().StringSliceVarP(&graphFields, flagGraphFieldsName, flagGraphFieldsShorthand, []string{flagGraphFieldsDefault},
		flagGraphFieldsUsage)

	// still under development - hide
	graphCmd.Hidden = true
}

func graphFunc(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeLoadCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	fmt.Print(nodeHeaderString)
	firstElement := true

	res, err := graph.ExecuteCypher(session, "MATCH (n) RETURN n, labels(n)[0] as class", nil)
	if err == nil {
		for res.Next() {
			record := res.Record()
			n, _ := record.Get("n")
			node := n.(neo4j.Node)
			c, _ := record.Get("class")
			class := c.(string)
			id := node.Props["ID"].(string)

			var name interface{}
			for _, fieldName := range graphFields {
				if node.Props[fieldName] != nil {
					name = node.Props[fieldName]
					break
				}
			}
			if name == nil {
				name = id
			}
			if firstElement {
				firstElement = false
			} else {
				fmt.Print(",")
			}

			nodeData := map[string]string{
				"id":          fmt.Sprintf("%s-%s", class, id),
				"class":       class,
				"description": fmt.Sprintf("%v", name),
			}
			jb, _ := json.Marshal(nodeData)
			fmt.Print(string(jb))
		}
	}

	fmt.Print(linkHeaderString)
	firstElement = true

	res, err = graph.ExecuteCypher(session, "MATCH (n1)-[r]->(n2) RETURN n1.ID as sourceID, labels(n1)[0] as sourceClass, n2.ID as targetID, labels(n2)[0] as targetClass, type(r) as relationship", nil)
	if err == nil {
		for res.Next() {
			record := res.Record()
			sourceID, _ := record.Get("sourceID")
			sourceClass, _ := record.Get("sourceClass")
			targetID, _ := record.Get("targetID")
			targetClass, _ := record.Get("targetClass")
			relationship, _ := record.Get("relationship")

			if firstElement {
				firstElement = false
			} else {
				fmt.Print(",")
			}

			linkData := map[string]string{
				"source":       fmt.Sprintf("%s-%v", sourceClass, sourceID),
				"target":       fmt.Sprintf("%s-%v", targetClass, targetID),
				"relationship": fmt.Sprintf("%v", relationship),
			}
			jb, _ := json.Marshal(linkData)
			fmt.Print(string(jb))
		}
	}

	fmt.Print(linkFooterString)
}

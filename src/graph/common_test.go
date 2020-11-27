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
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getDefinitionCypherString(t *testing.T) {
	t.Run("ZeroFields", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{})
		assert.Equal(t, fmt.Sprintf(mergeCypherPrefix, "class1"), cypher)
	})

	t.Run("OneField", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{"name1": "value1"})
		assert.Equal(t, fmt.Sprintf(mergeCypherPrefix, "class1")+" SET "+fmt.Sprintf(mergeCypherField,
			"name1", "name1"), cypher)
	})

	t.Run("TwoFields", func(t *testing.T) {
		cypher := getDefinitionCypherString("class1", definition.Fields{
			"name1": "value1", "name2": "value2"})
		assert.Equal(t, fmt.Sprintf(mergeCypherPrefix, "class1")+" SET "+fmt.Sprintf(mergeCypherField,
			"name1", "name1")+","+fmt.Sprintf(mergeCypherField,
			"name2", "name2"), cypher)
	})
}

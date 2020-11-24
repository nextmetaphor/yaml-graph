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

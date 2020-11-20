package definition

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	logWarnCannotLoadYAMLFile           = "cannot load YAML file [%s]"
	logWarnCannotParseYAMLFile          = "cannot parse YAML file [%s]"
	logWarnNoDefinitionsFoundInYAMLFile = "no definitions found in YAML file [%s]"
)

type (
	// Reference TODO
	Reference struct {
		// Class TODO
		Class string `yaml:"Class"`

		// Relationship TODO
		Relationship string `yaml:"Relationship"`
	}

	// Definition TODO
	Definition map[string]interface{}

	// Specification TODO
	Specification struct {
		// Class allows the class for the definitions within the document to be specified.
		// If this is not specified, the subdirectory immediately below the definition root directory
		// is used as the definition class
		Class *string `yaml:"Class,omitempty"`

		// References allows relationships to other classes to be specified. If this is not specified,
		// the parent directories are used to specify these references
		References []Reference `yaml:"References,omitempty"`

		// Definitions TODO
		Definitions []Definition `yaml:"Definitions,omitempty"`
	}
)

func loadSpecificationFromFile(filename string) (*Specification, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf(logWarnCannotLoadYAMLFile, filename))
		return nil, err
	}

	spec := &Specification{}
	err = yaml.Unmarshal(yamlFile, spec)
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf(logWarnCannotParseYAMLFile, filename))

		return nil, err
	}

	// if no definitions are found, return an error and an nil Specification
	if len(spec.Definitions) == 0 {
		log.Warn().Err(err).Msg(fmt.Sprintf(logWarnNoDefinitionsFoundInYAMLFile, filename))
		return nil, errors.New(fmt.Sprintf(logWarnNoDefinitionsFoundInYAMLFile, filename))
	}

	return spec, nil
}

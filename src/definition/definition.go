package definition

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	definitionSuffix                    = ".yaml"
	logWarnCannotLoadYAMLFile           = "cannot load YAML file [%s]"
	logWarnCannotParseYAMLFile          = "cannot parse YAML file [%s]"
	logWarnNoDefinitionsFoundInYAMLFile = "no definitions found in YAML file [%s]"
	logErrorCannotReadRootDirectory     = "cannot read root directory [%]"
	logDebugFoundDefinitionSubdirectory = "found definition subdirectory [%s]"
	logErrorCannotProcessFiles          = "cannot process files in root directory [%s]"
	logWarnCannotProcessFile            = "cannot process files in directory [%s]"
	logDebugProcessingFile              = "processing file [%s] in directory [%s]"
	logDebugIgnoringFile                = "ignoring file [%s] in directory [%s]"
)

type (
	// Reference TODO
	Reference struct {
		// Class TODO
		Class string `yaml:"Class"`

		// ID TODO
		ID string `yaml:"ID"`

		// Relationship TODO
		Relationship string `yaml:"Relationship"`
	}

	// Definition TODO
	Definition struct {
		Fields     map[string]interface{} `yaml:"Fields"`
		References []Reference            `yaml:"References"`
	}

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
		Definitions map[string]Definition `yaml:"Definitions,omitempty"`
	}

	processFileFuncType = func(filePath string, fileInfo os.FileInfo) (err error)
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
		return nil, fmt.Errorf(logWarnNoDefinitionsFoundInYAMLFile, filename)
	}

	// TODO debug
	return spec, nil
}

func processFiles(rootDir string, processFileFunc processFileFuncType) error {
	err := filepath.Walk(rootDir,
		func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Warn().Err(err).Msgf(logWarnCannotProcessFile, filePath)
				return err
			}
			if !fileInfo.IsDir() {
				if strings.HasSuffix(fileInfo.Name(), definitionSuffix) {
					log.Debug().Msg(fmt.Sprintf(logDebugProcessingFile, fileInfo.Name(), filePath))
					return processFileFunc(filePath, fileInfo)

				} else {
					log.Debug().Msg(fmt.Sprintf(logDebugIgnoringFile, fileInfo.Name(), filePath))
				}
			}
			return nil
		})

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf(logErrorCannotProcessFiles, rootDir))
		return err
	}

	return nil
}

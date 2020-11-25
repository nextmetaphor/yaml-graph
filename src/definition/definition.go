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
	definitionFormat                     = ".%s"
	logDebugCannotLoadYAMLFile           = "cannot load YAML file [%s]"
	logDebugCannotParseYAMLFile          = "cannot parse YAML file [%s]"
	logDebugNoDefinitionsFoundInYAMLFile = "no definitions found in YAML file [%s]"
	logErrorCannotReadRootDirectory      = "cannot read root directory [%]"
	logDebugFoundDefinitionSubdirectory  = "found definition subdirectory [%s]"
	logErrorCannotProcessFiles           = "cannot process files in root directory [%s]"
	logWarnCannotProcessFile             = "cannot process files in directory [%s]"
	logDebugProcessingFile               = "processing file [%s] in directory [%s]"
	logDebugIgnoringFile                 = "ignoring file [%s] in directory [%s]"
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

	// Fields TODO
	Fields map[string]interface{}

	// Definition TODO
	Definition struct {
		Fields         Fields                   `yaml:"Fields"`
		References     []Reference              `yaml:"References"`
		SubDefinitions map[string]Specification `yaml:"SubDefinitions"`
	}

	// Specification TODO
	Specification struct {
		// Class allows the class for all of the definitions within the document to be specified.
		Class string `yaml:"Class,omitempty"`

		// References allows relationships to other classes for all of the definitions within the document to be
		// specified.
		References []Reference `yaml:"References,omitempty"`

		// Definitions TODO
		Definitions map[string]Definition `yaml:"Definitions,omitempty"`
	}

	processFileFuncType = func(filePath string, fileInfo os.FileInfo) (err error)
)

// LoadSpecificationFromFile TODO
func LoadSpecificationFromFile(filename string) (*Specification, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugCannotLoadYAMLFile, filename))
		return nil, err
	}

	spec := &Specification{}
	err = yaml.Unmarshal(yamlFile, spec)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugCannotParseYAMLFile, filename))

		return nil, err
	}

	// if no definitions are found, return an error and an nil Specification
	if len(spec.Definitions) == 0 {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugNoDefinitionsFoundInYAMLFile, filename))
		return nil, fmt.Errorf(logDebugNoDefinitionsFoundInYAMLFile, filename)
	}

	// TODO debug
	return spec, nil
}

// ProcessFiles TODO
func ProcessFiles(rootDir, fileExtension string, processFileFunc processFileFuncType) error {
	err := filepath.Walk(rootDir,
		func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Warn().Err(err).Msgf(logWarnCannotProcessFile, filePath)
				return err
			}
			if !fileInfo.IsDir() {
				if strings.HasSuffix(fileInfo.Name(), fmt.Sprintf(definitionFormat, fileExtension)) {
					log.Debug().Msg(fmt.Sprintf(logDebugProcessingFile, fileInfo.Name(), filePath))
					return processFileFunc(filePath, fileInfo)
				}

				log.Debug().Msg(fmt.Sprintf(logDebugIgnoringFile, fileInfo.Name(), filePath))
			}
			return nil
		})

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf(logErrorCannotProcessFiles, rootDir))
		return err
	}

	return nil
}

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
	"github.com/rs/zerolog"
)

const (
	appName    = "yaml-graph"
	appVersion = "0.3.7"

	commandRootUse      = appName
	commandRootUseShort = appName + ": generate graphs from YAML definition files"
	commandRootUseLong  = "Define data in YAML then generate graph representations to model relationships"

	commandVersionUse    = "version"
	commandVersionShort  = "Print the version number of " + appName
	commandVersionString = appVersion

	commandLoadUse      = "load"
	commandLoadUseShort = "Load definition files into graph representation"

	commandGraphUse      = "graph"
	commandGraphUseShort = "Generate HTML graph from definition files"

	commandValidateUse      = "validate"
	commandValidateUseShort = "Validate definition files"

	commandJSONUse      = "json"
	commandJSONUseShort = "generate JSON tree"

	commandReportUse      = "report"
	commandReportUseShort = "Generate report from graph representation"

	flagFileExtension          = "ext"
	flagFileExtensionShorthand = "e"
	flagFileExtensionDefault   = "yaml"
	flagFileExtensionUsage     = "file extension for definitions"

	flagDBURLName      = "dbURL"
	flagDBURLShorthand = "d"
	flagDBURLDefault   = "bolt://localhost:7687"
	flagDBURLUsage     = "URL of graph database"

	flagUsernameName      = "username"
	flagUsernameShorthand = "u"
	flagUsernameDefault   = "username"
	flagUsernameUsage     = "username for graph database"

	flagPasswordName      = "password"
	flagPasswordShorthand = "p"
	flagPasswordDefault   = "password"
	flagPasswordUsage     = "password for graph database"

	flagLogLevelName      = "logLevel"
	flagLogLevelShorthand = "l"
	flagLogLevelDefault   = int8(zerolog.WarnLevel)
	flagLogLevelUsage     = "log level (0=debug, 1=info, 2=warn, 3=error)"

	flagSourceName      = "source"
	flagSourceShorthand = "s"
	flagSourceUsage     = "Source directory to read definitions from (required)"
	flagSourceDefault   = "definition"

	flagReportFieldsFileName      = "fields"
	flagReportFieldsFileShorthand = "f"
	flagReportFieldsFileUsage     = "report fields file (required)"

	flagJSONDefinitionName      = "json"
	flagJSONDefinitionShorthand = "j"
	flagJSONDefinitionUsage     = "JSON definition file (required)"

	flagReportTemplateFileName      = "template"
	flagReportTemplateFileShorthand = "t"
	flagReportTemplateFileUsage     = "report template file (required)"

	flagDefinitionFormatName      = "format"
	flagDefinitionFormatShorthand = "f"
	flagDefinitionFormatUsage     = "Definition format file (required)"

	flagLoadDefinitionsName  = "load"
	flagLoadDefinitionsUsage = "load definitions"

	exitCodeRootCmdFailed     = 1
	exitCodeLoadCmdFailed     = 2
	exitCodeValidateCmdFailed = 3
	exitCodeJSONCmdFailed     = 4
	exitCodeTemplateCmdFailed = 5
)

var (
	// variable for flagJSONDefinitionName parameter
	jsonDefinition string

	// variable for flagSourceName parameter
	// note: we allow multiple source directories to enable the union of two definition directories
	sourceDir []string

	// variable for flagReportDefinitionName parameter
	reportDefinition string

	// variable for flagFileExtension parameter
	fileExtension string

	// variable for flagUsernameName parameter
	username string

	// variable for flagPasswordName parameter
	password string

	// variable for flagDBURLName parameter
	dbURL string

	// variable for flagLogLevelName parameter
	logLevel int8

	// variable for flagReportTemplateFileName parameter
	templateName string

	// variable for flagReportFieldsFileName parameter
	templateFormat string

	// variable for flagLoadDefinitionsName parameter
	loadDefinitions bool

	// variable for flagReportFieldsFileName parameter
	definitionFormatFile string
)

package cmd

import (
	"github.com/rs/zerolog"
)

const (
	appName    = "yaml-graph"
	appVersion = "0.1"

	commandRootUse      = appName
	commandRootUseShort = appName + ": generate graphs from YAML definition files"
	commandRootUseLong  = "Define data in YAML then generate graph representations to model relationships"

	commandVersionUse   = "version"
	commandVersionShort = "Print the version number of " + appName
	commandVersionLong  = appName + " " + appVersion

	commandParseUse      = "parse"
	commandParseUseShort = "Parse definition files into graph representation"

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

	exitCodeRootCmdFailed  = 1
	exitCodeParseCmdFailed = 2
)

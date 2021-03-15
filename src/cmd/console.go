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
	"github.com/nextmetaphor/yaml-graph/cui"
	"github.com/spf13/cobra"
)

var (
	consoleCmd = &cobra.Command{
		Use:   commandConsoleUse,
		Short: commandConsoleUseShort,
		Run:   console,
	}
)

func init() {
	rootCmd.AddCommand(consoleCmd)
}

func console(_ *cobra.Command, _ []string) {
	cui.OpenConsole(dbURL, username, password)
}

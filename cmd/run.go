// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/palantir/godel-format-plugin/formatplugin"
	"github.com/palantir/godel-format-plugin/formatplugin/config"
	godelconfig "github.com/palantir/godel/v2/framework/godel/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [flags] [files]",
	Short: "Format specified files (if no files are specified, format all project Go files)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var formatCfg config.Format
		if formatConfigFileFlagVal != "" {
			cfgVal, err := readFormatConfigFromFile(formatConfigFileFlagVal)
			if err != nil {
				return err
			}
			formatCfg = cfgVal
		}

		if godelConfigFileFlagVal != "" {
			exclude, err := godelconfig.ReadGodelConfigExcludesFromFile(godelConfigFileFlagVal)
			if err != nil {
				return err
			}
			formatCfg.Exclude.Add(exclude)
		}

		param, err := formatCfg.ToParam(cliFormatterFactory)
		if err != nil {
			return err
		}

		// no formatters specified
		if len(assetsFlagVal) == 0 {
			return nil
		}
		return formatplugin.Run(param, projectDirFlagVal, verifyFlagVal, args, cmd.OutOrStdout())
	},
}

func init() {
	runCmd.Flags().BoolVar(&verifyFlagVal, "verify", false, "verify files match formatting without applying formatting")
	rootCmd.AddCommand(runCmd)
}

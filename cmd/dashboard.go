/*
Copyright © 2023 Julian Easterling <julian@julianscorner.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:                "dashboard",
	Short:              "Open the Kubernetes development environment dashboard",
	Long:               "Open the Kubernetes development environment dashboard",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if isMinikubeEnv() {
			cobra.CheckErr(runMinikube("dashboard"))
		}

		if isVagrantEnv() {
			cobra.CheckErr(executeExternalProgram("octant", "--kubeconfig=./.kubectl.cfg"))
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		cobra.CheckErr(ensureKubectlfile())
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

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
	"errors"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Output status of the Kubernetes development environment",
	Long:  "Output status of the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		if isMinikubeEnv() {
			cobra.CheckErr(runMinikube("status"))
		}

		if isVagrantEnv() {
			cobra.CheckErr(executeExternalProgram("vagrant", "status"))
		}

		err := ensureKubectlfile()

		if err == nil {
			nodesCmd.Run(cmd, args)
			podsCmd.Run(cmd, args)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if !isVagrantEnv() && !isMinikubeEnv() {
			cobra.CheckErr(errors.New("a Kubernetes development environment not found"))
		}

		if isVagrantEnv() {
			cobra.CheckErr(ensureVagrantfile())
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

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

var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Aliases: []string{"down"},
	Short:   "Destroy the Kubernetes development environment",
	Long:    "Destroy the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			force = askForConfirmation("Are you sure you want to destroy the environment?")
		}

		if force {
			if isVagrantEnv() {
				cobra.CheckErr(vagrantDestroy())
			}

			if isMinikubeEnv() {
				cobra.CheckErr(minikubeDestroy())
			}
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
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().BoolP("force", "f", false, "destroy without confirmation")
}

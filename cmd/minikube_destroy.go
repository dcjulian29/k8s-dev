/*
Copyright © 2024 Julian Easterling <julian@julianscorner.com>

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

var minikube_destroyCmd = &cobra.Command{
	Use:     "destroy",
	Aliases: []string{"down"},
	Short:   "Destroy the Kubernetes development Minikube environment",
	Long:    "Destroy the Kubernetes development Minikube environment",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			force = askForConfirmation("Are you sure you want to destroy the environment?")
		}

		if force {
			cobra.CheckErr(removeMinikube())
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if !isMinikubeEnv() {
			cobra.CheckErr(errors.New("the Minikube environment does not exist"))
		}

		if !isMinikubeRunning() {
			cobra.CheckErr(errors.New("the Minikube environment is not running"))
		}
	},
}

func init() {
	minikubeCmd.AddCommand(minikube_destroyCmd)

	minikube_destroyCmd.Flags().BoolP("force", "f", false, "destroy without confirmation")
}

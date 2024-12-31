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

var resetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"recreate"},
	Short:   "Reset the Kubernetes development vagrant environment",
	Long:    "Reset the Kubernetes development vagrant environment",
	Run: func(cmd *cobra.Command, args []string) {
		recreate, _ := cmd.Flags().GetBool("recreate")

		if recreate {
			cobra.CheckErr(vagrantDestroy())

			provision, _ := cmd.Flags().GetBool("provision")

			cobra.CheckErr(vagrantUp("", provision))

			cobra.CheckErr(executeExternalProgram("ansible-playbook", "playbooks/init.yml"))
		} else {
			cobra.CheckErr(executeExternalProgram("ansible-playbook", "playbooks/reset.yml"))
		}

		deploy, _ := cmd.Flags().GetBool("deploy")

		if deploy {
			deployCmd.Run(cmd, args)
		}

		nodes, _ := cmd.Flags().GetBool("nodes")

		if nodes {
			nodesCmd.Run(cmd, args)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if isMinikubeEnv() {
			cobra.CheckErr(errors.New("'reset' is only available for a vagrant environment"))
		}

		if !isVagrantEnv() {
			cobra.CheckErr(errors.New("the vagrant environment does not exist"))
		} else {
			cobra.CheckErr(ensureVagrantfile())
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolP("provision", "p", true, "run the Vagrant provisioner")
	resetCmd.Flags().Bool("recreate", false, "Recreate the vagrant hosts")
	resetCmd.Flags().BoolP("nodes", "n", true, "Show nodes of development environment")
	resetCmd.Flags().BoolP("deploy", "d", false, "deploy the Kubernetes cluster")
}

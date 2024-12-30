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
	"fmt"
	"strings"

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
			cobra.CheckErr(vagrantUp(strings.Join(args, " "), true))
			cobra.CheckErr(executeExternalProgram("ansible-playbook", "playbooks/init.yml"))
		} else {
			cobra.CheckErr(executeExternalProgram("ansible-playbook", "playbooks/reset.yml"))
		}

		nodes, _ := cmd.Flags().GetBool("nodes")
		pods, _ := cmd.Flags().GetBool("pods")

		if nodes {
			output, err := executeCommand("kubectl", "--kubeconfig=./.kubectl.cfg", "get", "nodes")

			cobra.CheckErr(err)

			printSubMessage("Cluster Nodes")
			fmt.Println(output)
		}

		if pods {
			output, err := executeCommand("kubectl", "--kubeconfig=.kubectl.cfg", "get", "pods", "--all-namespaces")

			cobra.CheckErr(err)

			printSubMessage("Cluster Pods")
			fmt.Println(output)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		cobra.CheckErr(ensureVagrantfile())
		cobra.CheckErr(ensureKubectlfile())
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolP("nodes", "n", true, "Show nodes of deployed cluster")
	resetCmd.Flags().BoolP("pods", "p", false, "Show pods of deployed cluster")
	resetCmd.Flags().Bool("recreate", false, "Recreate the vagrant hosts")
}

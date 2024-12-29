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

	"github.com/spf13/cobra"
)

var miniKube_deploy bool = false

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the Kubernetes development environment",
	Long:  "Deploy the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		if miniKube_deploy {
			env := []string{
				"K8S_AUTH_VERIFY_SSL=false",
			}

			cobra.CheckErr(executeExternalProgramEnv("ansible-playbook", env, "playbooks/config.yml"))
		} else {
			cobra.CheckErr(executeExternalProgram("ansible-playbook", "playbooks/init.yml"))
		}

		nodes, _ := cmd.Flags().GetBool("nodes")
		pods, _ := cmd.Flags().GetBool("pods")

		if nodes {
			output, err := executeCommand("kubectl", "--kubeconfig=./.kubectl.cfg", "get", "nodes")

			cobra.CheckErr(err)

			printSubMessage("Environment Nodes")
			fmt.Println(output)
		}

		if pods {
			output, err := executeCommand("kubectl", "--kubeconfig=./.kubectl.cfg", "get", "pods", "--all-namespaces")

			cobra.CheckErr(err)

			printSubMessage("Environment Pods")
			fmt.Println(output)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		cobra.CheckErr(ensureVagrantfile())

		err := ensureMinikubeRunning()

		if err != nil {
			// Vagrant environments initialize Kubernetes in the VM(s) via playbook
			// so deployment would have to be forced if already deployed.
			err = ensureKubectlfile()

			if err == nil {
				force, _ := cmd.Flags().GetBool("force")

				if !force {
					cobra.CheckErr(fmt.Errorf("%s has already been deployed", "kubernetes"))
				}
			}
		} else {
			miniKube_deploy = true
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().BoolP("nodes", "n", true, "Show nodes of deployed environment")
	deployCmd.Flags().BoolP("pods", "p", false, "Show pods of deployed environment")
	deployCmd.Flags().BoolP("force", "f", false, "force redeployment of environment")
}

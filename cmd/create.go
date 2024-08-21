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

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"up"},
	Short:   "Create the Kubernetes development Vagrant environment",
	Long:    "Create the Kubernetes development Vagrant environment",
	Run: func(cmd *cobra.Command, args []string) {
		provision, _ := cmd.Flags().GetBool("provision")
		cobra.CheckErr(vagrantUp(strings.Join(args, " "), provision))

		deploy, _ := cmd.Flags().GetBool("deploy")

		if deploy {
			deployCmd.Run(cmd, args)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		cobra.CheckErr(ensureVagrantfile())
		err := ensureKubectlfile()

		if err == nil {

			printMessage("kubernetes cluster is currently deployed")

			force, _ := cmd.Flags().GetBool("force")

			var destroy bool

			if force {
				destroy = true
			} else {
				destroy = askForConfirmation("Are you sure you want to recreate the cluster?")
			}

			if destroy {
				vagrantDestroy(strings.Join(args, " "), true)
				removeFile(".kubectl.cfg")
			} else {
				cobra.CheckErr(fmt.Errorf("cluster cannot be created while it exists"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolP("provision", "p", true, "run the Vagrant provisioner")
	createCmd.Flags().BoolP("deploy", "d", false, "deploy the Kubernetes cluster")
	createCmd.Flags().BoolP("force", "f", false, "force recreation of the Kubernetes cluster")
}

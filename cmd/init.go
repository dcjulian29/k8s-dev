/*
Copyright © 2023 Julian Easterling julian@julianscorner.com

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

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize"},
	Short:   "Initialize the Kubernetes development vagrant environment files",
	Long:    "Initialize the Kubernetes development vagrant environment files",
	Run: func(cmd *cobra.Command, args []string) {
		if workingDirectory != folderPath {
			cobra.CheckErr(createFolder(folderPath))
		}

		ensureRootDirectory()

		force, _ := cmd.Flags().GetBool("force")

		cobra.CheckErr(testNeedForce(force))

		printMessage("Initializing the development vagrant environment...")

		cobra.CheckErr(makeDirectory("collections"))
		cobra.CheckErr(makeDirectory("group_vars"))
		cobra.CheckErr(makeDirectory("playbooks"))
		cobra.CheckErr(makeDirectory("roles"))
		cobra.CheckErr(ansible_cfg())
		cobra.CheckErr(ansible_lint())

		servers, _ := cmd.Flags().GetInt("servers")
		agents, _ := cmd.Flags().GetInt("agents")
		box, _ := cmd.Flags().GetString("box")

		cobra.CheckErr(inventory_file(servers, agents))
		cobra.CheckErr(requirements_yml())
		cobra.CheckErr(vagrant_file(servers, agents, box))

		cobra.CheckErr(all_yml())
		cobra.CheckErr(k3s_cluster_yml())

		cobra.CheckErr(deploy_yml())
		cobra.CheckErr(reboot_yml())
		cobra.CheckErr(reset_yml())

		params := []string{"install", "-r", "./requirements.yml"}

		printSubMessage("restoring collections")
		cobra.CheckErr(executeExternalProgram("ansible-galaxy", append([]string{"collection"}, params...)...))

		printSubMessage("restoring roles")
		cobra.CheckErr(executeExternalProgram("ansible-galaxy", append([]string{"role"}, params...)...))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().IntP("agents", "a", 3, "number of work agent nodes")
	initCmd.Flags().IntP("servers", "s", 2, "number of control servers")
	initCmd.Flags().StringP("box", "b", "debian/bookworm64", "vagrant box image")

	initCmd.Flags().BoolP("force", "f", false, "overwrite an existing development folder")
}

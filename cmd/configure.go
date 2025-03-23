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

var configCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config"},
	Short:   "Configure the Kubernetes development environment",
	Long:    "Configure the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		env := []string{
			"K8S_AUTH_VERIFY_SSL=false",
		}

		cobra.CheckErr(executeExternalProgramEnv("ansible-playbook", env, "playbooks/config.yml"))

		pods, _ := cmd.Flags().GetBool("pods")

		if pods {
			podsCmd.Run(cmd, args)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		err := ensureKubectlfile()

		if err == nil {
			cobra.CheckErr(fmt.Errorf("%s has not been created yet", "kubernetes"))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolP("pods", "p", false, "Show pods of the configured environment")
}

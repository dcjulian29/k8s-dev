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

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Show pods of development environment",
	Long:  "Show pods of development environment",
	Run: func(cmd *cobra.Command, args []string) {
		params := []string{
			"--kubeconfig=./.kubectl.cfg",
			"--insecure-skip-tls-verify=true",
			"get",
			"pods",
		}

		namespace, _ := cmd.Flags().GetString("namespace")

		if len(namespace) == 0 {
			params = append(params, "--all-namespaces")
		} else {
			params = append(params, "--namespace", namespace)
		}

		output, err := executeCommand("kubectl", params...)

		cobra.CheckErr(err)

		fmt.Println(output)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		cobra.CheckErr(ensureVagrantfile())
		cobra.CheckErr(ensureKubectlfile())
	},
}

func init() {
	rootCmd.AddCommand(podsCmd)

	podsCmd.Flags().StringP("namespace", "n", "", "Show pods of the specified namespace")
}

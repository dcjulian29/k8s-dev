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

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Use helm in the Kubernetes development environment",
	Long:  "Use helm in the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		params := append(args, "--kubeconfig=./.kubectl.cfg", "--kube-context=default")

		output, err := executeCommand("helm", params...)

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
	rootCmd.AddCommand(helmCmd)
}

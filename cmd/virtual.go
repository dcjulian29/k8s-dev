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
	"fmt"

	"github.com/spf13/cobra"
)

var virtualCmd = &cobra.Command{
	Use:                "virtual",
	Short:              "Control virtual machines in the Kubernetes development vagrant environment",
	Long:               "Control virtual machines in the Kubernetes development vagrant environment",
	Aliases:            []string{"virtctl"},
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		params := append(args, "--kubeconfig=./.kubectl.cfg")

		output, err := executeCommand("virtctl", params...)

		cobra.CheckErr(err)

		fmt.Println(output)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if isMinikubeEnv() {
			cobra.CheckErr(errors.New("'virtual' is only available for a vagrant environment"))
		}

		if !isVagrantEnv() {
			cobra.CheckErr(errors.New("the vagrant environment does not exist"))
		} else {
			cobra.CheckErr(ensureVagrantfile())
		}

		cobra.CheckErr(ensureKubectlfile())
	},
}

func init() {
	rootCmd.AddCommand(virtualCmd)
}

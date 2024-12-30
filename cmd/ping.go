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

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the Kubernetes development vagrant environment",
	Long:  "Ping the Kubernetes development vagrant environment",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(executeExternalProgram("vagrant", "provision"))
		cobra.CheckErr(executeExternalProgram("ansible", "-m", "ping", "all"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
		if isMinikubeEnv() {
			cobra.CheckErr(errors.New("the ping command isn't applicable to a minikube environment"))
		}

		cobra.CheckErr(ensureVagrantfile())
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}

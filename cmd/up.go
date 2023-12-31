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
	"strings"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [node]",
	Short: "Bring the Kubernetes development vagrant environment online",
	Long:  "Bring the Kubernetes development vagrant environment online",
	Run: func(cmd *cobra.Command, args []string) {
		provision, _ := cmd.Flags().GetBool("provision")
		vagrantUp(strings.Join(args, " "), provision)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolP("provision", "p", true, "run the vagrant provisioner")
}

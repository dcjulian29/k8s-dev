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

var destroyCmd = &cobra.Command{
	Use:   "destroy [node]",
	Short: "Destroy the Kubernetes development vagrant environment",
	Long:  "Destroy the Kubernetes development vagrant environment",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		vagrantDestroy(strings.Join(args, " "), force)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().BoolP("force", "f", false, "destroy without confirmation")
}

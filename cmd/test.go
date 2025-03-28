/*
Copyright © 2024 Julian Easterling <julian@julianscorner.com>

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
	"os"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <role>",
	Args:  cobra.ExactArgs(1),
	Short: "Test a role against Kubernetes development environment",
	Long:  "Test a role against Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		step, _ := cmd.Flags().GetBool("step")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if fileExists(".tmp/play.yml") {
			cobra.CheckErr(os.Remove(".tmp/play.yml"))
		}

		if !dirExists(".tmp") {
			cobra.CheckErr(os.Mkdir(".tmp", 0755))
		}

		file, err := os.Create(".tmp/play.yml")

		if err != nil {
			cobra.CheckErr(err)
		}

		defer file.Close()

		content := "---\n- hosts: 127.0.0.1\n  gather_facts: false\n\n  vars:\n    k8s_config: ../.kubectl.cfg\n\n  roles:\n"
		content = fmt.Sprintf("%s%s", content, fmt.Sprintf("    - role: %s\n", name))

		if _, err = file.WriteString(content); err != nil {
			cobra.CheckErr(err)
		}

		var param []string

		if verbose {
			param = append(param, "-v")
		}

		if step {
			param = append(param, "--step")
		}

		param = append(param, ".tmp/play.yml")

		cobra.CheckErr(executeExternalProgram("ansible-playbook", param...))

		cobra.CheckErr(os.RemoveAll(".tmp"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if !isVagrantEnv() && !isMinikubeEnv() {
			cobra.CheckErr(errors.New("a Kubernetes development environment not found"))
		}

		if isVagrantEnv() {
			cobra.CheckErr(ensureVagrantfile())
		}

		cobra.CheckErr(ensureKubectlfile())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolP("verbose", "v", false, "tell Ansible to print more debug messages")
	testCmd.Flags().Bool("step", false, "one-step-at-a-time: confirm each task before running")
}

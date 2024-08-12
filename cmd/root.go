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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.szostok.io/version/extension"
)

var (
	folderPath       string
	workingDirectory string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-dev",
	Short: "k8s-dev enables development of Kubernetes charts and Ansible playbooks and roles.",
	Long: `k8s-dev integrates with Vagrant to enable users to define, develop, and test Helm charts
and Ansible playbooks and roles. It allows users to define and manage infrastructure resources and
uses the providers automation engine to provision, develop, and test a repeatable environment.`,
}

func Execute() {
	workingDirectory, _ = os.Getwd()

	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "k8s-dev"),
		),
	)

	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	pwd, _ := os.Getwd()

	rootCmd.PersistentFlags().StringVar(&folderPath, "path", pwd, "path to development folder")
}

func ensureKubectlfile() error {
	if !fileExists(".kubectl.cfg") {
		return fmt.Errorf("can't find the kubectl.cfg file")
	}

	return nil
}

func ensureRootDirectory() {
	if workingDirectory != folderPath {
		err := os.Chdir(folderPath)
		cobra.CheckErr(err)
	}
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", s)

		response, _ := reader.ReadString('\n')

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" || response == "" {
			return false
		}
	}
}

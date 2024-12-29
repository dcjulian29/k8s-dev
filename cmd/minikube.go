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
	b64 "encoding/base64"
	"os"

	"github.com/spf13/cobra"
)

var minikubeCmd = &cobra.Command{
	Use:   "minikube",
	Short: "Manage creation and destruction of the Minikube Kubernetes environment",
	Long:  "Manage creation and destruction of the Minikube Kubernetes environment",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(minikubeCmd)
}

func isMinikubeEnv() bool {
	return dirExists("./.minikube")
}

func isMinikubeRunning() bool {
	if !isMinikubeEnv() {
		return false
	}

	env := []string{
		"KUBECONFIG=./.kubectl.cfg",
	}

	_, err := executeCommandEnv("minikube", env, "status")

	if err != nil {
		return err.Error() == "exit status 4"
	}

	return true
}

func removeMinikube() error {
	err := removeFile("./.kubectl.cfg")

	if err != nil {
		return err
	}

	err = executeExternalProgram("minikube", "stop")

	if err != nil {
		return err
	}

	err = executeExternalProgram("minikube", "delete")

	if err != nil {
		return err
	}

	return os.RemoveAll("./.minikube")
}

func convertToBase64(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

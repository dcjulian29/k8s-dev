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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var minikube_createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"up"},
	Short:   "Create the Kubernetes development Minikube environment",
	Long:    "Create the Kubernetes development Minikube environment",
	Run: func(cmd *cobra.Command, args []string) {
		cni, _ := cmd.Flags().GetString("cni")
		nodes, _ := cmd.Flags().GetInt("total-nodes")
		ha, _ := cmd.Flags().GetBool("ha")

		cni = strings.ToLower(cni)
		cobra.CheckErr(validateNetwork(cni))

		param := []string{
			"start",
			"--cni=" + cni,
			"--listen-address=0.0.0.0",
		}

		if nodes > 0 {
			param = append(param, "--nodes="+strconv.Itoa(nodes))
		}

		if ha {
			param = append(param, "--ha")
		}

		env := []string{
			"KUBECONFIG=./.kubectl.cfg",
		}

		cobra.CheckErr(executeExternalProgramEnv("minikube", env, param...))

		configureKubectl()

		deploy, _ := cmd.Flags().GetBool("deploy")

		if deploy {
			deployCmd.Run(cmd, args)
		}

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if isMinikubeRunning() {
			printMessage("kubernetes is currently deployed")

			force, _ := cmd.Flags().GetBool("force")

			var destroy bool

			if force {
				destroy = true
			} else {
				destroy = askForConfirmation("Are you sure you want to recreate the environment?")
			}

			if destroy {
				cobra.CheckErr(removeMinikube())
			} else {
				cobra.CheckErr(fmt.Errorf("environment cannot be created while it exists"))
			}
		}
	},
}

func init() {
	minikubeCmd.AddCommand(minikube_createCmd)

	minikube_createCmd.Flags().String("cni", "flannel", "CNI plug-in to use. Valid options: calico, cilium, flannel")
	minikube_createCmd.Flags().IntP("total-nodes", "t", 0, "number of total nodes to create ('0' indicates auto)")
	minikube_createCmd.Flags().Bool("ha", false, "create highly available multi-control plane")
	minikube_createCmd.Flags().BoolP("deploy", "d", false, "deploy the Kubernetes cluster")
	minikube_createCmd.Flags().BoolP("force", "f", false, "force recreation of the Kubernetes cluster")
}

func validateNetwork(provider string) error {
	validNetworks := []string{"calico", "cilium", "flannel"}

	for _, validNetwork := range validNetworks {
		if provider == validNetwork {
			return nil
		}
	}

	return fmt.Errorf("'%s' is invalid for CNI. Valid options: calico, cilium, flannel", provider)
}

func configureKubectl() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	current, err := os.Getwd()
	cobra.CheckErr(err)

	cobra.CheckErr(ensureDir(makePath(current, ".minikube")))

	src := makePath(home, ".minikube", "ca.crt")
	cobra.CheckErr(copyFile(src, makePath(current, makePath(".minikube", "ca.crt"))))
	content, err := readFile("./.minikube/ca.crt")
	cobra.CheckErr(err)

	cobra.CheckErr(replaceInFile("./.kubectl.cfg", "certificate-authority: ", "certificate-authority-data: "))
	cobra.CheckErr(replaceInFile("./.kubectl.cfg", src, convertToBase64(content)))

	src = makePath(home, ".minikube", "profiles", "minikube", "client.crt")
	cobra.CheckErr(copyFile(src, makePath(current, makePath(".minikube", "client.crt"))))
	content, err = readFile("./.minikube/client.crt")
	cobra.CheckErr(err)

	cobra.CheckErr(replaceInFile("./.kubectl.cfg", "client-certificate: ", "client-certificate-data: "))
	cobra.CheckErr(replaceInFile("./.kubectl.cfg", src, convertToBase64(content)))

	src = makePath(home, ".minikube", "profiles", "minikube", "client.key")
	cobra.CheckErr(copyFile(src, makePath(current, makePath(".minikube", "client.key"))))
	content, err = readFile("./.minikube/client.key")
	cobra.CheckErr(err)

	cobra.CheckErr(replaceInFile("./.kubectl.cfg", "client-key: ", "client-key-data: "))
	cobra.CheckErr(replaceInFile("./.kubectl.cfg", src, convertToBase64(content)))

	ip, err := GetHostIP()
	cobra.CheckErr(err)

	cobra.CheckErr(replaceInFile("./.kubectl.cfg", "127.0.0.1", ip))
}

// func toIndentYaml(content, indent string) string {
// 	lines := strings.Split(content, "\n")

// 	for i, line := range lines {
// 		lines[i] = fmt.Sprintf("%s%s", indent, line)
// 	}

// 	output := strings.Join(lines, "\n")

// 	return output
// }

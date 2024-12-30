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
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"up"},
	Short:   "Create the Kubernetes development environment",
	Long:    "Create the Kubernetes development environment",
	Run: func(cmd *cobra.Command, args []string) {
		minikube, _ := cmd.Flags().GetBool("minikube")

		if minikube {
			cni, _ := cmd.Flags().GetString("cni")
			nodes, _ := cmd.Flags().GetInt("total-nodes")
			ha, _ := cmd.Flags().GetBool("ha")

			if len(cni) == 0 {
				cni = "flannel"
			} else {
				cni = strings.ToLower(cni)
			}

			cobra.CheckErr(validateMinikubeNetwork(cni))

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

			cobra.CheckErr(runMinikube(param...))

			configureMinikubeKubectl()
		} else {
			provision, _ := cmd.Flags().GetBool("provision")
			cobra.CheckErr(vagrantUp(strings.Join(args, " "), provision))
		}

		deploy, _ := cmd.Flags().GetBool("deploy")

		if deploy {
			deployCmd.Run(cmd, args)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureRootDirectory()

		if isMinikubeEnv() || isVagrantEnv() {
			cobra.CheckErr(errors.New("environment is already created"))
		}

		vagrant, _ := cmd.Flags().GetBool("vagrant")
		minikube, _ := cmd.Flags().GetBool("minikube")
		needForce := false

		if minikube {
			printMessage("using Minikube environment...")
			if isMinikubeRunning() {
				needForce = true
			}
		} else {
			if vagrant {
				ha, _ := cmd.Flags().GetBool("ha")

				if ha {
					cobra.CheckErr(errors.New("'ha' is not valid in a Vagrant environment"))
				}

				nodes, _ := cmd.Flags().GetInt("total-nodes")

				if nodes > 0 {
					cobra.CheckErr(errors.New("'total-nodes' is not valid in a Vagrant environment"))
				}

				cni, _ := cmd.Flags().GetString("cni")

				if len(cni) > 0 {
					cobra.CheckErr(errors.New("'cni' is not valid in a Vagrant environment"))
				}

				printMessage("using Vagrant environment...")

				cobra.CheckErr(ensureVagrantfile())

				err := ensureKubectlfile()

				if err == nil {
					needForce = true
				}
			} else {
				cobra.CheckErr(errors.New("environment not selected"))
			}
		}

		if needForce {
			printSubMessage("kubernetes environment is currently deployed")

			force, _ := cmd.Flags().GetBool("force")

			var destroy bool

			if force {
				destroy = true
			} else {
				destroy = askForConfirmation("Are you sure you want to recreate the environment?")
			}

			if destroy {
				if minikube {
					cobra.CheckErr(minikubeDestroy())
				} else {
					cobra.CheckErr(vagrantDestroy())
				}
			} else {
				cobra.CheckErr(errors.New("environment cannot be created while it exists"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Bool("vagrant", true, "Use Vagrant to create environment")
	createCmd.Flags().BoolP("provision", "p", true, "run the Vagrant provisioner")

	createCmd.Flags().Bool("minikube", false, "Use Minikube to create environment")
	createCmd.Flags().String("cni", "", "CNI plug-in to use. Valid options: calico, cilium, flannel")
	createCmd.Flags().IntP("total-nodes", "t", 0, "number of total nodes to create ('0' indicates auto)")
	createCmd.Flags().Bool("ha", false, "create highly available multi-control plane")

	createCmd.Flags().BoolP("deploy", "d", false, "deploy the Kubernetes cluster")
	createCmd.Flags().BoolP("force", "f", false, "force recreation of the Kubernetes cluster")
}

func configureMinikubeKubectl() {
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

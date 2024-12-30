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
)

func ensureVagrantfile() error {
	if !fileExists("Vagrantfile") {
		return fmt.Errorf("can't find the Vagrantfile")
	}

	return nil
}

func isVagrantEnv() bool {
	return dirExists("./.vagrant")
}

func vagrantDestroy() error {
	fmt.Println(Info("destroying all Vagrant machines..."))

	err := removeFile("./.kubectl.cfg")

	if err != nil {
		return err
	}

	err = executeExternalProgram("vagrant", "destroy", "--force")

	if err != nil {
		return err
	}

	err = os.RemoveAll("./.vagrant")

	if err != nil {
		return err
	}

	return nil
}

func vagrantUp(name string, provision bool) error {
	if err := ensureVagrantfile(); err != nil {
		return err
	}

	param := []string{
		"up",
	}

	if len(name) > 0 {
		fmt.Println(Info(fmt.Sprintf("Bringing '%s' online...", name)))

		param = append(param, name)
	} else {
		fmt.Println(Info("Bringing all vagrant machines online..."))
	}

	if provision {
		param = append(param, "--provision")
	}

	return executeExternalProgram("vagrant", param...)
}

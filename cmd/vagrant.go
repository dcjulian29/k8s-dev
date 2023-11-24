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

import "fmt"

func vagrantHalt(name string, force bool) error {
	param := []string{
		"halt",
	}

	if len(name) > 0 {
		fmt.Println(Info(fmt.Sprintf("Halting '%s' machine...", name)))

		param = append(param, name)
	} else {
		fmt.Println(Info("Halting all vagrant machines..."))
	}

	return executeExternalProgram("vagrant", param...)
}

func vagrantUp(name string, provision bool) error {
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

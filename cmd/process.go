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
	"os"
	"os/exec"
)

func executeExternalProgram(program string, params ...string) error {
	return executeExternalProgramEnv(program, []string{""}, params...)
}

func executeExternalProgramEnv(program string, env []string, params ...string) error {
	cmd := exec.Command(program, params...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = append(os.Environ(), env...)

	return cmd.Run()
}

func executeCommand(program string, params ...string) (string, error) {
	return executeCommandEnv(program, []string{""}, params...)
}

func executeCommandEnv(program string, env []string, params ...string) (string, error) {
	cmd := exec.Command(program, params...)
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), env...)

	output, err := cmd.CombinedOutput()

	return string(output[:]), err
}

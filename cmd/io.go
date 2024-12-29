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
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
)

func copyFile(src, dst string) error {
	if fileExists(src) {
		if fileExists(dst) {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}

		source, err := os.Open(src)
		if err != nil {
			return err
		}

		defer source.Close()

		destination, err := os.Create(dst)

		if err != nil {
			return err
		}

		defer destination.Close()

		_, err = io.Copy(destination, source)

		if err != nil {
			return err
		}
	}

	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func makePath(parts ...string) string {
	path := ""

	p := "/"

	if runtime.GOOS == "windows" {
		p = "\\"
	}

	for _, part := range parts {
		if len(path) == 0 {
			path = part
		} else {
			path = path + p + part
		}
	}

	return path
}

func removeFile(filePath string) error {
	if fileExists(filePath) {
		return os.Remove(filePath)
	}

	return nil
}

func readFile(filePath string) (string, error) {
	if fileExists(filePath) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}

	return "", fmt.Errorf("'%s' does not exist", filePath)
}

func replaceInFile(filename, original, replacement string) error {
	input, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, original) {
			lines[i] = strings.Replace(lines[i], original, replacement, -1)
		}
	}

	output := strings.Join(lines, "\n")

	err = os.WriteFile(filename, []byte(output), 0644)

	if err != nil {
		return err
	}

	return nil
}

func GetHostIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()

	return localAddr, nil
}

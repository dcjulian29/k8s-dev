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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "create"},
	Short:   "Initialize the Kubernetes development vagrant environment",
	Long:    "Initialize the Kubernetes development vagrant environment",
	Run:     initialize,
}

func ansible_cfg() {
	printSubMessage("creating 'ansible.cfg'")

	file, err := os.Create("ansible.cfg")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`[defaults]
any_errors_fatal            = true
collections_path            = ./collections
duplicate_dict_key          = error
error_on_undefined_vars     = true
gathering                   = smart
host_key_checking           = false
inventory                   = ./hosts.ini
log_path                    = ./ansible.log
roles_path                  = ./roles:./k3s-ansible/roles
stdout_callback             = community.general.yaml
`)

	_, err = file.Write(content)

	cobra.CheckErr(err)
}

func ansible_lint() {
	printSubMessage("creating '.ansible-lint'")
	file, err := os.Create(".ansible-lint")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`---
enable_list:
  - args
  - empty-string-compare
  - no-log-password
  - no-same-owner
  - name[prefix]
  - yaml
exclude_paths:
  - roles/
  - k3s-ansible/
kinds:
  - playbook: "playbooks/*.yml"
profile: production
skip_list:
  - experimental
`)

	_, err = file.Write(content)

	cobra.CheckErr(err)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().IntP("agents", "a", 3, "number of work agent nodes")
	initCmd.Flags().IntP("servers", "s", 2, "number of control servers")
	initCmd.Flags().StringP("box", "b", "debian/bookworm64", "vagrant box image")

	initCmd.Flags().BoolP("force", "f", false, "overwrite an existing development environment")
}

func initialize(cmd *cobra.Command, args []string) {
	printMessage("Initializing development environment...")

	if workingDirectory != folderPath {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			printSubMessage("creating development environment folder")
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				cobra.CheckErr(fmt.Errorf("unable to create development environment folder"))
			}
		}

		ensureRootDirectory()
	}

	force, _ := cmd.Flags().GetBool("force")
	if fileExists("ansible.cfg") && !force {
		cobra.CheckErr(fmt.Errorf("The folder for the development environment already contains an environment and forceFlag was not provided."))
	}

	makeDirectory("collections")
	makeDirectory("group_vars")
	makeDirectory("playbooks")
	makeDirectory("roles")
	ansible_cfg()
	ansible_lint()

	servers, _ := cmd.Flags().GetInt("servers")
	agents, _ := cmd.Flags().GetInt("agents")
	box, _ := cmd.Flags().GetString("box")

	inventory_file(servers, agents)
	vagrant_file(servers, agents, box)
	k3s_cluster()
}

func inventory_file(servers, agents int) {
	printSubMessage("creating 'hosts.ini'")

	file, err := os.Create("hosts.ini")

	cobra.CheckErr(err)

	defer file.Close()

	var sb strings.Builder

	sb.WriteString("[k3s_cluster:children]\nmaster\nnode\n\n[master]\n")

	for i := 1; i <= servers; i++ {
		sb.WriteString(fmt.Sprintf("control-%d ansible_host=192.168.57.1%d\n", i, i))
	}

	sb.WriteString("\n[node]\n")

	for i := 1; i <= agents; i++ {
		sb.WriteString(fmt.Sprintf("work-%d ansible_host=192.168.57.2%d\n", i, i))
	}

	sb.WriteString("\n[all:vars]\n")
	sb.WriteString("ansible_user=vagrant\n")
	sb.WriteString("ansible_ssh_common_args='-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o CheckHostIP=no'\n")
	sb.WriteString("ansible_port=22\n")
	sb.WriteString("ansible_ssh_private_key_file=~/.ssh/insecure_private_key\n")

	content := []byte(sb.String())

	_, err = file.Write(content)

	cobra.CheckErr(err)
}

func k3s_cluster() {
	printSubMessage("creating 'group_vars/k3s_cluster.yml'")
	file, err := os.Create("group_vars/k3s_cluster.yml")

	cobra.CheckErr(err)

	defer file.Close()

	content := []byte(`---
apiserver_endpoint: "192.168.57.10"
custom_registries: false
extra_args: >-
  --flannel-iface={{ flannel_iface }}
  --node-ip={{ k3s_node_ip }}
extra_server_args: >-
  {{ extra_args }}
  {{ '--node-taint node-role.kubernetes.io/master=true:NoSchedule' if k3s_master_taint else '' }}
  --tls-san {{ apiserver_endpoint }}
  --disable servicelb
  --disable traefik
extra_agent_args: >-
  {{ extra_args }}
flannel_iface: "{{ 'enp0s8' if ansible_distribution == 'Ubuntu' else 'eth1' }}"
k3s_master_taint: "{{ true if groups['node'] | default([]) | length >= 1 else false }}"
k3s_node_ip: '{{ ansible_facts[flannel_iface]["ipv4"]["address"] }}'
k3s_token: Sup3rS3cr3t!
kube_vip_iface: "{{ flannel_iface }}"
k3s_version: v1.25.12+k3s1
kube_vip_tag_version: "v0.5.12"
log_destination: true
metal_lb_controller_tag_version: "v0.13.9"
metal_lb_ip_range: "192.168.57.30-192.168.57.99"
metal_lb_mode: layer2
metal_lb_speaker_tag_version: "v0.13.9"
metal_lb_type: native
proxmox_lxc_configure: false
systemd_dir: /etc/systemd/system
system_timezone: "American/New_York"
`)

	_, err = file.Write(content)

	cobra.CheckErr(err)
}

func makeDirectory(path string) {
	printSubMessage(fmt.Sprintf("creating '%s/'", path))

	err := ensureDir(path)
	cobra.CheckErr(err)
}

func vagrant_file(servers, agents int, box string) {
	printSubMessage("creating 'Vagrantfile'")

	file, err := os.Create("Vagrantfile")

	cobra.CheckErr(err)

	defer file.Close()

	filevars := fmt.Sprintf("IMAGE_NAME = \"%s\"\nSERVER_NUMBER = %d\nAGENT_NUMBER = %d\n\n", box, servers, agents)
	content := []byte(filevars + `Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provision "shell", inline: "ping -c 1 192.168.57.1"
  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.cpus = 1
    vb.check_guest_additions = false
    vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
    vb.customize [ "modifyvm", :id, "--graphicscontroller", "vmsvga"]
    vb.customize [ "modifyvm", :id, "--ioapic", "on"]
    vb.customize [ "modifyvm", :id, "--nicpromisc2", "allow-vms" ]
  end

  (1..SERVER_NUMBER).each do |i|
    config.vm.define "control-#{i}" do |c|
      c.vm.box = IMAGE_NAME
      c.vm.hostname = "control-#{i}"
      c.vm.network "private_network", ip: "192.168.57.#{i + 10}"
      c.vm.network :forwarded_port, guest: 22, host: "80#{i + 10}", id: 'ssh'
      c.vm.provider "virtualbox" do |vb|
        vb.memory = 1048
      end
    end
  end

  (1..AGENT_NUMBER).each do |i|
    config.vm.define "work-#{i}" do |c|
      c.vm.box = IMAGE_NAME
      c.vm.hostname = "work-#{i}"
      c.vm.network "private_network", ip: "192.168.57.#{i + 20}"
      c.vm.network :forwarded_port, guest: 22, host: "80#{i + 20}", id: 'ssh'
      c.vm.provider "virtualbox" do |vb|
        vb.memory = 2048
      end
    end
  end
end
`)

	_, err = file.Write(content)

	cobra.CheckErr(err)
}

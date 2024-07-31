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
)

func ansible_cfg() error {
	return createFile("ansible.cfg", []byte(`[defaults]
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
`))
}

func ansible_lint() error {
	return createFile(".ansible-lint", []byte(`---
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
`))
}

func createFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		printSubMessage("creating development folder...")
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			return fmt.Errorf("unable to create development folder")
		}
	}

	return nil
}

func createFile(filePath string, content []byte) error {
	printSubMessage(fmt.Sprintf("creating '%s'", filePath))
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(content)

	if err != nil {
		return err
	}

	return nil
}

func init_yml() error {
	return createFile("playbooks/init.yml", []byte(`---
- name: Initialize node
  hosts: k3s_cluster
  become: true

  roles:
    - role: dcjulian29.base

  post_tasks:
    - name: Ensure the SSH user doesn't need password to sudo
      ansible.builtin.lineinfile:
        dest: /etc/sudoers
        regexp: '^{{ ansible_user }}'
        line: '{{ ansible_user }} ALL=(ALL) NOPASSWD: ALL'
        state: present
        validate: 'visudo -cf %s'
        mode: "0644"

    - name: Reboot and wait each node to come back up
      ansible.builtin.reboot:

- name: Initialize new kubernetes cluster
  hosts: k3s_cluster
  become: true

  roles:
    - role: prereq
    - role: download
    - role: k3s_custom_registries
      when: custom_registries

- name: Setup kubernetes servers
  hosts: master
  become: true

  roles:
    - role: k3s_server

- name: Setup kubernetes agents
  hosts: node
  become: true

  roles:
    - role: k3s_agent

- name: Configure k3s cluster
  hosts: master
  become: true

  roles:
    - role: k3s_server_post

  post_tasks:
    - name: Retrieve kubernetes configuration for this cluster
      ansible.builtin.fetch:
        src: ~/.kube/config
        dest: ../.kubectl.cfg
        flat: true
      when:
        - ansible_hostname == hostvars[groups['master'][0]]['ansible_hostname']
`))
}

func inventory_file(servers, agents int) error {
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

	return createFile("hosts.ini", []byte(sb.String()))
}

func all_yml() error {
	return createFile("group_vars/all.yml", []byte(`---
`))
}

func k3s_cluster_yml() error {
	return createFile("group_vars/k3s_cluster.yml", []byte(`---
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
`))
}

func makeDirectory(path string) error {
	printSubMessage(fmt.Sprintf("creating '%s/'", path))

	return ensureDir(path)
}

func reboot_yml() error {
	return createFile("playbooks/reboot.yml", []byte(`---
- name: Reboot kubernetes cluster
  hosts: k3s_cluster
  become: true
  gather_facts: false

  tasks:
    - name: Reboot and wait each node to come back up
      ansible.builtin.reboot:
`))
}

func reset_yml() error {
	return createFile("playbooks/reset.yml", []byte(`---
- name: Reset kubernetes cluster
  hosts: k3s_cluster
  become: true

  roles:
    - role: reset

  post_tasks:
    - name: Reboot and wait each node to come back up
      ansible.builtin.reboot:

- name: Initialize new kubernetes cluster
  hosts: k3s_cluster
  any_errors_fatal: true
  become: true

  roles:
    - role: prereq
    - role: download
    - role: k3s_custom_registries
      when: custom_registries

- name: Setup kubernetes servers
  hosts: master
  any_errors_fatal: true
  become: true

  roles:
    - role: k3s_server

- name: Setup kubernetes agents
  hosts: node
  any_errors_fatal: true
  become: true

  roles:
    - role: k3s_agent

- name: Configure k3s cluster
  hosts: master
  any_errors_fatal: true
  become: true

  roles:
    - role: k3s_server_post

  post_tasks:
    - name: Retrieve kubernetes configuration for this cluster
      ansible.builtin.fetch:
        src: ~/.kube/config
        dest: ../.kubectl.cfg
        flat: true

      when:
        - ansible_hostname == hostvars[groups['master'][0]]['ansible_hostname']
`))
}

func requirements_yml() error {
	return createFile("requirements.yml", []byte(`---
collections:
  - name: ansible.posix
  - name: ansible.utils
  - name: community.general
  - name: kubernetes.core
roles:
  - name: dcjulian29.base
    src: https://github.com/dcjulian29/ansible-role-base.git
`))
}

func testNeedForce(force bool) error {
	msg := "'%s' already exists, would be over written, and force was not provided."

	folders := []string{
		"collections",
		"group_vars",
		"playbooks",
		"roles",
		"k3s-ansible",
	}

	files := []string{
		"ansible.cfg",
		".ansible-lint",
		"host.ini",
		"requirements.yml",
		"Vagrantfile",
	}

	for _, f := range folders {
		if dirExists(f) && !force {
			return fmt.Errorf(msg, f)
		}
	}

	for _, f := range files {
		if fileExists(f) && !force {
			return fmt.Errorf(msg, f)
		}
	}

	return nil
}

func vagrant_file(servers, agents int, box string) error {
	filevars := fmt.Sprintf("IMAGE_NAME = \"%s\"\nSERVER_NUMBER = %d\nAGENT_NUMBER = %d\n\n", box, servers, agents)

	return createFile("Vagrantfile", []byte(filevars+`Vagrant.configure("2") do |config|
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
`))
}

# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV["LC_ALL"] = "en_US.UTF-8"

$script = <<SCRIPT
# install Go
GO_VERSION=1.12
GO_HASH=d7d1f1f88ddfe55840712dc1747f37a790cbcaa448f6c9cf51bbe10aa65442f5
curl -sSL -o go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
echo "${GO_HASH} go.tar.gz" | sha256sum -c -
tar -C /usr/local -xvf go.tar.gz
rm go.tar.gz
ln -s /usr/local/go/bin/go* /usr/local/bin/

SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.box_check_update = false

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "512"
  end

  # forward ssh agent to easily ssh into the different machines
  config.ssh.forward_agent = true
  # always use Vagrants insecure key
  config.ssh.insert_key = false

  (1..2).each do |n|
    config.vm.define "node-#{n}" do |node|
      node.vm.hostname = "node-#{n}"

      node.vm.provision "shell", inline: $script
    end
  end
end

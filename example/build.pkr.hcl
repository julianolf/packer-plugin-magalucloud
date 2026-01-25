# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    magalucloud = {
      source  = "github.com/julianolf/magalucloud"
      version = ">= 0.0.1"
    }
  }
}

var "token" {
  type      = string
  default   = "${env("MGC_API_KEY")}"
  sensitive = true
}

var "source_image" {
  type    = string
  default = "cloud-ubuntu-22.04 LTS"
}

var "machine_type" {
  type    = string
  default = "BV1-1-10"
}

var "region" {
  type    = string
  default = "br-se1"
}

var "ssh_key" {
  type    = string
  default = "ssh-ed25519"
}

source "magalucloud" "example" {
  token        = var.token
  source_image = var.source_image
  machine_type = var.machine_type
  region       = var.region
  ssh_key      = var.ssh_key
}

build {
  name = "custom-ubuntu"

  sources = ["source.magalucloud.example"]

  provisioner "shell-local" {
    inline = ["echo Hello Packer"]
  }
}

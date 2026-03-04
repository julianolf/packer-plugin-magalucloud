# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    magalucloud = {
      source  = "github.com/julianolf/magalucloud"
      version = ">= 0.1.0"
    }
  }
}

variable "api_key" {
  type      = string
  default   = env("MGC_API_KEY")
  sensitive = true
}

variable "region" {
  type    = string
  default = "br-se1"
}

variable "source_image" {
  type    = string
  default = "cloud-debian-12 LTS"
}

variable "machine_type" {
  type    = string
  default = "BV1-1-10"
}

variable "name_prefix" {
  type    = string
  default = "packer-example"
}

variable "ssh_username" {
  type    = string
  default = "debian"
}

source "magalucloud" "example" {
  api_key      = var.api_key
  region       = var.region
  source_image = var.source_image
  machine_type = var.machine_type
  image_name   = "${var.name_prefix}-${formatdate("YYYYMMDDhhmmss", timestamp())}"
  ssh_username = var.ssh_username
}

build {
  name = "example"

  sources = ["source.magalucloud.example"]

  provisioner "shell" {
    inline = ["echo Hello Packer"]
  }
}

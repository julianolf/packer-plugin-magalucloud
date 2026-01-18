# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    magalucloud = {
      version = ">=v0.1.0"
      source  = "github.com/julianolf/magalucloud"
    }
  }
}

source "magalucloud-my-builder" "foo-example" {
  mock = local.foo
}

source "magalucloud-my-builder" "bar-example" {
  mock = local.bar
}

build {
  sources = [
    "source.magalucloud-my-builder.foo-example",
  ]

  source "source.magalucloud-my-builder.bar-example" {
    name = "bar"
  }

  provisioner "magalucloud-my-provisioner" {
    only = ["magalucloud-my-builder.foo-example"]
    mock = "foo: ${local.foo}"
  }

  provisioner "magalucloud-my-provisioner" {
    only = ["magalucloud-my-builder.bar"]
    mock = "bar: ${local.bar}"
  }

  post-processor "magalucloud-my-post-processor" {
    mock = "post-processor mock-config"
  }
}

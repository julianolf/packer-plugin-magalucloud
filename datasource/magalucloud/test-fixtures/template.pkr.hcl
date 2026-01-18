# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

data "magalucloud-my-datasource" "test" {
  mock = "mock-config"
}

locals {
  foo = data.magalucloud-my-datasource.test.foo
  bar = data.magalucloud-my-datasource.test.bar
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo foo: ${local.foo}",
      "echo bar: ${local.bar}",
    ]
  }
}

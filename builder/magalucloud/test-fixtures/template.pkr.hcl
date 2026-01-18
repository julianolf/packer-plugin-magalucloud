# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

source "magalucloud-my-builder" "basic-example" {
  mock = "mock-config"
}

build {
  sources = [
    "source.magalucloud-my-builder.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo build generated data: ${build.GeneratedMockData}",
    ]
  }
}

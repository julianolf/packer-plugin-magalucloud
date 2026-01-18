# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

source "scaffolding-my-builder" "basic-example" {
  mock = "mock-config"
}

build {
  sources = [
    "source.scaffolding-my-builder.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo build generated data: ${build.GeneratedMockData}",
    ]
  }
}

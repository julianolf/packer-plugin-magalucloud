# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "magalucloud-my-provisioner" {
    mock = "my-mock-config"
  }
}

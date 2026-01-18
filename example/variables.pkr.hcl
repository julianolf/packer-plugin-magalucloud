# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

locals {
  foo = data.magalucloud-my-datasource.mock-data.foo
  bar = data.magalucloud-my-datasource.mock-data.bar
}

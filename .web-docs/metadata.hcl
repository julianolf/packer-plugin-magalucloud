# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

integration {
  name        = "Magalu Cloud"
  description = "The magalucloud plugin can be used with HashiCorp Packer to create custom images on MGC."
  identifier  = "packer/hashicorp/magalucloud"
  flags       = ["hcp-ready"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/julianolf/packer-plugin-magalucloud/blob/main/LICENSE"
  }
  component {
    type = "builder"
    name = "Magalu Cloud"
    slug = "magalucloud"
  }
}

# Copyright (c) Juliano Fernandes 2026
# SPDX-License-Identifier: MPL-2.0

variable "api_key" {
  type      = string
  default   = env("MGC_API_KEY")
  sensitive = true
}

variable "access_key" {
  type      = string
  default   = env("MGC_ACCESS_KEY")
  sensitive = true
}

variable "secret_key" {
  type      = string
  default   = env("MGC_SECRET_KEY")
  sensitive = true
}

variable "region" {
  type    = string
  default = "br-se1"
}

variable "bucket" {
  type    = string
  default = env("BUCKET")
}

variable "name_prefix" {
  type    = string
  default = "packer-testacc"
}

source "file" "test" {
  source = "test-fixtures/image.qcow2"
  target = "test-fixtures/test.qcow2"
}

build {
  name = "test"

  sources = ["source.file.test"]

  post-processor "magalucloud-import" {
    api_key    = var.api_key
    access_key = var.access_key
    secret_key = var.secret_key
    region     = var.region
    bucket     = var.bucket
    image_name = "${var.name_prefix}-${formatdate("YYYYMMDDhhmmss", timestamp())}"
  }
}

Type: `magalucloud-import`
Artifact BuilderId: `packer.post-processor.magalucloud-import`

The Magalu Cloud Image Import post-processor takes a Qcow2 (QEMU Copy On Write)
image and imports it as an MGC image available to the Magalu Cloud Virtual Machines service.

~> This post-processor is for advanced users. Please ensure you read the
[Importing Custom Images](https://docs.magalu.cloud/docs/computing/images/import-custom-images/overview)
documentation before using it.

## How Does it Work?

The import process works by uploading a temporary copy of the image
to a bucket in MGC Object Storage and calling an import task on the Images service.
Once completed, a custom image is created. The temporary image copy
in the MGC bucket can be discarded after the import is complete.

~> **Note**: To prevent Packer from deleting the image copy in the MGC bucket, set the `skip_cleanup` configuration option to `true`.

**Required**

* `api_key` (string) - The API key used for authorization with Magalu Cloud.
* `access_key` (string) - The key pair ID used for authorization with Magalu Cloud Object Storage.
* `secret_key` (string) - The key pair secret used for authorization with Magalu Cloud Object Storage.
* `region` (string) - The region where the image will be built.
* `bucket` (string) - The name of the MGC bucket where the image will be uploaded.
* `image_name` (string) - The name for the new image.

**Optional**

* `url` (string) - The API URL used to interact with Magalu Cloud services (defaults to the selected region URL).
* `endpoint` (string) - The API URL used to interact with MGC Object Storage (defaults to the selected region URL).
* `filename` (string) - The name of the temporary copy of the image in MGC Object Storage (defaults to the artifact filename).
* `platform` (string) - The system platform (defaults to `linux`).
* `architecture` (string) - The architecture or processor type that this image supports (defaults to `x86/64`).
* `license` (string) - The license associated with the software in the image (defaults to `unlicensed`).
* `version` (string) - The version of the image (empty by default).
* `description` (string) - The description of the image (empty by default).
* `uefi` (bool) - Indicates whether the operating system in the image supports UEFI as a bootloader (defaults to `false`).
* `expires` (string) - The amount of time the signed URL of the temporary image in MGC Object Storage is valid (defaults to one hour, `"1h"`). Valid time units are `"ns"`, `"us"` (or `"µs"`), `"ms"`, `"s"`, `"m"`, and `"h"`.
* `skip_cleanup` (bool) - Skip the cleanup stage and keep the uploaded image in the MGC bucket (defaults to `false`).

### Example Usage

Here is a basic example. This assumes that a Qcow2 image already
exists in the template root directory and that the MGC bucket
has been created.

```hcl id="zq0v9k"
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

source "file" "ubuntu" {
  source = "${path.root}/jammy-server-cloudimg-amd64-disk-kvm.img"
  target = "ubuntu-22.04-lts-amd64.qcow2"
}

build {
  name = "custom-ubuntu"

  sources = ["source.file.ubuntu"]

  post-processor "magalucloud-import" {
    api_key    = var.api_key
    access_key = var.access_key
    secret_key = var.secret_key
    region     = "br-se1"
    bucket     = "custom-images"
    image_name = "ubuntu-22_04"
  }
}
```

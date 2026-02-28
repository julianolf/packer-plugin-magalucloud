Type: `magalucloud`
Artifact BuilderId: `julianolf.magalucloud`

The `magalucloud` Packer builder is able to create [images](https://docs.magalu.cloud/docs/computing/images/general-overview)
for use with Magalu Cloud [Virtual Machines](https://docs.magalu.cloud/docs/computing/virtual-machine/overview)
based on existing images.

It is possible to build images from scratch, but not with the `magalucloud` Packer builder.
This process is recommended only for advanced users. Please see [Importing Custom Images](https://docs.magalu.cloud/docs/computing/images/import-custom-images/overview)
for more information.

**Required**

* `api_key` (string) - The API key used for authorization with Magalu Cloud.
* `region` (string) - The region where the image will be built.
* `source_image` (string) - The ID or name of the base image to use.
* `machine_type` (string) - The ID or name of the machine type to use to build the image.

**Optional**

* `availability_zone` (string) - The availability zone where the image will be built (defaults to `{region}-a`).
* `image_name` (string) - The name for the new image (defaults to `packer-{UUIDv4}`).
* `url` (string) - The API URL used to interact with Magalu Cloud (defaults to the selected region URL).

### Example Usage

```hcl
variable "api_key" {
  type      = string
  default   = env("MGC_API_KEY")
  sensitive = true
}

source "magalucloud" "example" {
  api_key      = var.api_key
  region       = "br-se1"
  source_image = "cloud-ubuntu-22.04 LTS"
  machine_type = "BV1-1-10"
}

build {
  name = "custom-ubuntu"

  sources = ["source.magalucloud.example"]
}
```

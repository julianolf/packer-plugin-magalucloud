Type: `magalucloud`
Artifact BuilderId: `julianolf.magalucloud`

The `magalucloud` Packer builder is able to create
[images](https://docs.magalu.cloud/docs/computing/images/general-overview) for use with
Magalu Cloud [Virtual Machines](https://docs.magalu.cloud/docs/computing/virtual-machine/overview)
based on existing images.

It is possible to build images from scratch, but not with the `magalucloud`
Packer builder. The process is recommended only for advanced users, please see
[Importing Custom Images](https://docs.magalu.cloud/docs/computing/images/import-custom-images/overview)
for more information.

**Required**

- `token` (string) - The API token to use for authorization on Magalu Cloud.
- `region` (string) - The region where the image will be built.
- `source_image` (string) - The ID or name of the base image to use.
- `machine_type` (string) - The ID or name of the machine type to use to build the image.
- `ssh_key` (string) - The name of the SSH key to use to connect to the build instance.

**Optional**

- `image_name` (string) - The name for the new image (default to packet-{ISO-date}).
- `url` (string) - The API URL to interact with Magalu Cloud (default to the selected region URL).

### Example Usage


```hcl

source "magalucloud" "example" {
  token        = "${env("MGC_API_KEY")}"
  source_image = "cloud-ubuntu-22.04 LTS"
  machine_type = "BV1-1-10"
  region       = "br-se1"
  ssh_key      = "ssh-ed25519"
}

build {
  name = "custom-ubuntu"

  sources = ["source.magalucloud.example"]
}
```

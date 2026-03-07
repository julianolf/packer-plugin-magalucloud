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

#### Basic Example

Here is a basic example. This assumes that a Qcow2 image already
exists in the template root directory and that the MGC bucket
has been created.

```hcl
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

source "file" "example" {
  source = "debian.qcow2"
  target = "image.qcow2"
}

build {
  name = "example"

  sources = ["source.file.example"]

  post-processor "magalucloud-import" {
    api_key    = var.api_key
    access_key = var.access_key
    secret_key = var.secret_key
    region     = "br-se1"
    bucket     = "custom-images"
    image_name = "custom-debian"
  }
}
```

#### QEMU Builder Example

Here is a complete example for building an Alpine Linux image with nginx pre-installed and importing as a custom image on Magalu Cloud.

```hcl
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

variables {
  region      = "br-se1"
  bucket      = "custom-images"
  private_key = "~/.ssh/id_ed25519"
  public_key  = "~/.ssh/id_ed25519.pub"
}

locals {
  meta_data = {
    instance-id    = "alpine-01"
    local-hostname = "alpine"
  }
  user_data = {
    ssh_authorized_keys = [
      file(var.public_key)
    ]
  }
}

source "qemu" "example" {
  iso_url              = "https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/cloud/generic_alpine-3.23.3-x86_64-bios-cloudinit-r0.qcow2"
  iso_checksum         = "md5:7dab3974c59468394b5763b040c6e209"
  disk_image           = true
  format               = "qcow2"
  disk_interface       = "virtio"
  net_device           = "virtio-net"
  cpus                 = 1
  memory               = 1024
  headless             = true
  ssh_username         = "alpine"
  ssh_private_key_file = var.private_key
  ssh_timeout          = "10m"
  vm_name              = "alpine-nginx.qcow2"
  output_directory     = "output"
  shutdown_command     = "doas poweroff"
  cd_label             = "cidata"
  cd_content = {
    "meta-data" = yamlencode(local.meta_data)
    "user-data" = "#cloud-config\n${yamlencode(local.user_data)}"
  }
}

build {
  sources = ["source.qemu.example"]

  provisioner "shell" {
    inline = [
      "doas apk update",
      "doas apk add nginx",
      "doas rc-update add nginx default",
      "doas rc-service nginx start"
    ]
  }

  provisioner "shell" {
    inline = [
      "doas cloud-init clean --logs --machine-id --seed",
      "doas rm -rf /tmp/*",
      "doas rm -rf /var/cache/apk/*",
      "doas rm -f /etc/ssh/ssh_host_*",
      "doas rm -f /root/.ssh/authorized_keys /home/alpine/.ssh/authorized_keys",
      "doas rm -f /root/.bash_history /home/alpine/.bash_history"
    ]
  }

  post-processor "magalucloud-import" {
    api_key     = var.api_key
    access_key  = var.access_key
    secret_key  = var.secret_key
    region      = var.region
    bucket      = var.bucket
    image_name  = "alpine-nginx"
    version     = "0.1.0"
    description = "Alpine Linux with nginx pre-installed"
  }
}
```

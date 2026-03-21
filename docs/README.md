### Installation

To install this plugin, copy and paste this code into your Packer configuration,
then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    name = {
      source  = "github.com/julianolf/magalucloud"
      version = ">= 0.2.0"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage the installation of this plugin.

```sh
$ packer plugins install github.com/julianolf/magalucloud
```

### Components

#### Builders

* [magalucloud](/packer/integrations/julianolf/magalucloud/latest/components/builder/magalucloud) - The
  magalucloud builder creates images from existing ones by launching an instance, provisioning it, and then
  exporting it as a reusable image.

**Note**: Magalu Cloud currently does not support creating custom images directly from virtual machine instances. As a workaround, the plugin creates an instance snapshot, which can be used to create new instances.

#### Post-Processors

* [magalucloud-import](/packer/integrations/julianolf/magalucloud/latest/components/post-processor/magalucloud-import) -
  The magalucloud-import post-processor imports an existing image as an MGC custom image that can be
  used to launch instances.

### Authentication

Authenticating with Magalu Cloud services requires an API key, a key pair ID, and a key pair secret.

The `magalucloud` builder requires only the API key, while `magalucloud-import` requires all three credentials,
since it uploads the existing image to the object storage service.

To better understand each key, its purpose, and how to generate them, read the
[documentation](https://docs.magalu.cloud/docs/devops-tools/api-keys/overview).

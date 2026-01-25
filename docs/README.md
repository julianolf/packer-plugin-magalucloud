### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    name = {
      source  = "github.com/julianolf/magalucloud"
      version = ">=0.0.1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/julianolf/magalucloud
```

### Components

#### Builders

- [builder](/packer/integrations/julianolf/magalucloud/latest/components/builder/magalucloud) - The
  magalucloud builder creates images from existing ones, by launching an instance, provisioning it,
  then exporting it as a reusable image.

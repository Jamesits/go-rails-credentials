# Rails Credentials

Golang library, CLI tool and OpenTofu/Terraform provider for Ruby on Rails credentials files encryption/decryption.

![Works - On My Machine](https://img.shields.io/badge/Works-On_My_Machine-2ea44f) ![Project Status - Feature Complete](https://img.shields.io/badge/Project_Status-Feature_Complete-2ea44f)

## Usage

### Library

GoDoc: [![Go Reference](https://pkg.go.dev/badge/github.com/jamesits/go-rails-credentials/pkg/credentials.svg)](https://pkg.go.dev/github.com/jamesits/go-rails-credentials/pkg/credentials)

See [edit.go](cmd/rails-credentials/edit.go) for a complete example.

### CLI

- `rails-credentials show` as a drop-in replacement for `rails credentials:show`
- `rails-credentials edit` as a drop-in replacement for `rails credentials:edit`

Environment variables:

- `RAILS_ENV` and `RAILS_MASTER_KEY` work as intended
- `VISUAL` or `EDITOR` specifies the editor to use for `edit` command

Command line arguments:

- Run under the root directory of your Rails project or set `--base-dir <dir>` to your project directory
- If your files are not at the default location, use `--master-key-file <path>` and `--credentials-file <path>` to set the paths explicitly; `config.credentials.{content,key}_path` does not work
- See the embedded help (`rails-credentials --help`) for detailed usage

Notes:

- Rails refuse to work if `master.key` has a newline at the end; our parser is more relax on this issue
- `rails credentials:diff` is not planned for now; contributions are welcomed

### OpenTofu / Terraform Provider

Decrypt the credentials on the fly (can also be used as a credentials validator):

```hcl
data "railscred_file" "example" {
  master_key        = file("${path.module}/config/master.key")
  encrypted_content = file("${path.module}/config/credentials.yml.enc")
}

output "credentials" {
  value     = data.railscred_file.example.content
  sensitive = true
}
```

Manage the plaintext credentials inside the Tofu config:

```hcl
# generate a random master key
resource "railscred_master_key" "example" {}

# plaintext credentials
data "railscred_inline" "example" {
  master_key = railscred_master_key.example.master_key
  content    = <<-EOT
# smtp:
#   user_name: my-smtp-user
#   password: my-smtp-password
#
# aws:
#   access_key_id: 123
#   secret_access_key: 345

# Used as the base secret for all MessageVerifiers in Rails, including the one protecting cookies.
secret_key_base:
EOT
}

output "encrypted_credentials" {
  value     = data.railscred_file.example.encrypted_content
}

# Example of using them in Kubernetes

resource "kubernetes_secret_v1" "rails" {
  metadata {
    name      = "rails"
    namespace = "application"
  }
  data = {
    "RAILS_MASTER_KEY" = railscred_master_key.example.master_key
  }
}

resource "kubernetes_secret_v1" "rails_credentials" {
  metadata {
    name      = "rails-credentials"
    namespace = "application"
  }
  data = {
    "credentials.yml.enc" = data.railscred_inline.example.content
  }
}

resource "kubernetes_deployment_v1" "rails" {
  metadata {
    name      = "rails"
    namespace = "application"
  }
  spec {
    template {
      spec {
        volume {
          name = "rails-credentials"
          secret {
            secret_name = "rails-credentials"
            items {
              key  = "credentials.yml.enc"
              path = "credentials.yml.enc"
            }
          }
        }
        container {
          env_from {
            secret_ref {
              name = "rails"
            }
          }
          volume_mount {
            name       = "rails-credentials"
            mount_path = "/app/config/credentials.yml.enc"
            sub_path   = "credentials.yml.enc"
            read_only  = true
          }
        }
      }
    }
  }
}
```

## Development

### Building

```shell
goreleaser build --snapshot --clean
```

### Tofu Provider Testing

To use the Tofu provider locally:

```shell
cat > .terraformrc <<EOF
provider_installation {
    dev_overrides {
        "jamesits/railscred" = "./dist/provider_linux_amd64_v1"
    }
    direct {}
}
EOF

export TF_CLI_CONFIG_FILE="./.terraformrc"
```

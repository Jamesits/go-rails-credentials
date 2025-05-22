# Rails Credentials CLI

Golang library, standalone CLI tool and Tofu provider for Rails credentials files operations.

![Works - On My Machine](https://img.shields.io/badge/Works-On_My_Machine-2ea44f) ![Project Status - Premature](https://img.shields.io/badge/Project_Status-Premature-yellow)

## Usage

### CLI

- `rails-credentails show` as a drop-in replacement for `rails credentials:show`
- `rails-credentials edit` as a drop-in replacement for `rails credentials:edit`

Notes:

- Run under the root directory of your Rails project; alternatively, use `--base-dir <dir>` to set your project directory
- `RAILS_ENV` and `RAILS_MASTER_KEY` will work as intended; use `--master-key-file <path>` and `--credentials-file <path>` if the files are not at the default path
- See the embedded help for detailed usage

### OpenTofu / Terraform provider

TBD.

## Development

### Building

```shell
goreleaser build --snapshot --clean
```

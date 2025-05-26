resource "railscred_master_key" "example" {}

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

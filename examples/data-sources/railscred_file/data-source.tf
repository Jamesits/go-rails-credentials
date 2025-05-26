data "railscred_file" "example" {
  master_key        = file("${path.module}/config/master.key")
  encrypted_content = file("${path.module}/config/credentials.yml.enc")
}

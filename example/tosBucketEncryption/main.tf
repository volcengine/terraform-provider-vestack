resource "vestack_tos_bucket" "foo" {
  bucket_name   = "tf-acc-test-bucket1"
  public_acl    = "private"
  az_redundancy = "multi-az"
  project_name  = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "vestack_kms_keyring" "foo" {
  keyring_name = "acc-test-keyring"
  description  = "acc-test"
  project_name = "default"
}

resource "vestack_tos_bucket_encryption" "foo" {
  bucket_name = vestack_tos_bucket.foo.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm       = "kms"
      kms_data_encryption = "AES256"
      kms_master_key_id   = vestack_kms_keyring.foo.id
    }
  }
}

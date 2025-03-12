# # https://www.linode.com/docs/guides/how-to-build-your-infrastructure-using-terraform-and-linode/#use-linode-object-storage-to-store-state
# terraform {
#   backend "s3" {
#     endpoints = {
#       s3 = "https://jp-osa-1.linodeobjects.com"
#     }
#     profile                     = "linode-s3"
#     skip_region_validation      = true
#     skip_credentials_validation = true
#     skip_requesting_account_id  = true
#     skip_s3_checksum            = true
#     bucket                      = "terraform-backend"
#     key                         = "tf-state.json"
#     region                      = "jp-osa-1"
#   }
# }

provider "aws" {
  region = local.region
}

data "aws_availability_zones" "available" {}

data "aws_caller_identity" "self" {}

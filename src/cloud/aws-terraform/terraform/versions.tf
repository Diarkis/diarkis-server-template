terraform {
  required_version = ">= 1.3.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.61"
    }
  }
  backend "s3" {
    region = "ap-northeast-1"
  }
}

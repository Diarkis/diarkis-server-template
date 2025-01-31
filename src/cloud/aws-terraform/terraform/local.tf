locals {
  name   = "diarkis"
  region = "ap-northeast-1"

  vpc_cidr = "10.1.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
  instance_types = ["m6i.large"]
  _env = {
    dev = {
      prefix = "dev"
      instance_types = ["t3.medium"]
      need_operation_tools = true
    },
    stg = {
      prefix = "stg"
      instance_types = ["m6i.large"]
      need_operation_tools = false
    },
    mnt = {
      prefix = "mnt"
      instance_types = ["m6i.large"]
      need_operation_tools = false
    },
    prd = {
      prefix = "prd"
      instance_types = ["m6i.large"]
      need_operation_tools = true
    },
  }
  env = local._env[var.env]
}

module "prometheus" {
  source = "terraform-aws-modules/managed-service-prometheus/aws"
  version = "3.0.0"
  workspace_alias = "${local.env.prefix}-${local.name}"
  count = local.env.need_operation_tools ? 1 : 0
}

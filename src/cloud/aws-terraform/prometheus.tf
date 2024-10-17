module "prometheus" {
  source = "terraform-aws-modules/managed-service-prometheus/aws"
  version = "3.0.0"

  workspace_alias = "${local.name}"
}

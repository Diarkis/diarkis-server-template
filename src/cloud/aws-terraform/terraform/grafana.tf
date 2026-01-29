resource "aws_grafana_workspace" "diarkis" {
  name = "${local.env.prefix}-${local.name}"
  count = local.env.need_operation_tools ? 1 : 0
  account_access_type      = "CURRENT_ACCOUNT"
  authentication_providers = ["SAML", "AWS_SSO"]
  permission_type          = "SERVICE_MANAGED"
  role_arn                 = aws_iam_role.grafana-assume[count.index].arn
}

resource "aws_iam_role" "grafana-assume" {
  name = "grafana-assume"
  count = local.env.need_operation_tools ? 1 : 0
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "grafana.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "grafana-policy" {
  name   = "grafana_policy"
  count = local.env.need_operation_tools ? 1 : 0
  role   = aws_iam_role.grafana-assume[count.index].id
  policy =  jsonencode({
    Version = "2012-10-17"
    Statement = [
    		{
			"Action": [
				"aps:ListWorkspaces",
				"aps:DescribeWorkspace",
				"aps:QueryMetrics",
				"aps:GetLabels",
				"aps:GetSeries",
				"aps:GetMetricMetadata"
			],
			"Effect": "Allow",
			"Resource": "*"
		},
    ]
  })
}

module "eks_al2" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"
  cluster_name    = "${local.env.prefix}-${local.name}"
  cluster_version = "1.31"

  cluster_endpoint_public_access  = true

  # EKS Addons
  cluster_addons = {
    coredns                = {}
    eks-pod-identity-agent = {}
    kube-proxy             = {}
    vpc-cni                = {}
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  enable_cluster_creator_admin_permissions = true

  create_node_security_group = true

  node_security_group_additional_rules = {
    ingress_node_communications = {
      description = "Ingress Node to Node"
      protocol    = "-1"
      from_port   = 0
      to_port     = 65535
      type        = "ingress"
      self        = true
      #source_cluster_security_group = true
    }
    egress_node_communications = {
      description = "Egress Node to Node"
      protocol    = "-1"
      from_port   = 0
      to_port     = 65535
      type        = "egress"
      self        = true
      #source_cluster_security_group = true
    }
    ingress_diarkis_udp = {
      description = "allow diarkis ports from all"
      protocol    = "udp"
      from_port   = 7000
      to_port     = 8000
      type        = "ingress"
      cidr_blocks = ["0.0.0.0/0"]
      ipv6_cidr_blocks = ["::/0"]
    }
    ingress_diarkis_tcp = {
      description = "allow diarkis ports from all"
      protocol    = "tcp"
      from_port   = 7000
      to_port     = 8000
      type        = "ingress"
      cidr_blocks = ["0.0.0.0/0"]
      ipv6_cidr_blocks = ["::/0"]
    }
  }
  eks_managed_node_group_defaults = {
    iam_role_policy_statements = [{
        sid = "Autoscaling"
        actions = [
                "autoscaling:DescribeAutoScalingGroups",
                "autoscaling:DescribeAutoScalingInstances",
                "autoscaling:DescribeLaunchConfigurations",
                "autoscaling:DescribeScalingActivities",
                "autoscaling:DescribeTags",
                "ec2:DescribeInstanceTypes",
                "ec2:DescribeLaunchTemplateVersions",
                "autoscaling:SetDesiredCapacity",
                "autoscaling:TerminateInstanceInAutoScalingGroup",
                "ec2:DescribeImages",
                "ec2:GetInstanceTypesFromInstanceRequirements",
                "eks:DescribeNodegroup"
                ]
        resources = ["*"]
        effect = "Allow"
    }]
  }

  eks_managed_node_groups = {
    diarkis-private = {
      ami_type       = "AL2_x86_64"
      instance_types = local.env.instance_types
      subnet_ids = module.vpc.private_subnets
      min_size = 1
      max_size = 10
      desired_size = 2
      labels = {
         "diarkis.io/network" = "private"
      }
      metadata_options = {
        "http_endpoint": "enabled",
        "http_put_response_hop_limit": 2,
        "http_tokens": "optional"
      }
    }
    diarkis-public = {
      ami_type       = "AL2_x86_64"
      instance_types = local.env.instance_types
      subnet_ids = module.vpc.public_subnets

      min_size = 1
      max_size = 10
      desired_size = 2
      labels = {
         "diarkis.io/network" = "public"
      }
      metadata_options = {
        "http_endpoint": "enabled",
        "http_put_response_hop_limit": 2,
        "http_tokens": "optional"
      }

    }
  }
}

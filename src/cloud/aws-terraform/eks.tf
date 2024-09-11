module "eks_al2" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = "${local.name}"
  cluster_version = "1.30"

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
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
    ingress_diarkis = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 7000
      to_port     = 8000
      type        = "ingress"
      cidr_blocks = ["0.0.0.0/0"]
      ipv6_cidr_blocks = ["::/0"]
    }
  }
  eks_managed_node_groups = {
    diarkis-private = {
      ami_type       = "AL2_x86_64"
      instance_types = local.instance_types
      subnet_ids = module.vpc.private_subnets

      min_size = 1
      max_size = 10
      desired_size = 2
    }
    diarkis-public = {
      ami_type       = "AL2_x86_64"
      instance_types = ["m6i.large"]
      subnet_ids = module.vpc.public_subnets

      min_size = 1
      max_size = 10
      desired_size = 2
    }
  }
}

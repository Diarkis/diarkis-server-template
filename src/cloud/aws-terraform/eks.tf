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

  eks_managed_node_groups = {
    diarkis-private = {
      ami_type       = "AL2_x86_64"
      instance_types = ["m6i.large"]
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

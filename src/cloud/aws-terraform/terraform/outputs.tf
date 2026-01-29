output "eks-connect-command" {
  description = "cli command which generates kubeconfig"
  value       = "aws eks update-kubeconfig --name ${local.env.prefix}-${local.name} --region ${local.region}"
}

output "auth-docker-command" {
  description = "authenticate to ecr command"
  value       = "aws ecr get-login-password --region ${local.region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.self.account_id}.dkr.ecr.ap-northeast-1.amazonaws.com"
}

output "http-docker-push-command" {
  description = "cli command which push http server image"
  value       = "docker push ${data.aws_caller_identity.self.account_id}.dkr.ecr.ap-northeast-1.amazonaws.com/diarkis-http"
}

output "udp-docker-push-command" {
  description = "cli command which push udp server image"
  value       = "docker push ${data.aws_caller_identity.self.account_id}.dkr.ecr.ap-northeast-1.amazonaws.com/diarkis-udp"
}

output "mars-docker-push-command" {
  description = "cli command which push mars server image"
  value       = "docker push ${data.aws_caller_identity.self.account_id}.dkr.ecr.ap-northeast-1.amazonaws.com/diarkis-mars"
}

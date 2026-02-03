module "prometheus" {
  source          = "terraform-aws-modules/managed-service-prometheus/aws"
  version         = "3.0.0"
  workspace_alias = "${local.env.prefix}-${local.name}"
  count           = local.env.need_operation_tools ? 1 : 0
}

# Prometheus Scraper for EKS cluster monitoring
resource "aws_prometheus_scraper" "eks_scraper" {
  count = local.env.need_operation_tools ? 1 : 0

  scrape_configuration = <<EOT
global:
  scrape_interval: 30s

scrape_configs:
  - job_name: 'kubernetes-apiservers'
    kubernetes_sd_configs:
      - role: endpoints
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: false
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
        action: keep
        regex: default;kubernetes;https

  - job_name: 'kubernetes-nodes'
    kubernetes_sd_configs:
      - role: node
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: false
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)

  - job_name: 'kubernetes-cadvisor'
    kubernetes_sd_configs:
      - role: node
    scheme: https
    metrics_path: /metrics/cadvisor
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: false
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)

  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: kubernetes_pod_name
  - job_name: diarkis
    metrics_path: /metrics/prometheus/v/3
    static_configs:
      - targets:
          - 172.20.85.238
EOT

  source {
    eks {
      cluster_arn = module.eks_al2.cluster_arn
      subnet_ids  = module.vpc.private_subnets
    }
  }

  destination {
    amp {
      workspace_arn = module.prometheus[0].workspace_arn
    }
  }

  tags = {
    Name        = "${local.env.prefix}-${local.name}-prometheus-scraper"
    Environment = local.env.prefix
  }

  depends_on = [
    module.prometheus,
    module.eks_al2
  ]
}

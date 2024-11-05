output "kubeconfig" {
	value = linode_lke_cluster.lke-1.kubeconfig
	sensitive = true
}

output "api_endpoints" {
	value = linode_lke_cluster.lke-1.api_endpoints
}

output "status" {
   value = linode_lke_cluster.lke-1.status
}

output "id" {
   value = linode_lke_cluster.lke-1.id
}

output "pool" {
   value = linode_lke_cluster.lke-1.pool
}

output "firewall_id" {
  value = linode_firewall.diarkis-firewall.id
}

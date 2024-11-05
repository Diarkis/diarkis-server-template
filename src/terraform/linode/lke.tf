resource "linode_lke_cluster" "lke-1" {
    k8s_version = var.lke-1.k8s_version
    label = var.lke-1.label
    region = var.region
    tags = var.lke-1.tags

    dynamic "pool" {
        for_each = var.pools
        content {
            type  = pool.value["type"]
            count = pool.value["count"]
        }
    }
}


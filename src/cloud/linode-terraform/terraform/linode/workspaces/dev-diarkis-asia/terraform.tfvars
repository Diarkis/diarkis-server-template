region = "jp-osa"

lke-1 = {
  k8s_version = "1.31"
  label       = "dev-diarkis-asia"
  tags        = ["diarkis"]
}

pools = [
  {
    type = "g6-standard-1" // Linode 2 GB (1 CPU, 2 GB RAM, 50 GB Storage)
    # type = "g6-standard-4" // Linode 8 GB (4 CPU, 8 GB RAM, 160 GB Storage)
    # type = "g6-dedicated-4" // Dedicated 8 GB (4 CPU, 8 GB RAM, 160 GB Storage)
    count = 3
  }
]

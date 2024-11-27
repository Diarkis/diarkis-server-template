resource "linode_firewall" "diarkis-firewall" {
  label = "diarkis"

  inbound {
    label    = "allow-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-diarkis-tcp"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "7000-8000"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
  inbound {
    label    = "allow-diarkis-udp"
    action   = "ACCEPT"
    protocol = "UDP"
    ports    = "7000-8000"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-kubelet"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "10250"
    ipv4     = ["192.168.128.0/17"]
  }

  inbound {
    label    = "allow-nodeport"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "30000-32767"
    ipv4     = ["192.168.128.0/17"]
  }

  inbound_policy = "DROP"

  outbound_policy = "ACCEPT"

}

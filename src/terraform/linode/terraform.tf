terraform {
	required_providers {
		linode = {
			source	= "linode/linode"
		}
	}
}

provider "linode" {
	token = "${var.token}"
}

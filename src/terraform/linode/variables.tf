variable "token" {}

variable "region" {
	default = ""
}

variable "lke-1" {
  type = object({
		k8s_version = string
		label = string
    tags = list(string)
  })
	default = {
		k8s_version = ""
		label = "diarkis-cluster-terraform"
		tags = ["diarkis"]
	}
}

variable "pools" {
  type = list(object({
    type = string
    count = number
  }))
  default = [
    {
      type = "g6-standard-1"
      count = 1
    }
  ]
}

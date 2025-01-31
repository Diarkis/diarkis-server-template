resource "aws_ecr_repository" "http" {
  name                 = "diarkis-http"
  image_tag_mutability = "MUTABLE"
  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "udp" {
  name                 = "diarkis-udp"
  image_tag_mutability = "MUTABLE"
  force_delete = true
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "tcp" {
  name                 = "diarkis-tcp"
  image_tag_mutability = "MUTABLE"
  force_delete = true
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "mars" {
  name                 = "diarkis-mars"
  image_tag_mutability = "MUTABLE"
  force_delete = true
  image_scanning_configuration {
    scan_on_push = true
  }
}

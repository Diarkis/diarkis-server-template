resource "aws_ecr_repository" "diarkis-http" {
  name                 = "diarkis-http"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "diarkis-udp" {
  name                 = "diarkis-udp"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "diarkis-mars" {
  name                 = "diarkis-mars"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

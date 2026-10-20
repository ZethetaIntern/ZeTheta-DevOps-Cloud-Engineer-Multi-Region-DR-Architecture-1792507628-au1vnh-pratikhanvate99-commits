module "network_primary" {
  source       = "./modules/network"
  name         = "${var.name}-mumbai"
  region       = var.primary_region
  cidr         = var.vpc_cidr_primary
  subnet_cidrs = var.subnet_cidrs_primary
}

module "network_dr" {
  source = "./modules/network"
  providers = {
    aws = aws.dr
  }
  name         = "${var.name}-hyderabad"
  region       = var.dr_region
  cidr         = var.vpc_cidr_dr
  subnet_cidrs = var.subnet_cidrs_dr
}

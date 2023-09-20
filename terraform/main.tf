provider "aws" {
  region = var.region
}

module "vpc" {
  source = "./vpc"
   cidr_block = "10.0.0.0/16"
}

module "iam" {
  source = "./iam"
}

module "eks" {
  source     = "./eks"
  cluster_name = var.cluster_name
   vpc_id = module.vpc.vpc_id
   subnet_ids = module.vpc.subnet_ids
}


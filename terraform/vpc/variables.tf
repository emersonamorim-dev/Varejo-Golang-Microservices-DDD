variable "cidr_block" {
  description = "CIDR for VPC"
  default     = "10.0.0.0/16"
}

variable "public_subnets" {
  description = "CIDR blocks for public subnets"
  default = [
    "10.0.1.0/24",
    "10.0.2.0/24"
  ]
}

variable "private_subnets" {
  description = "CIDR blocks for private subnets"
  default = [
    "10.0.3.0/24",
    "10.0.4.0/24"
  ]
}

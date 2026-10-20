variable "primary_region" {
  type    = string
  default = "ap-south-1"
}

variable "dr_region" {
  type    = string
  default = "ap-south-2"
}

variable "name" {
  type    = string
  default = "paysecure"
}

variable "vpc_cidr_primary" {
  type    = string
  default = "10.40.0.0/16"
}

variable "vpc_cidr_dr" {
  type    = string
  default = "10.41.0.0/16"
}

variable "subnet_cidrs_primary" {
  type    = list(string)
  default = ["10.40.0.0/20", "10.40.16.0/20", "10.40.32.0/20"]
}

variable "subnet_cidrs_dr" {
  type    = list(string)
  default = ["10.41.0.0/20", "10.41.16.0/20", "10.41.32.0/20"]
}

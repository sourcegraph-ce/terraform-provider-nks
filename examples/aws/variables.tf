# AWS configuration variables
variable "aws_code" {
  description = "AWS cluster type code"
  default     = "aws"
}

variable "aws_k8s_version" {
  description = "AWS kubernetes version"
  default     = "v1.11.1"
}

variable "aws_platform" {
  description = "AWS platform type"
  default     = "coreos"
}

variable "aws_region" {
  description = "AWS region"
  default     = "us-east-2"
}

variable "aws_zone" {
  description = "AWS zone"
  default     = "us-east-2a"
}

variable "aws_network_id" {
  description = "AWS network ID"
  default     = "__new__"
}

variable "aws_network_cidr" {
  description = "AWS network CIDR"
  default     = "10.0.0.0/16"
}

variable "aws_subnet_id" {
  description = "AWS subnet ID"
  default     = "__new__"
}

variable "aws_subnet_cidr" {
  description = "AWS subnet CIDR"
  default     = "10.0.0.0/24"
}

variable "aws_zone2" {
  description = "AWS zone for second master"
  default     = "us-east-2b"
}

variable "aws_subnet_id2" {
  description = "AWS subnet ID for second master"
  default     = "__new__"
}

variable "aws_subnet_cidr2" {
  description = "AWS subnet CIDR for second master"
  default     = "10.0.1.0/24"
}

variable "aws_master_size" {
  description = "AWS master node size"
  default     = "t2.medium"
}

variable "aws_worker_size" {
  description = "AWS worker node size"
  default     = "t2.medium"
}

# AWS configuration variables
variable "aws_code" {
  description = "AWS cluster type code"
  default     = "aws"
}

variable "aws_keyset" {
  description = "AWS keyset ID"
  default     = 3625
}

variable "aws_k8s_version" {
  description = "AWS kubernetes version"
  default     = "v1.8.3"
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
  default     = "vpc-14ca497c"
}

variable "aws_network_cidr" {
  description = "AWS network CIDR"
  default     = "172.31.0.0/16"
}

variable "aws_subnet_id" {
  description = "AWS subnet ID"
  default     = "subnet-f4295c9c"
}

variable "aws_subnet_cidr" {
  description = "AWS subnet CIDR"
  default     = "172.31.0.0/24"
}

variable "aws_master_size" {
  description = "AWS master node size"
  default     = "t2.medium"
}

variable "aws_worker_size" {
  description = "AWS worker node size"
  default     = "t2.medium"
}

# Azure configuration variables
variable "azure_code" {
  description = "Azure cluster type code"
  default     = "azure"
}

variable "azure_keyset" {
  description = "Azure keyset ID"
  default     = 1671
}

variable "azure_k8s_version" {
  description = "Azure kubernetes version"
  default     = "v1.8.7"
}

variable "azure_platform" {
  description = "Azure platform type"
  default     = "coreos"
}

variable "azure_region" {
  description = "Azure region"
  default     = "eastus"
}

variable "azure_resource_group" {
  description = "Azure resource group"
  default     = "__new__"
}

variable "azure_network_id" {
  description = "Azure network ID"
  default     = "__new__"
}

variable "azure_network_cidr" {
  description = "Azure network CIDR"
  default     = "172.23.0.0/16"
}

variable "azure_subnet_id" {
  description = "Azure subnet ID"
  default     = "__new__"
}

variable "azure_subnet_cidr" {
  description = "Azure subnet CIDR"
  default     = "172.23.1.0/24"
}

variable "azure_master_size" {
  description = "Azure master node size"
  default     = "standard_f1"
}

variable "azure_worker_size" {
  description = "Azure worker node size"
  default     = "standard_f1"
}

# DigitalOcean configuration variables
variable "digitalocean_code" {
  description = "DigitalOcean cluster type code"
  default     = "do"
}

variable "digitalocean_keyset" {
  description = "DigitalOcean keyset ID"
  default     = 3556
}

variable "digitalocean_k8s_version" {
  description = "DigitalOcean kubernetes version"
  default     = "v1.8.3"
}

variable "digitalocean_platform" {
  description = "DigitalOcean platform type"
  default     = "coreos"
}

variable "digitalocean_region" {
  description = "DigitalOcean region"
  default     = "nyc1"
}

variable "digitalocean_master_size" {
  description = "DigitalOcean master node size"
  default     = "2gb"
}

variable "digitalocean_worker_size" {
  description = "DigitalOcean worker node size"
  default     = "2gb"
}

# GCE configuration variables
variable "gce_code" {
  description = "GCE cluster type code"
  default     = "gce"
}

variable "gce_keyset" {
  description = "GCE keyset ID"
  default     = 3553
}

variable "gce_k8s_version" {
  description = "GCE kubernetes version"
  default     = "v1.8.3"
}

variable "gce_platform" {
  description = "GCE platform type"
  default     = "coreos"
}

variable "gce_region" {
  description = "GCE region"
  default     = "us-west1-a"
}

variable "gce_master_size" {
  description = "GCE master node size"
  default     = "n1-standard-1"
}

variable "gce_worker_size" {
  description = "GCE worker node size"
  default     = "n1-standard-1"
}

# GKE configuration variables
variable "gke_code" {
  description = "GKE cluster type code"
  default     = "gke"
}

variable "gke_keyset" {
  description = "GKE keyset ID"
  default     = 1797
}

variable "gke_k8s_version" {
  description = "GKE kubernetes version"
  default     = "latest"
}

variable "gke_platform" {
  description = "GKE platform type"
  default     = "gci"
}

variable "gke_region" {
  description = "GKE region"
  default     = "us-west1-a"
}

variable "gke_master_size" {
  description = "GKE master node size"
  default     = "n1-standard-1"
}

variable "gke_worker_size" {
  description = "GKE worker node size"
  default     = "n1-standard-1"
}

# Packet configuration variables
variable "packet_code" {
  description = "Packet cluster type code"
  default     = "packet"
}

variable "packet_keyset" {
  description = "Packet keyset ID"
  default     = 3880
}

variable "packet_k8s_version" {
  description = "Packet kubernetes version"
  default     = "v1.8.7"
}

variable "packet_platform" {
  description = "Packet platform type"
  default     = "coreos"
}

variable "packet_region" {
  description = "Packet region"
  default     = "sjc1"
}

variable "packet_project_id" {
  description = "Packet project ID"
  default     = "93125c2a-8b78-4d4f-a3c4-7367d6b7cca8"
}

variable "packet_master_size" {
  description = "Packet master node size"
  default     = "baremetal_0"
}

variable "packet_worker_size" {
  description = "Packet worker node size"
  default     = "baremetal_0"
}

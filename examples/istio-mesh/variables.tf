# Organization
variable "organization_name" {
  description = "NKS organization name"
  default     = ""
}

# Cluster
variable "a_cluster_name" {
  description = "NKS cluster name"
  default     = ""
}

variable "b_cluster_name" {
  description = "NKS cluster name"
  default     = ""
}

# Keyset
variable "ssh_keyset_name" {
  description = "NKS ssh keyset name"
  default     = ""
}

variable "provider_keyset_name" {
  description = "Cloud provider keyset name"
  default     = ""
}

# Istio mesh
variable "istio_mesh_name" {
  description = "NKS istio mesh name"
  default     = ""
}

variable "istio_mesh_type" {
  description = "NKS istio mesh type"
  default     = ""
}

# Cloud provider configuration variables
variable "provider_code" {
  description = "Cloud provider type code"
  default     = ""
}

variable "provider_k8s_version" {
  description = "Cloud provider kubernetes version"
  default     = ""
}

variable "provider_etcd_type" {
  description = "Cloud provider etcd type"
  default     = ""
}

variable "provider_channel" {
  description = "Cloud provider channel"
  default     = ""
}

variable "provider_platform" {
  description = "Cloud provider platform type"
  default     = ""
}

variable "provider_region" {
  description = "Cloud provider region"
  default     = ""
}

variable "provider_resource_group" {
  description = "Cloud provider resource group"
  default     = ""
}

variable "provider_network_id" {
  description = "Cloud provider network ID"
  default     = ""
}

variable "provider_network_cidr" {
  description = "Cloud provider network CIDR"
  default     = ""
}

variable "provider_subnet_id" {
  description = "Cloud provider subnet ID"
  default     = ""
}

variable "provider_subnet_cidr" {
  description = "Cloud provider subnet CIDR"
  default     = ""
}

variable "provider_master_size" {
  description = "Cloud provider master node size"
  default     = ""
}

variable "provider_worker_size" {
  description = "Cloud provider worker node size"
  default     = ""
}

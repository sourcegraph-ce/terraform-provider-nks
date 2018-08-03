variable "azure_code" {
  description = "Azure cluster type code"
  default     = "azure"
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
  default     = "10.0.0.0/16"
}

variable "azure_subnet_id" {
  description = "Azure subnet ID"
  default     = "__new__"
}

variable "azure_subnet_cidr" {
  description = "Azure subnet CIDR"
  default     = "10.0.0.0/24"
}

variable "azure_master_size" {
  description = "Azure master node size"
  default     = "standard_f1"
}

variable "azure_worker_size" {
  description = "Azure worker node size"
  default     = "standard_f1"
}

variable "packet_code" {
  description = "Packet cluster type code"
  default     = "packet"
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

variable "packet_master_size" {
  description = "Packet master node size"
  default     = "baremetal_0"
}

variable "packet_worker_size" {
  description = "Packet worker node size"
  default     = "baremetal_0"
}

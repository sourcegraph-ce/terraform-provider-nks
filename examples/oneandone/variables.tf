variable "oneandone_code" {
  description = "OneAndOne cluster type code"
  default     = "oneandone"
}

variable "oneandone_k8s_version" {
  description = "OneAndOne kubernetes version"
  default     = "v1.8.7"
}

variable "oneandone_platform" {
  description = "OneAndOne platform type"
  default     = "coreos"
}

variable "oneandone_region" {
  description = "OneAndOne region"
  default     = "US"
}

variable "oneandone_master_size" {
  description = "OneAndOne master node size"
  default     = "m"
}

variable "oneandone_worker_size" {
  description = "OneAndOne worker node size"
  default     = "m"
}

variable "gce_code" {
  description = "GCE cluster type code"
  default     = "gce"
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

variable "gce_region2" {
  description = "GCE region"
  default     = "us-west1-b"
}

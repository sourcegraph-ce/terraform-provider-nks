variable "gke_code" {
  description = "GKE cluster type code"
  default     = "gke"
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

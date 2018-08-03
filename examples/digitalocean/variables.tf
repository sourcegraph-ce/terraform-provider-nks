variable "digitalocean_code" {
  description = "DigitalOcean cluster type code"
  default     = "do"
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

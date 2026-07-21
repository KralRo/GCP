variable "client_slug" {
  type        = string
  description = "Client slug, e.g. streckerova"
}

variable "project_id" {
  type        = string
  description = "GCP project ID for this client, e.g. rk-streckerova"
}

variable "region" {
  type        = string
  description = "GCP region for Cloud Run and Firestore"
  default     = "europe-west1"
}

variable "env" {
  type        = string
  description = "Environment name: prod today, dev once it's enabled for this client"
  default     = "prod"
}

variable "image" {
  type        = string
  description = "Container image to deploy to Cloud Run"
  default     = "us-docker.pkg.dev/cloudrun/container/hello"
}

variable "firestore_location" {
  type        = string
  description = "Firestore location"
  default     = "eur3"
}

variable "manage_dns" {
  type        = bool
  description = "Whether Cloud DNS manages this client's domain (false if the domain stays at the client's own registrar)"
  default     = false
}

variable "domain" {
  type        = string
  description = "Client's domain name"
  default     = ""
}

variable "dns_records" {
  type = list(object({
    name    = string
    type    = string
    ttl     = number
    rrdatas = list(string)
  }))
  description = "DNS records for the domain - published in Cloud DNS when manage_dns = true, otherwise just documents what to set up at the registrar"
  default     = []
}

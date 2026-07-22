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

variable "create_zone" {
  type        = bool
  description = "Whether to create a new Cloud DNS managed zone (true) or reference one that already exists, e.g. auto-provisioned by Cloud Domains at registration (false). Only relevant when manage_dns = true."
  default     = true
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

variable "domain_mappings" {
  type        = list(string)
  description = "Hostnames to map to the Cloud Run service, e.g. [\"streckerova.kralroman.org\", \"admin.streckerova.kralroman.org\"]"
  default     = []
}

variable "admin_host" {
  type        = string
  description = "Hostname the admin panel should also be served on at its root path (in addition to /admin on the main domain). Empty disables host-based admin routing."
  default     = ""
}

variable "ci_deployer_service_account" {
  type        = string
  description = "Email of the shared github-actions-deployer SA (lives in rk-infra-krl) that runs Terraform in CI - needs actAs on the per-client runtime SA to deploy Cloud Run"
  default     = "github-actions-deployer@rk-infra-krl.iam.gserviceaccount.com"
}

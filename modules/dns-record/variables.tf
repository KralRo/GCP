variable "manage_dns" {
  type        = bool
  description = "Whether Cloud DNS manages this client's domain. Set false when the domain stays with the client's own registrar - the records output still lists what to set up manually there."
  default     = false
}

variable "create_zone" {
  type        = bool
  description = "Whether to create a new Cloud DNS managed zone (true) or reference one that already exists in this project, e.g. auto-provisioned by Cloud Domains at registration (false). Only relevant when manage_dns = true."
  default     = true
}

variable "project_id" {
  type        = string
  description = "GCP project ID the managed zone belongs to (only used when manage_dns = true)"
}

variable "domain" {
  type        = string
  description = "Client's domain name, e.g. example.com"
}

variable "records" {
  type = list(object({
    name    = string
    type    = string
    ttl     = number
    rrdatas = list(string)
  }))
  description = "DNS records to publish, e.g. the A/CNAME record pointing at the Cloud Run service"
}

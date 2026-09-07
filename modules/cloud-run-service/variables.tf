variable "project_id" {
  type        = string
  description = "GCP project ID the service is deployed into"
}

variable "region" {
  type        = string
  description = "Cloud Run region"
  default     = "europe-west1"
}

variable "service_name" {
  type        = string
  description = "Cloud Run service name, e.g. svc-prod or svc-dev"
}

variable "image" {
  type        = string
  description = "Container image URL to deploy"
}

variable "env_vars" {
  type        = map(string)
  description = "Environment variables passed to the container as plain values"
  default     = {}
}

variable "secret_env_vars" {
  type = map(object({
    secret  = string
    version = optional(string, "latest")
  }))
  description = "Environment variables sourced from a Secret Manager secret version instead of a plain value. `secret` is the Secret Manager secret ID (or full resource name); the runtime service account must already have roles/secretmanager.secretAccessor on it."
  default     = {}
}

variable "service_account" {
  type        = string
  description = "Email of the service account the container runs as. Empty uses the project's default compute service account."
  default     = ""
}

variable "allow_unauthenticated" {
  type        = bool
  description = "Whether to allow public unauthenticated access"
  default     = true
}

variable "min_instances" {
  type        = number
  description = "Minimum number of container instances (0 = scale to zero)"
  default     = 0
}

variable "max_instances" {
  type        = number
  description = "Maximum number of container instances"
  default     = 4
}

variable "cpu" {
  type        = string
  description = "CPU allocation per instance"
  default     = "1"
}

variable "memory" {
  type        = string
  description = "Memory allocation per instance"
  default     = "512Mi"
}

variable "domain_mappings" {
  type        = list(string)
  description = "Hostnames to map to this Cloud Run service via google_cloud_run_domain_mapping, e.g. [\"streckerova.kralroman.org\", \"admin.streckerova.kralroman.org\"]. Domain ownership must be verified (Search Console) before apply, or the mapping resource will fail to create."
  default     = []
}

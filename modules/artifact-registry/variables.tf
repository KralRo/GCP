variable "project_id" {
  type        = string
  description = "GCP project ID the repository belongs to"
}

variable "region" {
  type        = string
  description = "Region for the Artifact Registry repository"
  default     = "europe-west1"
}

variable "repository_id" {
  type        = string
  description = "Repository name, e.g. app"
  default     = "app"
}

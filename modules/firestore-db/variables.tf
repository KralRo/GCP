variable "project_id" {
  type        = string
  description = "GCP project ID the database belongs to"
}

variable "database_id" {
  type        = string
  description = "Firestore database ID, e.g. (default) for prod or dev for the dev environment"
}

variable "location_id" {
  type        = string
  description = "Firestore location (multi-region or region)"
  default     = "eur3"
}

terraform {
  required_version = ">= 1.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }

  # Bucket/prefix are provided per client via backend.hcl at `terraform init -backend-config=backend.hcl`.
  backend "gcs" {}
}

provider "google" {
  project = var.project_id
  region  = var.region
}

locals {
  firestore_database_id = var.env == "prod" ? "(default)" : var.env
}

# Identity the Cloud Run service runs as - separate from the github-actions-deployer
# SA used by CI, scoped to just what the running app needs (Firestore read/write).
resource "google_service_account" "runtime" {
  project      = var.project_id
  account_id   = "svc-${var.env}-runtime"
  display_name = "Cloud Run runtime (${var.env})"
}

resource "google_project_iam_member" "runtime_datastore" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.runtime.email}"
}

# Signs the admin session cookie. Stored as a plain Cloud Run env var (visible in
# Terraform state) rather than Secret Manager - acceptable for a single-editor admin
# panel today; revisit if that changes.
resource "random_password" "session_secret" {
  length  = 32
  special = false
}

module "artifact_registry" {
  source = "../../../modules/artifact-registry"

  project_id = var.project_id
  region     = var.region
}

module "cloud_run" {
  source = "../../../modules/cloud-run-service"

  project_id      = var.project_id
  region          = var.region
  service_name    = "svc-${var.env}"
  image           = var.image
  service_account = google_service_account.runtime.email

  env_vars = {
    FIRESTORE_PROJECT_ID  = var.project_id
    FIRESTORE_DATABASE_ID = local.firestore_database_id
    SESSION_SECRET        = random_password.session_secret.result
  }

  depends_on = [google_project_iam_member.runtime_datastore]
}

module "firestore" {
  source = "../../../modules/firestore-db"

  project_id  = var.project_id
  database_id = local.firestore_database_id
  location_id = var.firestore_location
}

module "dns" {
  source = "../../../modules/dns-record"

  manage_dns = var.manage_dns
  project_id = var.project_id
  domain     = var.domain
  records    = var.dns_records
}

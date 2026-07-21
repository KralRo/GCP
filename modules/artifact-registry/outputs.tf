output "repository_url" {
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repository_id}"
  description = "Base URL to push/pull images against, e.g. <url>/svc:<tag>"
}

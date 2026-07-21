output "url" {
  value       = google_cloud_run_v2_service.this.uri
  description = "The URL Cloud Run assigned to the service"
}

output "service_name" {
  value       = google_cloud_run_v2_service.this.name
  description = "The deployed Cloud Run service name"
}

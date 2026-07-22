output "url" {
  value       = google_cloud_run_v2_service.this.uri
  description = "The URL Cloud Run assigned to the service"
}

output "service_name" {
  value       = google_cloud_run_v2_service.this.name
  description = "The deployed Cloud Run service name"
}

output "domain_mapping_records" {
  value       = { for k, v in google_cloud_run_domain_mapping.this : k => v.status[0].resource_records }
  description = "DNS records Google Cloud Run needs for each mapped domain - set these up wherever that domain's DNS is actually hosted"
}

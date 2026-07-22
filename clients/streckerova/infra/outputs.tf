output "cloud_run_url" {
  value       = module.cloud_run.url
  description = "URL of the deployed Cloud Run service"
}

output "artifact_registry_url" {
  value       = module.artifact_registry.repository_url
  description = "Base URL to push/pull the service image against"
}

output "dns_records_to_configure" {
  value       = module.dns.records_to_configure
  description = "DNS records to set up - at the registrar if manage_dns = false"
}

output "domain_mapping_records" {
  value       = module.cloud_run.domain_mapping_records
  description = "DNS records to set up for each entry in domain_mappings"
}

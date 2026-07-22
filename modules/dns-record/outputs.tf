output "name_servers" {
  value       = var.manage_dns ? (var.create_zone ? google_dns_managed_zone.this[0].name_servers : data.google_dns_managed_zone.existing[0].name_servers) : null
  description = "Name servers to set at the registrar, if Cloud DNS manages this zone"
}

output "records_to_configure" {
  value       = var.records
  description = "Records to configure - at the client's registrar when manage_dns = false, informational otherwise"
}

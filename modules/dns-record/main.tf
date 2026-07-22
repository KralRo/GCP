resource "google_dns_managed_zone" "this" {
  count       = var.manage_dns && var.create_zone ? 1 : 0
  project     = var.project_id
  name        = replace(var.domain, ".", "-")
  dns_name    = "${var.domain}."
  description = "Managed zone for ${var.domain}"
}

# Cloud Domains auto-provisions a matching Cloud DNS zone at registration -
# reference that instead of creating a duplicate when create_zone = false.
data "google_dns_managed_zone" "existing" {
  count   = var.manage_dns && !var.create_zone ? 1 : 0
  project = var.project_id
  name    = replace(var.domain, ".", "-")
}

locals {
  zone_name = var.create_zone ? try(google_dns_managed_zone.this[0].name, null) : try(data.google_dns_managed_zone.existing[0].name, null)
}

resource "google_dns_record_set" "this" {
  for_each = var.manage_dns ? { for r in var.records : r.name => r } : {}

  project      = var.project_id
  managed_zone = local.zone_name
  name         = "${each.value.name}.${var.domain}."
  type         = each.value.type
  ttl          = each.value.ttl
  rrdatas      = each.value.rrdatas
}

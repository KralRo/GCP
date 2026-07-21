resource "google_dns_managed_zone" "this" {
  count       = var.manage_dns ? 1 : 0
  project     = var.project_id
  name        = replace(var.domain, ".", "-")
  dns_name    = "${var.domain}."
  description = "Managed zone for ${var.domain}"
}

resource "google_dns_record_set" "this" {
  for_each = var.manage_dns ? { for r in var.records : r.name => r } : {}

  project      = var.project_id
  managed_zone = google_dns_managed_zone.this[0].name
  name         = "${each.value.name}.${var.domain}."
  type         = each.value.type
  ttl          = each.value.ttl
  rrdatas      = each.value.rrdatas
}

client_slug = "streckerova"
project_id  = "rk-streckerova"
region      = "europe-west1"
env         = "prod"

domain_mappings = ["streckerova.kralroman.org", "admin.streckerova.kralroman.org"]
admin_host      = "admin.streckerova.kralroman.org"

# Cloud DNS manages kralroman.org itself (Roman's own domain - previously
# lived in a different, now-deleted GCP project, so this is a fresh zone).
# CNAME target is Google's documented value for Cloud Run domain mappings on
# a subdomain; cross-check against the `domain_mapping_records` output after
# apply and correct here if Google returns something different.
manage_dns = true
domain     = "kralroman.org"
dns_records = [
  { name = "streckerova", type = "CNAME", ttl = 300, rrdatas = ["ghs.googlehosted.com."] },
  { name = "admin.streckerova", type = "CNAME", ttl = 300, rrdatas = ["ghs.googlehosted.com."] },
]

client_slug = "streckerova"
project_id  = "rk-streckerova"
region      = "europe-west1"
env         = "prod"

domain_mappings = ["streckerova.kralroman.org", "admin.streckerova.kralroman.org", "stremi.cz", "www.stremi.cz"]
admin_host      = "admin.streckerova.kralroman.org"

# Cloud Domains auto-provisioned a Cloud DNS zone "kralroman-org" for this
# domain in this same project at registration - reference it, don't create
# a duplicate (create_zone = false). CNAME target is Google's documented
# value for Cloud Run domain mappings on a subdomain; cross-check against the
# `domain_mapping_records` output after apply and correct here if Google
# returns something different.
manage_dns  = true
create_zone = false
domain      = "kralroman.org"
dns_records = [
  { name = "streckerova", type = "CNAME", ttl = 300, rrdatas = ["ghs.googlehosted.com."] },
  { name = "admin.streckerova", type = "CNAME", ttl = 300, rrdatas = ["ghs.googlehosted.com."] },
]

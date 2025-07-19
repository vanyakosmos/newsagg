resource "google_artifact_registry_repository" "docker" {
  location      = var.gcp_region
  repository_id = "docker"
  format        = "DOCKER"
  vulnerability_scanning_config {
    enablement_config = "DISABLED"
  }
  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "keep-5-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count = 10
    }
  }
}

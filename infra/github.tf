resource "github_actions_variable" "docker_registry" {
  repository    = var.github_repo
  variable_name = "docker_registry"
  value         = "${var.gcp_region}-docker.pkg.dev"
  depends_on    = [google_artifact_registry_repository.docker]
}

resource "github_actions_variable" "image_name" {
  repository    = var.github_repo
  variable_name = "image_name"
  value         = local.image_name
}

resource "github_actions_variable" "gh_actions_sa" {
  repository    = var.github_repo
  variable_name = "gh_actions_sa"
  value         = google_service_account.github.email
}

resource "github_actions_variable" "workload_identity_provider" {
  repository    = var.github_repo
  variable_name = "workload_identity_provider"
  value         = google_iam_workload_identity_pool_provider.github_provider.name
}

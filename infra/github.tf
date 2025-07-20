# ENV VARS

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

# SECRETS

resource "github_actions_secret" "gcp_sa_key" {
  repository      = var.github_repo
  secret_name     = "gcp_sa_key"
  plaintext_value = base64decode(google_service_account_key.github.private_key)
}

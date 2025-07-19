# ENV VARS

resource "github_actions_variable" "gcp_project" {
  repository    = var.github_repo
  variable_name = "gcp_project"
  value         = var.gcp_project
}

resource "github_actions_variable" "gcp_region" {
  repository    = var.github_repo
  variable_name = "gcp_region"
  value         = var.gcp_region
}

resource "github_actions_variable" "docker_registry" {
  repository    = var.github_repo
  variable_name = "docker_registry"
  value         = "${var.gcp_region}-docker.pkg.dev"
  depends_on    = [google_artifact_registry_repository.docker]
}

resource "github_actions_variable" "image_name" {
  repository    = var.github_repo
  variable_name = "image_name"
  value         = "${github_actions_variable.docker_registry.value}/${var.gcp_project}/docker/newsagg"
}

resource "github_actions_variable" "bucket_name" {
  repository    = var.github_repo
  variable_name = "bucket_name"
  value         = google_storage_bucket.newsagg.name
}

# SERVICE ACCOUNTS

resource "github_actions_variable" "app_sa" {
  repository    = var.github_repo
  variable_name = "app_sa"
  value         = google_service_account.application.email
}

resource "github_actions_variable" "cloudscheduler_sa" {
  repository    = var.github_repo
  variable_name = "cloudscheduler_sa"
  value         = google_service_account.cloudscheduler.email
}

# SECRETS

resource "github_actions_secret" "gcp_sa_key" {
  repository      = var.github_repo
  secret_name     = "gcp_sa_key"
  plaintext_value = base64decode(google_service_account_key.github.private_key)
}

resource "github_actions_secret" "telegram_bot_token" {
  repository      = var.github_repo
  secret_name     = "telegram_bot_token"
  plaintext_value = var.telegram_bot_token
}

resource "github_actions_secret" "bucket_access_key" {
  repository      = var.github_repo
  secret_name     = "bucket_access_key"
  plaintext_value = google_storage_hmac_key.newsagg.access_id
}

resource "github_actions_secret" "bucket_secret_key" {
  repository      = var.github_repo
  secret_name     = "bucket_secret_key"
  plaintext_value = google_storage_hmac_key.newsagg.secret
}
